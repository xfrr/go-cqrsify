package saga

import (
	"context"
)

type Store interface {
	// Create a new instance; must fail if ID already exists.
	Create(ctx context.Context, s *Instance) error
	// Load instance by ID.
	Load(ctx context.Context, id string) (*Instance, error)
	// Save with optimistic lock on Revision (increment on success).
	Save(ctx context.Context, s *Instance) error
}
