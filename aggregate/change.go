package aggregate

import (
	"github.com/xfrr/go-cqrsify/event"
)

// Change represents an event that changes the state of an Aggregate.
// It is an alias for event.Event[any].
type Change = event.Event[any, any]

// ChangeApplier applies changes (events) to an Aggregate.
type ChangeApplier interface {
	// ApplyChange applies the given change (event) to the aggregate, updating its state accordingly.
	ApplyChange(Change)
}

// NextChange adds and applies a new change (event) to the given Aggregate.
func NextChange[AID comparable, EID comparable, P any](a Aggregate[AID], id EID, name string, payload P) {
	change := event.New(
		id,
		name,
		payload,
		event.WithAggregate(a.AggregateID(), a.AggregateName(), nextVersion(a)),
	)

	a.ApplyChange(change.Any())

	if r, ok := a.(ChangeCommitter); ok {
		r.RecordChange(change.Any())
	}
}
