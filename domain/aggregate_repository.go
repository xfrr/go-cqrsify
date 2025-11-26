package domain

import (
	"context"
)

// Repository is the interface that wraps the basic methods for managing the
// lifecycle of an aggregate.
type Repository[T Aggregate[ID], ID comparable] interface {
	// Exists checks if the aggregate with the given id exists in the repository.
	Exists(context.Context, ID) (bool, error)
	// Get retrieves an aggregate by its ID.
	Get(context.Context, ID) (T, error)
	// Save saves the aggregate to the repository.
	Save(context.Context, Aggregate[ID]) error
}

// VersionedRepository is the interface that wraps the basic methods for managing the
// lifecycle of an aggregate with versioning support.
type VersionedRepository[T VersionedAggregate[ID], ID comparable] interface {
	Repository[Aggregate[ID], ID]

	// GetVersion retrieves an aggregate by its ID and version.
	GetVersion(context.Context, ID, AggregateVersion) (T, error)

	// ExistsVersion checks if the aggregate with the given id and version exists in the repository.
	ExistsVersion(context.Context, ID, AggregateVersion) (bool, error)
}

// SearchableRepository is the interface that wraps the basic methods for managing the
// lifecycle of an aggregate with search capabilities.
type SearchableRepository[T Aggregate[ID], ID comparable] interface {
	Repository[T, ID]

	// Search loads all the aggregates from the repository that match the given options.
	Search(context.Context, *SearchCriteriaOptions) ([]T, error)
}

// EventSourcedRepository is the interface that wraps the basic methods for managing the
// lifecycle of an event-sourced aggregate.
type EventSourcedRepository[T EventSourcedAggregate[ID], ID comparable] interface {
	// Exists checks if the aggregate with the given id exists in the repository.
	Exists(context.Context, T) (bool, error)
	// Load loads the aggregate with the given id from the repository.
	Load(context.Context, T) error
	// Save saves the aggregate to the repository.
	Save(context.Context, T) error
	// LoadVersion loads the aggregate with the given id and version from the repository.
	LoadVersion(context.Context, T, AggregateVersion) error
	// ExistsVersion checks if the aggregate with the given id and version exists in the repository.
	ExistsVersion(context.Context, T, AggregateVersion) (bool, error)
}

type SearchableEventSourcedRepository[T EventSourcedAggregate[ID], ID comparable] interface {
	EventSourcedRepository[T, ID]

	// Search loads all the aggregates from the repository that match the given options.
	Search(context.Context, *SearchCriteriaOptions) ([]T, error)
}

// NotFoundError is returned when an aggregate is not found in the repository.
type NotFoundError[ID comparable] struct {
	ID ID
}

func (e NotFoundError[ID]) Error() string {
	return "aggregate not found"
}

func NewNotFoundError[ID comparable](id ID) NotFoundError[ID] {
	return NotFoundError[ID]{ID: id}
}
