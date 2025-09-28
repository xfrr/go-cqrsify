package messagingnats

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/xfrr/go-cqrsify/messaging"
)

// JetStreamMessageConsumer is a consumer that uses NATS JetStream.
type JetStreamMessageConsumer struct {
	mu sync.Mutex

	conn       *nats.Conn
	js         jetstream.JetStream
	streamName string

	subjectBuilder SubjectBuilderFunc
	serializer     messaging.MessageSerializer
	deserializer   messaging.MessageDeserializer
	handlers       map[string][]messaging.MessageHandler[messaging.Message]
	errorHandler   messaging.ErrorHandler
}

func NewJetStreamMessageConsumer(
	_ context.Context,
	conn *nats.Conn,
	streamName string,
	serializer messaging.MessageSerializer,
	deserializer messaging.MessageDeserializer,
	opts ...JetStreamMessageBusOption,
) (*JetStreamMessageConsumer, error) {
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

	p := &JetStreamMessageConsumer{
		mu:             sync.Mutex{},
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

// Subscribe implements messaging.MessageBus.
func (p *JetStreamMessageConsumer) Subscribe(ctx context.Context, msgType string, handler messaging.MessageHandler[messaging.Message]) (messaging.UnsubscribeFunc, error) {
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
		m, deserializeErr := p.deserializer.Deserialize(jmsg.Data())
		if deserializeErr != nil {
			p.errorHandler(nil, fmt.Errorf("failed to deserialize message: %w", deserializeErr))
			return
		}

		if err = handler.Handle(ctx, m); err != nil {
			p.errorHandler(m, fmt.Errorf("failed to handle message: %w", err))
			// TODO: check if its temporary or permanent error to decide ack/nack
			return
		}

		// Check if the message is replayable
		if rmsg, ok := m.(messaging.ReplyableMessage); ok && jmsg.Reply() != "" {
			replyCtx, cancel := context.WithTimeout(ctx, messaging.DefaultReplyTimeoutSeconds*time.Second)
			defer cancel()

			replyMsg, replyErr := rmsg.GetReply(replyCtx)
			if replyErr != nil {
				p.errorHandler(m, fmt.Errorf("failed to get reply message: %w", replyErr))
				return
			}

			replyData, serializeErr := p.serializer.Serialize(replyMsg)
			if serializeErr != nil {
				p.errorHandler(replyMsg, fmt.Errorf("failed to serialize reply message: %w", serializeErr))
				return
			}

			replySubject := jmsg.Headers().Get(replyHeaderKey)
			if err = p.conn.Publish(replySubject, replyData); err != nil {
				p.errorHandler(replyMsg, fmt.Errorf("failed to send reply message: %w", err))
				return
			}
		}

		if err = jmsg.Ack(); err != nil {
			p.errorHandler(m, fmt.Errorf("failed to ack message: %w", err))
			return
		}
	})
	if err != nil {
		return nil, fmt.Errorf("failed to subscribe to subject %s: %w", msgType, err)
	}

	return p.unsubscribeFn(msgType, sub, handler), nil
}

func (p *JetStreamMessageConsumer) unsubscribeFn(
	subject string,
	sub jetstream.ConsumeContext,
	handler messaging.MessageHandler[messaging.Message],
) func() {
	return func() {
		sub.Drain()
		p.mu.Lock()
		defer p.mu.Unlock()

		handlers := p.handlers[subject]
		for i := range handlers {
			if &handlers[i] == &handler {
				handlers = append(handlers[:i], handlers[i+1:]...)
				break
			}
		}

		if len(handlers) == 0 {
			delete(p.handlers, subject)
		}
	}
}
