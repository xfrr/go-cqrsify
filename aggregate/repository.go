package aggregate

import (
	"context"
	"errors"
)

var (
	// ErrNotFound is returned when the aggregate is not found in the repository.
	ErrNotFound = errors.New("aggregate not found")
)

// Repository is the interface that wraps the basic methods for managing the
// lifecycle of an aggregate.
type Repository[ID comparable] interface {
	// Delete deletes the aggregate from the repository.
	Delete(context.Context, Aggregate[ID]) error
	// Exists checks if the aggregate with the given id exists in the repository.
	Exists(context.Context, Aggregate[ID]) (bool, error)
	// Load loads the aggregate with the given id from the repository.
	Load(context.Context, Aggregate[ID]) error
	// Search loads all the aggregates from the repository that match the given options.
	Search(context.Context, *SearchCriteriaOptions) ([]Aggregate[ID], error)
	// Save saves the aggregate to the repository.
	Save(context.Context, Aggregate[ID]) error
}

// VersionedRepository is the interface that wraps the basic methods for managing the
// lifecycle of an aggregate with versioning support.
type VersionedRepository[ID comparable] interface {
	Repository[ID]

	// LoadVersion loads the aggregate with the given id and version from the repository.
	LoadVersion(context.Context, Aggregate[ID], Version) error

	// ExistsVersion checks if the aggregate with the given id and version exists in the repository.
	ExistsVersion(context.Context, Aggregate[ID], Version) (bool, error)
}
