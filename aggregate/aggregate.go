package aggregate

// Aggregate is the interface that wraps the basic functionality
// of a domain-driven design aggregate.
type Aggregate interface {
	// AggregateID returns the aggregate's AggregateID.
	AggregateID() ID

	// AggregateName returns the aggregate's name.
	AggregateName() string

	// AggregateChanges returns the aggregate's events.
	AggregateChanges() []Change

	// AggregateVersion returns the aggregate's version.
	AggregateVersion() Version

	// ApplyChange applies the given change (event) to the aggregate, updating its state accordingly.
	ApplyChange(Change)
}
