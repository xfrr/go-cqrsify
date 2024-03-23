package aggregate

import (
	"github.com/xfrr/cqrsify/event"
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
func NextChange[ID comparable, ChangePayload any](a Aggregate[ID], id, name string, payload ChangePayload) {
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
