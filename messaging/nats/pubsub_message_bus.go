package messagingnats

import (
	"github.com/xfrr/go-cqrsify/messaging"
)

// Ensure PubSubMessageBus implements the MessageBus interface.
var _ messaging.MessageBus = (*PubSubMessageBus)(nil)

// PubSubMessageBus is a NATS-based implementation of the MessageBus interface.
// It provides methods for publishing and subscribing to messages using NATS as the underlying message bus.
type PubSubMessageBus struct {
	*PubSubMessagePublisher
	*PubSubMessageConsumer
}

func NewPubSubMessageBus(
	pubSubPublisher *PubSubMessagePublisher,
	pubSubConsumer *PubSubMessageConsumer,
) PubSubMessageBus {
	return PubSubMessageBus{
		PubSubMessagePublisher: pubSubPublisher,
		PubSubMessageConsumer:  pubSubConsumer,
	}
}
