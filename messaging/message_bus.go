package messaging

import "context"

type MessageBus interface {
	MessagePublisher
	MessageSubscriber
}

// MessagePublisher is an interface for publishing messages to an message bus.
type MessagePublisher interface {
	// Publish emits one or more messages. Implementations should provide at-least-once delivery semantics
	// unless otherwise documented.
	Publish(ctx context.Context, msg ...Message) error
}

// MessageSubscriber is an interface for subscribing to messages from an message bus.
type MessageSubscriber interface {
	// Subscribe registers a handler for a given logical message name.
	// It returns an unsubscribe function that can be called to remove the subscription.
	Subscribe(ctx context.Context, subject string, h MessageHandler[Message]) (UnsubscribeFunc, error)
}

type UnsubscribeFunc func()
