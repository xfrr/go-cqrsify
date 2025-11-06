package messagingnats

import (
	"strings"

	"github.com/xfrr/go-cqrsify/messaging"
)

const replyHeaderKey = "Nats-Reply-Subject"

// Ensure JetstreamMessageBus implements the MessageBus interface.
var _ messaging.MessageBus = (*JetStreamMessageBus)(nil)

// JetStreamMessageBus is a NATS-based implementation of the MessageBus interface.
// It provides methods for publishing and subscribing to messages using NATS JetStream as the underlying message bus.
type JetStreamMessageBus struct {
	*JetstreamMessagePublisher
	*JetStreamMessageConsumer
}

func NewJetstreamMessageBus(
	publisher *JetstreamMessagePublisher,
	consumer *JetStreamMessageConsumer,
) *JetStreamMessageBus {
	return &JetStreamMessageBus{
		JetstreamMessagePublisher: publisher,
		JetStreamMessageConsumer:  consumer,
	}
}

// consumerNameFromMessageType generates a consumer name based on the message type.
func consumerNameFromMessageType(msgType string) string {
	// normalize the message type to be used as a consumer name
	// replace dots with underscores
	consumerName := strings.ReplaceAll(msgType, ".", "_")
	return consumerName
}
