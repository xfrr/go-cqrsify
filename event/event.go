package event

import "time"

// ID represents an event unique identifier.
type ID string

// String returns the identifier as a string.
func (id ID) String() string {
	return string(id)
}

// Aggregate represents a reference to an aggregate.
type Aggregate struct {
	ID      string
	Name    string
	Version int
}

// Event represents an event with the given payload type.
type Event[Payload any] interface {
	// ID returns the event's ID.
	ID() ID

	// Payload returns the event's payload.
	Payload() Payload

	// Reason returns the event's reason.
	Reason() string

	// Time returns the event's time.
	Time() time.Time

	// Aggregate returns the event's aggregate reference.
	Aggregate() Aggregate
}

// Cast attempts to cast the given event to the given payload type.
func Cast[To, From any](evt Event[From]) (Base[To], bool) {
	payload, ok := any(evt.Payload()).(To)
	if !ok {
		return Base[To]{}, false
	}

	return Base[To]{
		id:           evt.ID(),
		payload:      payload,
		reason:       evt.Reason(),
		time:         evt.Time(),
		aggregateRef: evt.Aggregate(),
	}, true
}
