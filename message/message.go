package message

import (
	"time"
)

// Message represents a message to be dispatched.
// Message is the common contract for Commands, Events and Queries.
// It carries cross-cutting metadata for observability and routing.
type Message interface {
	// Unique ID for this message instance (idempotency, log correlation).
	MessageID() string

	// CorrelationID links all messages spawned by the same high-level request.
	CorrelationID() string

	// CausationID points to the immediate parent message that caused this one.
	CausationID() string

	// UTC time when the message was created.
	Timestamp() time.Time

	// Arbitrary transport-safe metadata (trace ids, tenant, auth scope, etc).
	Metadata() map[string]string
}
