package messagingnats

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/xfrr/go-cqrsify/messaging"
)

const replyHeaderKey = "Nats-Reply-Subject"

// Ensure JetstreamMessageBus implements the MessageBus interface.
var _ messaging.MessageBus = (*JetstreamMessageBus)(nil)

// JetstreamMessageBus is a NATS-based implementation of the MessageBus interface.
// It provides methods for publishing and subscribing to messages using NATS JetStream as the underlying message bus.
type JetstreamMessageBus struct {
	mu sync.Mutex

	conn *nats.Conn
	js   jetstream.JetStream

	streamName     string
	subjectBuilder SubjectBuilder
	serializer     messaging.MessageSerializer
	deserializer   messaging.MessageDeserializer

	handlers     map[string][]messaging.MessageHandler[messaging.Message]
	errorHandler messaging.ErrorHandler
}

func NewJetstreamMessageBus(
	conn *nats.Conn,
	streamName string,
	serializer messaging.MessageSerializer,
	deserializer messaging.MessageDeserializer,
	opts ...MessageBusOption,
) (*JetstreamMessageBus, error) {
	js, err := jetstream.New(conn)
	if err != nil {
		return nil, err
	}

	busOptions := MessageBusOptions{
		subjectBuilder: DefaultSubjectBuilder,
		errorHandler:   messaging.DefaultErrorHandler,
	}
	for _, opt := range opts {
		opt(&busOptions)
	}

	p := &JetstreamMessageBus{
		conn:           conn,
		js:             js,
		streamName:     streamName,
		serializer:     serializer,
		deserializer:   deserializer,
		subjectBuilder: busOptions.subjectBuilder,
		errorHandler:   busOptions.errorHandler,
		handlers:       make(map[string][]messaging.MessageHandler[messaging.Message]),
	}

	return p, nil
}

// Publish implements messaging.MessageBus.
func (p *JetstreamMessageBus) Publish(ctx context.Context, msg ...messaging.Message) error {
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
func (p *JetstreamMessageBus) PublishRequest(ctx context.Context, msg messaging.Message) (messaging.Message, error) {
	// Create a temporary subscription to receive the reply
	msgSubject := p.subjectBuilder(msg)
	replySubject := fmt.Sprintf("%s.reply", msgSubject)

	// Publish the message with a header indicating the reply subject
	data, err := p.serializer.Serialize(msg)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize message: %w", err)
	}

	jsMsg := &nats.Msg{
		Subject: msgSubject,
		Reply:   replySubject,
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
	if err := replyMsg.Ack(); err != nil {
		return nil, fmt.Errorf("failed to ack reply message: %w", err)
	}

	return reply, nil
}

// Subscribe implements messaging.MessageBus.
func (p *JetstreamMessageBus) Subscribe(ctx context.Context, msgType string, handler messaging.MessageHandler[messaging.Message]) (messaging.UnsubscribeFunc, error) {
	consumerCfg := jetstream.ConsumerConfig{
		Durable:       consumerNameFromMessageType(msgType) + "_durable",
		DeliverPolicy: jetstream.DeliverAllPolicy,
		AckPolicy:     jetstream.AckExplicitPolicy,
		FilterSubject: msgType,
	}

	consumer, err := p.js.CreateOrUpdateConsumer(ctx, p.streamName, consumerCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create consumer: %w", err)
	}

	sub, err := consumer.Consume(func(jmsg jetstream.Msg) {
		m, err := p.deserializer.Deserialize(jmsg.Data())
		if err != nil {
			p.errorHandler(nil, fmt.Errorf("failed to deserialize message: %w", err))
			return
		}

		if err := handler.Handle(ctx, m); err != nil {
			p.errorHandler(m, fmt.Errorf("failed to handle message: %w", err))
			// TODO: check if its temporary or permanent error to decide ack/nack
			return
		}

		// Check if the message is replayable
		if rmsg, ok := m.(messaging.ReplyableMessage); ok && jmsg.Reply() != "" {
			replyCtx, cancel := context.WithTimeout(ctx, messaging.DefaultReplyTimeoutSeconds*time.Second)
			defer cancel()

			replyMsg, err := rmsg.GetReply(replyCtx)
			if err != nil {
				p.errorHandler(m, fmt.Errorf("failed to get reply message: %w", err))
				return
			}

			replyData, err := p.serializer.Serialize(replyMsg)
			if err != nil {
				p.errorHandler(replyMsg, fmt.Errorf("failed to serialize reply message: %w", err))
				return
			}

			replySubject := jmsg.Headers().Get(replyHeaderKey)
			if err := p.conn.Publish(replySubject, replyData); err != nil {
				p.errorHandler(replyMsg, fmt.Errorf("failed to send reply message: %w", err))
				return
			}
		}

		if err := jmsg.Ack(); err != nil {
			p.errorHandler(m, fmt.Errorf("failed to ack message: %w", err))
			return
		}
	})
	if err != nil {
		return nil, fmt.Errorf("failed to subscribe to subject %s: %w", msgType, err)
	}

	unsubscribe := func() {
		sub.Stop()
		p.mu.Lock()
		defer p.mu.Unlock()

		// Remove the handler from the map
		// If there are no more handlers for the subject, delete the entry
		delete(p.handlers, msgType)
	}

	// Store the handler
	p.mu.Lock()
	defer p.mu.Unlock()
	p.handlers[msgType] = append(p.handlers[msgType], handler)

	return unsubscribe, nil
}

// consumerNameFromMessageType generates a consumer name based on the message type.
func consumerNameFromMessageType(msgType string) string {
	// normalize the message type to be used as a consumer name
	// replace dots with underscores
	consumerName := strings.ReplaceAll(msgType, ".", "_")
	// truncate to 30 characters to avoid exceeding NATS limits
	if len(consumerName) > 30 {
		consumerName = consumerName[:30]
	}
	return consumerName
}
