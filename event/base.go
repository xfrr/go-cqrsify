package event

import "time"

var _ Event[any, any] = (*Base[any, any])(nil)

// Base represents the base implementation of an event.
// It is used to create new events with the given payload type.
type Base[ID comparable, Payload any] struct {
	id           ID
	payload      Payload
	reason       string
	time         time.Time
	aggregateRef *AggregateRef[any]
}

// ID returns the event's ID.
func (e Base[ID, Payload]) ID() ID {
	return e.id
}

// Payload returns the event's payload.
func (e Base[ID, Payload]) Payload() Payload {
	return e.payload
}

// Reason returns the event's reason.
func (e Base[ID, Payload]) Reason() string {
	return e.reason
}

// Time returns the event's time.
func (e Base[ID, Payload]) Time() time.Time {
	return e.time
}

// AggregateRef returns the event's aggregate reference.
func (e Base[ID, Payload]) Aggregate() *AggregateRef[any] {
	return e.aggregateRef
}

// Any returns the event's payload as an interface.
func (e Base[ID, Payload]) Any() *Base[any, any] {
	var aggregateRef *AggregateRef[any]
	if e.aggregateRef != nil {
		aggregateRef = &AggregateRef[any]{
			ID:      e.aggregateRef.ID,
			Name:    e.aggregateRef.Name,
			Version: e.aggregateRef.Version,
		}
	}

	return &Base[any, any]{
		id:           e.id,
		payload:      e.payload,
		reason:       e.reason,
		time:         e.time,
		aggregateRef: aggregateRef,
	}
}

// NewOption represents an option for creating a new event.
type NewOption func(*Base[any, any])

// WithAggregate sets the event's aggregate reference to the given value.
func WithAggregate[ID comparable](id ID, name string, version int) NewOption {
	return func(e *Base[any, any]) {
		e.aggregateRef = &AggregateRef[any]{
			ID:      id,
			Name:    name,
			Version: version,
		}
	}
}

// WithTime sets the event's time to the given value.
func WithTime(t time.Time) NewOption {
	return func(e *Base[any, any]) {
		e.time = t
	}
}

// New creates a new event with the given reason, and payload.
// Returns an event with the given options applied.
// The event's aggregate reference is set to an empty reference if the options are not provided.
func New[ID comparable, Payload any](id ID, reason string, payload Payload, opts ...NewOption) *Base[ID, Payload] {
	e := &Base[any, any]{
		id:      id,
		payload: payload,
		reason:  reason,
		time:    time.Now(),
	}

	var castedPayload Payload
	if e.payload != nil {
		castedPayload = e.payload.(Payload)
	}

	for _, opt := range opts {
		opt(e)
	}

	return &Base[ID, Payload]{
		id:           id,
		payload:      castedPayload,
		reason:       e.reason,
		time:         e.time,
		aggregateRef: e.aggregateRef,
	}
}
