package aggregate

import (
	"fmt"

	"github.com/xfrr/go-cqrsify/aggregate/event"
)

type RaiseEventError struct {
	msg   string
	cause error
}

func (e RaiseEventError) Error() string {
	return fmt.Sprintf("%s: %s", e.msg, e.cause)
}

func (e RaiseEventError) Unwrap() error {
	return e.cause
}

func NewRaiseEventErrorWithCause(msg string, cause error) RaiseEventError {
	return RaiseEventError{
		msg:   msg,
		cause: cause,
	}
}

func NewRaiseEventError(msg string) RaiseEventError {
	return RaiseEventError{
		msg:   msg,
		cause: nil,
	}
}

type EventRecorder interface {
	RecordEvent(Event)
}

// RaiseEvent creates a new event with the given name and payload and applies it to the given Aggregate.
//
// - If the Aggregate implements EventRecorder, the event will be recorded.
//
// - If the Aggregate implements EventCommitter, the event will be committed.
//
// - The event's version will be the next version of the Aggregate.
func RaiseEvent[AID comparable, EID comparable, P any](
	agg Aggregate[AID],
	eventID EID, eventName string, eventPayload P,
) error {
	if agg == nil {
		return NewRaiseEventError("aggregate is nil")
	}

	raisedEvent, err := event.New(
		eventID, eventName, eventPayload,
		event.WithAggregate(
			agg.AggregateID(),
			agg.AggregateName(),
			nextVersion(agg),
		),
	)
	if err != nil {
		return NewRaiseEventErrorWithCause("failed to create event", err)
	}

	agg.ApplyEvent(raisedEvent.Any())

	if r, ok := agg.(EventRecorder); ok {
		r.RecordEvent(raisedEvent.Any())
	}

	return nil
}
