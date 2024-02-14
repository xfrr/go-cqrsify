package aggregate

// Aggregate represents a domain-driven design and event-sourced aggregate.
type Aggregate interface {
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
