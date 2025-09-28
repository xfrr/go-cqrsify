package messagingnats

import (
	"context"
	"fmt"
	"strings"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/xfrr/go-cqrsify/messaging"
)

const replyHeaderKey = "Nats-Reply-Subject"

// Ensure JetstreamMessageBus implements the MessageBus interface.
var _ messaging.MessageBus = (*JetStreamMessageBus)(nil)

// JetStreamMessageBus is a NATS-based implementation of the MessageBus interface.
// It provides methods for publishing and subscribing to messages using NATS JetStream as the underlying message bus.
type JetStreamMessageBus struct {
	JetstreamMessagePublisher
	JetStreamMessageConsumer
}

func NewJetstreamMessageBus(
	ctx context.Context,
	conn *nats.Conn,
	streamName string,
	serializer messaging.MessageSerializer,
	deserializer messaging.MessageDeserializer,
	opts ...JetStreamMessageBusOption,
) (*JetStreamMessageBus, error) {
	js, err := jetstream.New(conn)
	if err != nil {
		return nil, err
	}

	busOptions := JetStreamMessageBusOptions{
		MessageBusOptions: MessageBusOptions{
			subjectBuilder: DefaultSubjectBuilder,
			errorHandler:   messaging.DefaultErrorHandler,
		},
		streamCfg: defaultStreamConfig(streamName),
	}
	for _, opt := range opts {
		opt(&busOptions)
	}

	// Create the stream if it doesn't exist
	_, err = js.CreateOrUpdateStream(ctx, busOptions.streamCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create or update stream: %w", err)
	}

	p := &JetStreamMessageBus{
		JetstreamMessagePublisher: JetstreamMessagePublisher{
			conn:           conn,
			js:             js,
			streamName:     streamName,
			serializer:     serializer,
			subjectBuilder: busOptions.subjectBuilder,
		},
		JetStreamMessageConsumer: JetStreamMessageConsumer{
			conn:           conn,
			js:             js,
			streamName:     streamName,
			serializer:     serializer,
			deserializer:   deserializer,
			subjectBuilder: busOptions.subjectBuilder,
			errorHandler:   busOptions.errorHandler,
			handlers:       make(map[string][]messaging.MessageHandler[messaging.Message]),
		},
	}

	return p, nil
}

// consumerNameFromMessageType generates a consumer name based on the message type.
func consumerNameFromMessageType(msgType string) string {
	// normalize the message type to be used as a consumer name
	// replace dots with underscores
	consumerName := strings.ReplaceAll(msgType, ".", "_")
	// truncate to 30 characters to avoid exceeding NATS limits
	const maxConsumerNameLength = 30
	if len(consumerName) > maxConsumerNameLength {
		consumerName = consumerName[:maxConsumerNameLength]
	}
	return consumerName
}
