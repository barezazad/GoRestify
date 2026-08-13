package rate_limiter

import (
	"sync"
	"time"
)

// Limiter is an in-memory fixed-window rate limiter with expired-key cleanup.
type Limiter struct {
	mu     sync.Mutex
	data   map[string]record
	limit  uint
	window time.Duration
	stopCh chan struct{}
}

type record struct {
	count     uint
	expiresAt time.Time
}

// New creates a limiter that allows `limit` requests per `window` per key.
// A background cleanup runs periodically to avoid unbounded memory growth.
func New(limit uint, window time.Duration) *Limiter {
	if limit == 0 {
		limit = 1
	}
	if window <= 0 {
		window = time.Minute
	}

	l := &Limiter{
		data:   make(map[string]record),
		limit:  limit,
		window: window,
		stopCh: make(chan struct{}),
	}

	go l.cleanupLoop()
	return l
}

// Allow increments the counter for key and returns true if under the limit.
func (l *Limiter) Allow(key string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()
	rec, ok := l.data[key]

	if !ok || now.After(rec.expiresAt) {
		l.data[key] = record{
			count:     1,
			expiresAt: now.Add(l.window),
		}
		return true
	}

	if rec.count >= l.limit {
		return false
	}

	rec.count++
	l.data[key] = rec
	return true
}

// Remaining returns how many requests are left for key in the current window.
func (l *Limiter) Remaining(key string) uint {
	l.mu.Lock()
	defer l.mu.Unlock()

	rec, ok := l.data[key]
	if !ok || time.Now().After(rec.expiresAt) {
		return l.limit
	}
	if rec.count >= l.limit {
		return 0
	}
	return l.limit - rec.count
}

// Reset clears the counter for key.
func (l *Limiter) Reset(key string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	delete(l.data, key)
}

// Stop stops the background cleanup goroutine.
func (l *Limiter) Stop() {
	select {
	case <-l.stopCh:
		return
	default:
		close(l.stopCh)
	}
}

func (l *Limiter) cleanupLoop() {
	ticker := time.NewTicker(l.window)
	defer ticker.Stop()

	for {
		select {
		case <-l.stopCh:
			return
		case <-ticker.C:
			l.cleanup()
		}
	}
}

func (l *Limiter) cleanup() {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()
	for key, rec := range l.data {
		if now.After(rec.expiresAt) {
			delete(l.data, key)
		}
	}
}

// Keep old names as thin wrappers for any existing callers.
type (
	// TTLMap is deprecated alias kept for compatibility.
	TTLMap = Limiter
	// Record is deprecated.
	Record = record
)

// NewTTLMap keeps the previous constructor signature (limit, ttlSeconds).
func NewTTLMap(limit, ttlSeconds uint) *Limiter {
	return New(limit, time.Duration(ttlSeconds)*time.Second)
}

// Set is kept for compatibility; prefer Allow.
func (l *Limiter) Set(key string, value uint) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.data[key] = record{
		count:     value,
		expiresAt: time.Now().Add(l.window),
	}
}

// Get is kept for compatibility; prefer Remaining / Allow.
func (l *Limiter) Get(key string) uint {
	l.mu.Lock()
	defer l.mu.Unlock()

	rec, exists := l.data[key]
	if !exists || time.Now().After(rec.expiresAt) {
		if exists {
			delete(l.data, key)
		}
		return 0
	}
	return rec.count
}
