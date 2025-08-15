package event

import (
	"time"
)

var _ Event[any, any] = (*Base[any, any])(nil)

// Base represents the base implementation of an event.
// It is used to create new events with the given payload type.
type Base[ID comparable, Payload any] struct {
	id           ID
	payload      Payload
	name         string
	occurredAt   time.Time
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

// Name returns the event unique name.
func (e Base[ID, Payload]) Name() string {
	return e.name
}

// OccurredAt returns the event's time.
func (e Base[ID, Payload]) OccurredAt() time.Time {
	return e.occurredAt
}

// AggregateRef returns the event's aggregate reference.
func (e Base[ID, Payload]) Aggregate() *AggregateRef[any] {
	return e.aggregateRef
}

// Any returns the event's payload as an interface.
func (e Base[ID, Payload]) Any() *Base[any, any] {
	return &Base[any, any]{
		id:           e.id,
		payload:      e.payload,
		name:         e.name,
		occurredAt:   e.occurredAt,
		aggregateRef: e.aggregateRef,
	}
}

type ValidationError struct {
	desc string
}

func (e ValidationError) Error() string {
	return e.desc
}

func NewValidationError(desc string) ValidationError {
	return ValidationError{desc: desc}
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

// WithOccurredAt sets the event's time to the given value.
func WithOccurredAt(t time.Time) NewOption {
	return func(e *Base[any, any]) {
		e.occurredAt = t
	}
}

// New creates a new event with the given name and payload.
// Returns an event with the given options applied.
// The event's aggregate reference is set to an empty reference if the options are not provided.
func New[ID comparable, Payload any](
	id ID,
	name string,
	payload Payload,
	opts ...NewOption,
) (*Base[ID, Payload], error) {
	baseEvent := &Base[any, any]{
		id:           id,
		payload:      payload,
		name:         name,
		occurredAt:   time.Now(),
		aggregateRef: nil,
	}

	for _, opt := range opts {
		opt(baseEvent)
	}

	if baseEvent.id == nil {
		return nil, NewValidationError("event ID is nil")
	}

	if baseEvent.name == "" {
		return nil, NewValidationError("event name is empty")
	}

	if baseEvent.payload == nil {
		return nil, NewValidationError("event payload is nil")
	}

	casted, ok := Cast[ID, Payload](baseEvent)
	if !ok {
		return nil, NewValidationError("failed to cast event")
	}

	return casted, nil
}
