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

// ApplyChange creates and applies the given event to the Aggregate, updating its state
// accordingly and returns the event with the next aggregate version.
// It also records the event if the Aggregate implements the EventRecorder interface.
func ApplyChange[EventPayload any](a Aggregate, evtid, evtname string, evtpayload EventPayload) event.Event[EventPayload] {
	version := NextVersion(a)

	opts := []event.NewOption{
		event.WithAggregate(event.Aggregate{
			ID:      a.AggregateID().String(),
			Name:    a.AggregateName(),
			Version: version,
		}),
	}

	evt := event.New(evtid, evtname, evtpayload, opts...)
	anevt := evt.Any()

	a.ApplyChange(anevt)

	if r, ok := a.(ChangeCommiter); ok {
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
