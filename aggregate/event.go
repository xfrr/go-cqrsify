package aggregate

import (
	"errors"
	"time"
)

var (
	ErrEventNameEmpty = errors.New("event name cannot be empty")
)

// Event represents an event on the event-sourcing context.
// It can be used to represent changes to the state of an aggregate.
type Event interface {
	// Name is a unique and human-readable name for the event.
	Name() string
	// Timestamp returns the time at which the event occurred.
	Timestamp() time.Time
	// AggregateRef returns a reference to the aggregate that the event belongs to.
	AggregateRef() *EventAggregateReference
}

// BaseEvent is the default implementation of the Event interface.
type BaseEvent struct {
	name         string
	timestamp    time.Time
	aggregateRef *EventAggregateReference
}

func (e BaseEvent) Name() string {
	return e.name
}

func (e BaseEvent) Timestamp() time.Time {
	return e.timestamp
}

func (e BaseEvent) AggregateRef() *EventAggregateReference {
	return e.aggregateRef
}

func NewEvent(name string, aggref *EventAggregateReference, opts ...EventOption) BaseEvent {
	event := &BaseEvent{
		name:         name,
		aggregateRef: aggref,
		timestamp:    time.Now(),
	}

	for _, opt := range opts {
		opt(event)
	}

	return *event
}
