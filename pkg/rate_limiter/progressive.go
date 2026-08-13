package rate_limiter

import (
	"sync"
	"time"
)

// ProgressiveLimiter locks a key for escalating durations after repeated attempt bursts.
// Example: 5 attempts → 60s lock, next 5 → 1h, next 5 → 24h.
type ProgressiveLimiter struct {
	mu               sync.Mutex
	data             map[string]*progressiveRecord
	attemptsPerStage uint
	lockouts         []time.Duration
	stopCh           chan struct{}
}

type progressiveRecord struct {
	attempts    uint
	stage       int
	lockedUntil time.Time
}

// NewProgressive creates a progressive limiter.
// attemptsPerStage is how many failures trigger a lock (e.g. 5).
// lockouts are lock durations per escalation stage (e.g. 1m, 1h, 24h).
func NewProgressive(attemptsPerStage uint, lockouts ...time.Duration) *ProgressiveLimiter {
	if attemptsPerStage == 0 {
		attemptsPerStage = 5
	}
	if len(lockouts) == 0 {
		lockouts = []time.Duration{time.Minute, time.Hour, 24 * time.Hour}
	}

	p := &ProgressiveLimiter{
		data:             make(map[string]*progressiveRecord),
		attemptsPerStage: attemptsPerStage,
		lockouts:         lockouts,
		stopCh:           make(chan struct{}),
	}
	go p.cleanupLoop()
	return p
}

// Locked reports whether key is currently blocked and how long until unlock.
func (p *ProgressiveLimiter) Locked(key string) (locked bool, retryAfter time.Duration) {
	p.mu.Lock()
	defer p.mu.Unlock()

	rec := p.data[key]
	if rec == nil {
		return false, 0
	}

	now := time.Now()
	if now.Before(rec.lockedUntil) {
		return true, rec.lockedUntil.Sub(now)
	}
	return false, 0
}

// RecordFailure counts a failed attempt. When the stage limit is hit, the key is locked
// for the duration of the current stage (then stage escalates).
func (p *ProgressiveLimiter) RecordFailure(key string) (locked bool, retryAfter time.Duration) {
	p.mu.Lock()
	defer p.mu.Unlock()

	now := time.Now()
	rec := p.data[key]
	if rec == nil {
		rec = &progressiveRecord{}
		p.data[key] = rec
	}

	if now.Before(rec.lockedUntil) {
		return true, rec.lockedUntil.Sub(now)
	}

	rec.attempts++
	if rec.attempts < p.attemptsPerStage {
		return false, 0
	}

	duration := p.lockouts[len(p.lockouts)-1]
	if rec.stage < len(p.lockouts) {
		duration = p.lockouts[rec.stage]
	}

	rec.lockedUntil = now.Add(duration)
	rec.stage++
	rec.attempts = 0
	return true, duration
}

// Reset clears all state for key (e.g. after successful login).
func (p *ProgressiveLimiter) Reset(key string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	delete(p.data, key)
}

// BlockInfo is a snapshot entry for admin unblock UI/API.
type BlockInfo struct {
	Key           string     `json:"key"`
	Attempts      uint       `json:"attempts"`
	Stage         int        `json:"stage"`
	Locked        bool       `json:"locked"`
	LockedUntil   *time.Time `json:"locked_until,omitempty"`
	RetryAfterSec int        `json:"retry_after_sec,omitempty"`
}

// Snapshot returns current limiter entries (locked and in-progress).
func (p *ProgressiveLimiter) Snapshot() []BlockInfo {
	p.mu.Lock()
	defer p.mu.Unlock()

	now := time.Now()
	out := make([]BlockInfo, 0, len(p.data))
	for key, rec := range p.data {
		info := BlockInfo{
			Key:      key,
			Attempts: rec.attempts,
			Stage:    rec.stage,
		}
		if now.Before(rec.lockedUntil) {
			until := rec.lockedUntil
			info.Locked = true
			info.LockedUntil = &until
			info.RetryAfterSec = int(rec.lockedUntil.Sub(now).Seconds())
			if info.RetryAfterSec < 1 {
				info.RetryAfterSec = 1
			}
		}
		out = append(out, info)
	}
	return out
}

// Stop stops the background cleanup goroutine.
func (p *ProgressiveLimiter) Stop() {
	select {
	case <-p.stopCh:
		return
	default:
		close(p.stopCh)
	}
}

func (p *ProgressiveLimiter) cleanupLoop() {
	ticker := time.NewTicker(time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-p.stopCh:
			return
		case <-ticker.C:
			p.cleanup()
		}
	}
}

func (p *ProgressiveLimiter) cleanup() {
	p.mu.Lock()
	defer p.mu.Unlock()

	now := time.Now()
	for key, rec := range p.data {
		// drop records that are unlocked and have no pending attempts
		if now.After(rec.lockedUntil) && rec.attempts == 0 && rec.stage == 0 {
			delete(p.data, key)
			continue
		}
		// drop fully expired long-idle escalations after last lock + idle
		if now.After(rec.lockedUntil) && rec.attempts == 0 && now.Sub(rec.lockedUntil) > 24*time.Hour {
			delete(p.data, key)
		}
	}
}
