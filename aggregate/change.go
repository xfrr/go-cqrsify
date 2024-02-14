package aggregate

import (
	"errors"

	"github.com/xfrr/cqrsify/event"
)

var (
	// ErrInvalidAggregateID is returned when the change's aggregate ID does not match the aggregate's ID.
	ErrInvalidAggregateID = errors.New("invalid aggregate ID")

	// ErrInvalidAggregateName is returned when the change's aggregate name does not match the aggregate's name.
	ErrInvalidAggregateName = errors.New("invalid aggregate name")

	// ErrInvalidVersion is returned when the change's version does not match the aggregate's version.
	ErrInvalidVersion = errors.New("invalid aggregate version")

	// ErrInvalidChangePayload is returned when the change's payload is nil.
	ErrInvalidChangePayload = errors.New("invalid change payload")

	// ErrInvalidChangeTime is returned when the change's time is zero.
	ErrInvalidChangeTime = errors.New("invalid change time")
)

// ChangeApplier applies changes (events) to an Aggregate.
type ChangeApplier interface {
	// ApplyChange applies the given change (event) to the aggregate, updating its state accordingly.
	ApplyChange(Change)
}

// Change represents an event that changes the state of an Aggregate.
// It is an alias for event.Event[any].
type Change = event.Event[any]

// ValidateChange validates the given change (event) against the Aggregate's state.
// It returns an error if the change is not valid.
func ValidateChange(a Aggregate, change Change) error {
	if change.Aggregate().ID != a.AggregateID().String() {
		return ErrInvalidAggregateID
	}

	if change.Aggregate().Name != a.AggregateName() {
		return ErrInvalidAggregateName
	}

	if change.Payload() == nil {
		return ErrInvalidChangePayload
	}

	if change.Aggregate().Version != NextVersion(a) {
		return ErrInvalidVersion
	}

	return nil
}

// ApplyChange creates a new change (event) and applies it to the given Aggregate.
// It returns the change with the next aggregate version.
// It also records the event if the Aggregate implements the ChangeCommitter interface.
func ApplyChange[ChangePayload any](a Aggregate, changeID, changeName string, changePayload ChangePayload) event.Event[ChangePayload] {
	version := NextVersion(a)

	opts := []event.NewOption{
		event.WithAggregate(event.Aggregate{
			ID:      a.AggregateID().String(),
			Name:    a.AggregateName(),
			Version: version,
		}),
	}

	evt := event.New(changeID, changeName, changePayload, opts...)
	anevt := evt.Any()

	a.ApplyChange(anevt)

	if r, ok := a.(ChangeCommitter); ok {
		r.RecordChange(anevt)
	}

	return evt
}

// UncommittedVersion returns the aggregate version
// based on the uncommitted events.
// If there are no uncommitted events, it returns the current version.
func UncommittedVersion(a Aggregate) int {
	if len(a.AggregateChanges()) == 0 {
		return int(a.AggregateVersion())
	}

	return a.AggregateChanges()[len(a.AggregateChanges())-1].Aggregate().Version
}

// NextVersion returns the next version of the aggregate.
func NextVersion(a Aggregate) int {
	return UncommittedVersion(a) + 1
}
