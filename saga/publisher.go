package saga

import "context"

type MessagePublisher interface {
	// Publish emits one or more messages. Implementations should provide at-least-once delivery semantics
	// unless otherwise documented.
	Publish(ctx context.Context, messages ...Message) error
}

// Message is a generic message interface.
type Message interface {
	// MessageType returns the type of the message.
	MessageType() string
}
