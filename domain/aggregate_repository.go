package domain

import (
	"context"
)

// Repository is the interface that wraps the basic methods for managing the
// lifecycle of an aggregate.
type Repository[ID comparable] interface {
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
	LoadVersion(context.Context, VersionedAggregate[ID], AggregateVersion) error

	// ExistsVersion checks if the aggregate with the given id and version exists in the repository.
	ExistsVersion(context.Context, VersionedAggregate[ID], AggregateVersion) (bool, error)
}

// EventSourcedRepository is the interface that wraps the basic methods for managing the
// lifecycle of an event-sourced aggregate.
type EventSourcedRepository[ID comparable] interface {
	// Exists checks if the aggregate with the given id exists in the repository.
	Exists(context.Context, EventSourcedAggregate[ID]) (bool, error)
	// Load loads the aggregate with the given id from the repository.
	Load(context.Context, EventSourcedAggregate[ID]) error
	// Search loads all the aggregates from the repository that match the given options.
	Search(context.Context, *SearchCriteriaOptions) ([]EventSourcedAggregate[ID], error)
	// Save saves the aggregate to the repository.
	Save(context.Context, EventSourcedAggregate[ID]) error
	// LoadVersion loads the aggregate with the given id and version from the repository.
	LoadVersion(context.Context, EventSourcedAggregate[ID], AggregateVersion) error
	// ExistsVersion checks if the aggregate with the given id and version exists in the repository.
	ExistsVersion(context.Context, EventSourcedAggregate[ID], AggregateVersion) (bool, error)
}

type NotFoundError[ID comparable] struct {
	ID ID
}

func (e NotFoundError[ID]) Error() string {
	return "aggregate not found"
}

func NewNotFoundError[ID comparable](id ID) NotFoundError[ID] {
	return NotFoundError[ID]{ID: id}
}
