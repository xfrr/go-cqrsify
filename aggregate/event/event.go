package event

import "time"

// ID represents an event unique identifier.
type ID comparable

// AggregateRef represents a reference to an aggregate.
type AggregateRef[ID comparable] struct {
	ID      ID
	Name    string
	Version int
}

// Event represents an event with the given payload type.
type Event[ID comparable, Payload any] interface {
	// ID returns the event's ID.
	ID() ID

	// Payload returns the event's payload.
	Payload() Payload

	// Name returns the event's name.
	Name() string

	// Time returns the event's time.
	OccurredAt() time.Time

	// Aggregate returns the event's aggregate reference.
	Aggregate() *AggregateRef[any]
}

// Cast attempts to cast the given event to the given payload type.
func Cast[OutID comparable, OutPayload any, InputID comparable, InputPayload any](
	evt Event[InputID, InputPayload],
) (*Base[OutID, OutPayload], bool) {
	id, ok := any(evt.ID()).(OutID)
	if !ok {
		return nil, false
	}

	payload, ok := any(evt.Payload()).(OutPayload)
	if !ok {
		return nil, false
	}

	return &Base[OutID, OutPayload]{
		id:           id,
		payload:      payload,
		name:         evt.Name(),
		occurredAt:   evt.OccurredAt(),
		aggregateRef: evt.Aggregate(),
	}, true
}
