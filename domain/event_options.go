package domain

import "time"

type EventOption func(*BaseEvent)

// WithEventTimestamp sets the timestamp of the event.
func WithEventTimestamp(t time.Time) EventOption {
	return func(e *BaseEvent) {
		e.timestamp = t
	}
}
