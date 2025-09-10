package domain

// EventApplier applies events to an Aggregate.
type EventApplier interface {
	// ApplyEvent applies the given event to the aggregate, updating its state accordingly.
	ApplyEvent(Event)
}
