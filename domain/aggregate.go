package domain

// Aggregate represents a domain-driven design and event-sourced aggregate.
type Aggregate[ID comparable] interface {
	// AggregateID returns the aggregate's AggregateID.
	AggregateID() ID

	// AggregateName returns the aggregate's name.
	AggregateName() string
}

type VersionedAggregate[ID comparable] interface {
	Aggregate[ID]

	// AggregateVersion returns the aggregate's version.
	AggregateVersion() AggregateVersion
}

type EventSourcedAggregate[ID comparable] interface {
	VersionedAggregate[ID]

	// AggregateEvents returns the aggregate's events.
	AggregateEvents() []Event

	// ApplyEvent applies the given event to the aggregate.
	ApplyEvent(Event) error

	// HandleEvent registers an event handler for the given event name on the aggregate.
	HandleEvent(eventName string, handler func(event Event) error)
}

func CastAggregate[OutID comparable, InID comparable](a Aggregate[InID]) (*BaseAggregate[OutID], bool) {
	id, ok := any(a.AggregateID()).(OutID)
	if !ok {
		return nil, false
	}

	var version AggregateVersion
	if va, ok := a.(VersionedAggregate[InID]); ok {
		version = va.AggregateVersion()
	}

	var events []Event
	if ea, ok := a.(EventSourcedAggregate[InID]); ok {
		events = ea.AggregateEvents()
	}

	return &BaseAggregate[OutID]{
		id:       id,
		name:     a.AggregateName(),
		version:  version,
		events:   events,
		handlers: map[string][]func(Event) error{},
	}, true
}
