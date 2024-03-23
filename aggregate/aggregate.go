package aggregate

// Aggregate represents a domain-driven design and event-sourced aggregate.
type Aggregate[ID comparable] interface {
	// AggregateID returns the aggregate's AggregateID.
	AggregateID() ID

	// AggregateName returns the aggregate's name.
	AggregateName() string

	// AggregateChanges returns the aggregate's events.
	AggregateChanges() []Change

	// AggregateVersion returns the aggregate's version.
	AggregateVersion() Version

	// ChangeApplier applies changes (events) to the aggregate.
	ChangeApplier
}

func Cast[OutID comparable, InID comparable](
	a Aggregate[InID],
) (*Base[OutID], bool) {
	id, ok := any(a.AggregateID()).(OutID)
	if !ok {
		return nil, false
	}

	return &Base[OutID]{
		id:      id,
		name:    a.AggregateName(),
		version: a.AggregateVersion(),
		changes: a.AggregateChanges(),
	}, true

}
