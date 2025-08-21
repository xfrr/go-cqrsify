package aggregate

import (
	"fmt"
)

type InvalidEventTypeError struct {
	Event    Event
	Expected any
}

func (e InvalidEventTypeError) Error() string {
	return fmt.Sprintf("invalid event type: %T, expected: %T", e.Event, e.Expected)
}

type EventRecorder interface {
	RecordEvent(Event)
}

// NextEvent applies the given event to the aggregate,
// increments the event's version, and appends it to the aggregate's
// uncommitted list of events (if the aggregate implements EventRecorder).
func NextEvent[T comparable](
	agg EventSourcedAggregate[T],
	event Event,
) error {
	if agg == nil {
		return ErrNilAggregate
	}

	agg.ApplyEvent(event)
	if r, ok := agg.(EventRecorder); ok {
		r.RecordEvent(event)
	}

	return nil
}

// HandleEvent registers an event handler for the given event name on the aggregate.
func HandleEvent[T EventSourcedAggregate[string], E any](
	agg T,
	eventName string,
	handler func(agg T, e E) error,
) {
	agg.HandleEvent(eventName, func(event Event) error {
		newEvent, ok := event.(E)
		if !ok {
			return InvalidEventTypeError{Event: event, Expected: newEvent}
		}
		return handler(agg, newEvent)
	})
}
