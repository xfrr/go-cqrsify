package retry

import (
	"context"
	"sync"
	"time"
)

// InMemoryDedupeStore is a simple, process-local store.
type InMemoryDedupeStore struct {
	mu   sync.Mutex
	keys map[string]time.Time // expiry
}

func NewInMemoryDedupeStore() *InMemoryDedupeStore {
	return &InMemoryDedupeStore{keys: map[string]time.Time{}}
}

func (s *InMemoryDedupeStore) Begin(_ context.Context, token string, ttl time.Duration) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	now := time.Now()
	if exp, ok := s.keys[token]; ok && exp.After(now) {
		return false, nil
	}
	s.keys[token] = now.Add(ttl)
	return true, nil
}
func (s *InMemoryDedupeStore) Commit(_ context.Context, token string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.keys, token)
	return nil
}
func (s *InMemoryDedupeStore) Rollback(_ context.Context, token string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.keys, token)
	return nil
}
