package domain

import "context"

// EventStore represents an event store that can save and retrieve events.
type EventStore[ID comparable] interface {
	EventSaver
	EventRetriever[ID]
	EventSearcher
}

// EventSaver represents an event store that can save events.
type EventSaver interface {
	// Save saves the given events to the event store.
	Save(ctx context.Context, events []Event) error
}

// EventRetriever represents an event store that can retrieve events.
type EventRetriever[ID comparable] interface {
	// RetrieveMany retrieves events for the given aggregate ID from the event store.
	RetrieveMany(ctx context.Context, aggregateID ID, opts ...RetrieveEventsOption) ([]Event, error)
}

// EventSearcher represents an event store that can search for events.
type EventSearcher interface {
	// Search searches for events in the event store that match the given criteria.
	Search(ctx context.Context, criteria *SearchCriteriaOptions) ([]Event, error)
}

// RetrieveEventsOption represents an option for retrieving events.
type RetrieveEventsOption func(*RetrieveEventsOptions)

// RetrieveEventsFromVersion specifies the version from which to start retrieving events (inclusive).
func RetrieveEventsFromVersion(version int) RetrieveEventsOption {
	return func(opts *RetrieveEventsOptions) {
		opts.FromVersion = version
	}
}

// RetrieveEventsToVersion specifies the version up to which to retrieve events (inclusive).
func RetrieveEventsToVersion(version int) RetrieveEventsOption {
	return func(opts *RetrieveEventsOptions) {
		opts.ToVersion = version
	}
}

// RetrieveEventsBatchSize specifies the maximum number of events to retrieve in a single batch.
func RetrieveEventsBatchSize(size int) RetrieveEventsOption {
	return func(opts *RetrieveEventsOptions) {
		opts.BatchSize = size
	}
}

type RetrieveEventsOptions struct {
	// FromVersion specifies the version from which to start retrieving events (inclusive).
	FromVersion int
	// ToVersion specifies the version up to which to retrieve events (inclusive).
	ToVersion int
	// BatchSize specifies the maximum number of events to retrieve in a single batch.
	BatchSize int
}
