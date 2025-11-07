package saga

import (
	"context"
	"errors"
	"sync"
)

var _ Store = (*StoreInMemory)(nil)

type StoreInMemory struct {
	mu   sync.RWMutex
	data map[string]Instance
}

func NewInMemoryStore() *StoreInMemory {
	return &StoreInMemory{data: map[string]Instance{}}
}

func (m *StoreInMemory) Create(_ context.Context, s *Instance) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.data[s.ID]; ok {
		return ErrConflict
	}
	cpy := *s
	m.data[s.ID] = cpy
	return nil
}

func (m *StoreInMemory) Load(_ context.Context, id string) (*Instance, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	v, ok := m.data[id]
	if !ok {
		return nil, ErrNotFound
	}

	cloned := cloneInstance(v)
	return &cloned, nil
}

func (m *StoreInMemory) Save(_ context.Context, s *Instance) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if s == nil {
		return errors.New("cannot save nil instance")
	}

	cur, ok := m.data[s.ID]
	if !ok {
		return ErrNotFound
	}
	// optimistic concurrency: revision must match
	if s.Revision != cur.Revision {
		return ErrConflict
	}
	s.IncrementRevision()
	m.data[s.ID] = cloneInstance(*s)
	return nil
}

func cloneInstance(in Instance) Instance {
	out := in
	out.Steps = make([]StepState, len(in.Steps))
	for i := range in.Steps {
		ss := in.Steps[i]
		cp := ss
		if ss.Data != nil {
			cp.Data = map[string]any{}
			for k, v := range ss.Data {
				cp.Data[k] = v
			}
		}
		out.Steps[i] = cp
	}
	if in.Input != nil {
		out.Input = map[string]any{}
		for k, v := range in.Input {
			out.Input[k] = v
		}
	}
	if in.Metadata != nil {
		out.Metadata = map[string]string{}
		for k, v := range in.Metadata {
			out.Metadata[k] = v
		}
	}
	return out
}
