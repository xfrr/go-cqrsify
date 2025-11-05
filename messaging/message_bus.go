package messaging

import "context"

type MessageBus interface {
	MessagePublisher
	MessageConsumer
}

type MessageBusReplier interface {
	MessagePublisherReplier
	MessageConsumerReplier
}

// MessagePublisher is an interface for publishing messages to an message bus.
//
//go:generate moq -pkg messagingmock -out mock/message_publisher.go . MessagePublisher:MessagePublisher
type MessagePublisher interface {
	// Publish emits one or more messages. Implementations should provide at-least-once delivery semantics
	// unless otherwise documented.
	Publish(ctx context.Context, messages ...Message) error
}

// MessagePublisherReplier is an interface for sending messages and waiting for replies.
//
//go:generate moq -pkg messagingmock -out mock/message_publisher_replier.go . MessagePublisherReplier:MessagePublisherReplier
type MessagePublisherReplier interface {
	// PublishRequest sends a message and waits for a reply.
	PublishRequest(ctx context.Context, msg Message) (Message, error)
}

// MessageConsumer is an interface for subscribing to messages from an message bus.
//
//go:generate moq -pkg messagingmock -out mock/message_consumer.go . MessageConsumer:MessageConsumer
type MessageConsumer interface {
	// Subscribe registers a handler for a given logical message name.
	// It returns an unsubscribe function that can be called to remove the subscription.
	Subscribe(ctx context.Context, h MessageHandler[Message]) (UnsubscribeFunc, error)
}

// MessageConsumerReplier is an interface for subscribing to messages with reply from an message bus.
//
//go:generate moq -pkg messagingmock -out mock/message_consumer_replier.go . MessageConsumerReplier:MessageConsumerReplier
type MessageConsumerReplier interface {
	SubscribeWithReply(ctx context.Context, h MessageHandlerWithReply[Message, Message]) (UnsubscribeFunc, error)
}
