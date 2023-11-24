package rate_limiter

import (
	"sync"
	"time"
)

// Record .
type Record struct {
	Data       uint
	Expiration time.Time
}

// TTLMap .
type TTLMap struct {
	data  map[string]Record
	Limit uint
	TTL   uint
	mu    sync.Mutex
}

// NewTTLMap .
func NewTTLMap(limit, ttl uint) *TTLMap {
	return &TTLMap{
		data:  make(map[string]Record),
		Limit: limit,
		TTL:   ttl,
	}
}

// Set .
func (m *TTLMap) Set(key string, value uint) {
	m.mu.Lock()
	defer m.mu.Unlock()

	TTL := time.Duration(m.TTL) * time.Second
	expiration := time.Now().Add(TTL)

	m.data[key] = Record{
		Data:       value,
		Expiration: expiration,
	}
}

// Get .
func (m *TTLMap) Get(key string) uint {
	m.mu.Lock()
	defer m.mu.Unlock()

	record, exists := m.data[key]
	if !exists {
		return 0
	}

	if time.Now().After(record.Expiration) {
		delete(m.data, key)
		return 0
	}

	return record.Data
}
