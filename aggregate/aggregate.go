package aggregate

// Aggregate represents a domain-driven design and event-sourced aggregate.
type Aggregate[ID comparable] interface {
	// AggregateID returns the aggregate's AggregateID.
	AggregateID() ID

	// AggregateName returns the aggregate's name.
	AggregateName() string

	// AggregateEvents returns the aggregate's events.
	AggregateEvents() []Event

	// AggregateVersion returns the aggregate's version.
	AggregateVersion() Version

	// EventApplier applies events to the aggregate.
	EventApplier
}

func Cast[OutID comparable, InID comparable](
	a Aggregate[InID],
) (*Base[OutID], bool) {
	id, ok := any(a.AggregateID()).(OutID)
	if !ok {
		return nil, false
	}

	return &Base[OutID]{
		id:       id,
		name:     a.AggregateName(),
		version:  a.AggregateVersion(),
		events:   a.AggregateEvents(),
		handlers: make(map[string][]func(Event)),
	}, true
}
