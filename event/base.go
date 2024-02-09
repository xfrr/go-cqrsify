package event

import "time"

// Base represents the base implementation of an event.
// It is used to create new events with the given payload type.
type Base[Payload any] struct {
	id           ID
	payload      Payload
	reason       string
	time         time.Time
	aggregateRef Aggregate
}

// ID returns the event's ID.
func (e Base[Payload]) ID() ID {
	return e.id
}

// Payload returns the event's payload.
func (e Base[Payload]) Payload() Payload {
	return e.payload
}

// Reason returns the event's reason.
func (e Base[Payload]) Reason() string {
	return e.reason
}

// Time returns the event's time.
func (e Base[Payload]) Time() time.Time {
	return e.time
}

// AggregateRef returns the event's aggregate reference.
func (e Base[Payload]) Aggregate() Aggregate {
	return e.aggregateRef
}

// Any returns the event's payload as an interface.
func (e Base[Payload]) Any() Event[any] {
	n := New[any](e.id.String(), e.reason, e.payload,
		WithAggregate(e.aggregateRef),
	)
	n.time = e.time
	return n
}

// NewOption represents an option for creating a new event.
type NewOption func(*Base[any])

// WithAggregate sets the event's aggregate reference to the given value.
func WithAggregate(ref Aggregate) NewOption {
	return func(e *Base[any]) {
		e.aggregateRef = ref
	}
}

// New creates a new event with the given ID, reason, and payload.
// Returns an event with the given options applied.
// The event's aggregate reference is set to an empty reference if the options are not provided.
func New[Payload any](id string, reason string, payload Payload, opts ...NewOption) *Base[Payload] {
	e := &Base[any]{
		id:      ID(id),
		payload: payload,
		reason:  reason,
		time:    time.Now(),
	}

	for _, opt := range opts {
		opt(e)
	}

	var p Payload
	if _, ok := e.payload.(Payload); ok {
		p = e.payload.(Payload)
	}

	return &Base[Payload]{
		id:           e.id,
		payload:      p,
		reason:       e.reason,
		time:         e.time,
		aggregateRef: e.aggregateRef,
	}
}
