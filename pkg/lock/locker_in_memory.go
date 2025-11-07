package lock

import (
	"context"
	"sync"
	"time"
)

var _ Locker = (*InMemoryLocker)(nil)

// InMemoryLocker is a simple in-memory lock implementation.
type InMemoryLocker struct {
	mu   sync.Mutex
	keys map[string]time.Time
}

// NewInMemoryLocker creates a new InMemoryLocker.
func NewInMemoryLocker() *InMemoryLocker {
	return &InMemoryLocker{keys: map[string]time.Time{}}
}

func (l *InMemoryLocker) TryLock(_ context.Context, key string, ttl time.Duration) (bool, error) {
	l.mu.Lock()
	defer l.mu.Unlock()
	now := time.Now()
	exp, ok := l.keys[key]
	if !ok || exp.Before(now) {
		l.keys[key] = now.Add(ttl)
		return true, nil
	}
	return false, nil
}

func (l *InMemoryLocker) Unlock(_ context.Context, key string) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	delete(l.keys, key)
	return nil
}

func (l *InMemoryLocker) Refresh(_ context.Context, key string, ttl time.Duration) (bool, error) {
	l.mu.Lock()
	defer l.mu.Unlock()
	_, ok := l.keys[key]
	if !ok {
		return false, nil
	}
	l.keys[key] = time.Now().Add(ttl)
	return true, nil
}
