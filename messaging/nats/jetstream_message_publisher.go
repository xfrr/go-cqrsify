package messagingnats

import (
	"context"
	"fmt"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/xfrr/go-cqrsify/messaging"
)

// JetstreamMessagePublisher is a publisher that uses NATS JetStream.
type JetstreamMessagePublisher struct {
	conn       *nats.Conn
	js         jetstream.JetStream
	streamName string

	subjectBuilder SubjectBuilder
	serializer     messaging.MessageSerializer
	deserializer   messaging.MessageDeserializer
}

func NewJetStreamMessagePublisher(
	_ context.Context,
	conn *nats.Conn,
	streamName string,
	serializer messaging.MessageSerializer,
	deserializer messaging.MessageDeserializer,
	opts ...JetStreamMessageBusOption,
) (*JetstreamMessagePublisher, error) {
	js, err := jetstream.New(conn)
	if err != nil {
		return nil, err
	}

	busOptions := JetStreamMessageBusOptions{
		MessageBusOptions: MessageBusOptions{
			subjectBuilder: DefaultSubjectBuilder,
			errorHandler:   messaging.DefaultErrorHandler,
		},
	}
	for _, opt := range opts {
		opt(&busOptions)
	}

	p := &JetstreamMessagePublisher{
		conn:           conn,
		js:             js,
		streamName:     streamName,
		serializer:     serializer,
		deserializer:   deserializer,
		subjectBuilder: busOptions.subjectBuilder,
	}

	return p, nil
}

// Publish implements messaging.MessageBus.
func (p *JetstreamMessagePublisher) Publish(ctx context.Context, msg ...messaging.Message) error {
	for _, m := range msg {
		data, err := p.serializer.Serialize(m)
		if err != nil {
			return err
		}

		opts := []jetstream.PublishOpt{}
		if m.MessageID() != "" {
			opts = append(opts, jetstream.WithMsgID(m.MessageID()))
		}

		subject := p.subjectBuilder(m)
		_, err = p.js.Publish(ctx, subject, data, opts...)
		if err != nil {
			return err
		}
	}

	return nil
}

// PublishRequest sends a request message and waits for a single reply.
func (p *JetstreamMessagePublisher) PublishRequest(ctx context.Context, msg messaging.Message) (messaging.Message, error) {
	msgSubject := p.subjectBuilder(msg)
	replySubject := fmt.Sprintf("%s.reply.%d", msgSubject, time.Now().UnixNano())
	// If the message has a MessageID, use it to create a unique reply subject
	// This helps in correlating replies in case of multiple requests
	// being sent simultaneously
	if msg.MessageID() != "" {
		replySubject = fmt.Sprintf("%s.reply.%s", msgSubject, msg.MessageID())
	}

	// Publish the message with a header indicating the reply subject
	data, err := p.serializer.Serialize(msg)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize message: %w", err)
	}

	jsMsg := &nats.Msg{
		Subject: msgSubject,
		Data:    data,
		Header: nats.Header{
			replyHeaderKey: []string{replySubject},
		},
	}

	pubAck, err := p.js.PublishMsg(ctx, jsMsg)
	if err != nil {
		return nil, fmt.Errorf("failed to publish request message: %w", err)
	}

	// Create a temporary consumer to receive the reply
	consumerCfg := jetstream.ConsumerConfig{
		Name:          consumerNameFromMessageType(msg.MessageType()) + fmt.Sprintf("_reply_%d", pubAck.Sequence),
		DeliverPolicy: jetstream.DeliverAllPolicy,
		AckPolicy:     jetstream.AckExplicitPolicy,
		MaxDeliver:    3,
		FilterSubject: replySubject,
		BackOff:       []time.Duration{time.Second, 2 * time.Second, 5 * time.Second},
	}

	consumer, err := p.js.CreateConsumer(ctx, p.streamName, consumerCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create consumer for reply: %w", err)
	}

	// Receive the reply message
	replyMsg, err := consumer.Next()
	if err != nil {
		return nil, fmt.Errorf("failed to receive reply message: %w", err)
	}

	// Deserialize the reply message
	reply, err := p.deserializer.Deserialize(replyMsg.Data())
	if err != nil {
		termErr := replyMsg.TermWithReason("deserialization_failed")
		if termErr != nil {
			return nil, fmt.Errorf("failed to terminate message after deserialization failure: %w", termErr)
		}
		return nil, fmt.Errorf("failed to deserialize reply message: %w", err)
	}

	// Acknowledge the reply message
	if err = replyMsg.Ack(); err != nil {
		return nil, fmt.Errorf("failed to ack reply message: %w", err)
	}

	return reply, nil
}
