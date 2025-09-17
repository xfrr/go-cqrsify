package messagingnats

import (
	"context"
	"fmt"
	"sync"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/xfrr/go-cqrsify/messaging"
)

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

// Subscribe implements messaging.MessageBus.
func (p *JetstreamMessageBus) Subscribe(ctx context.Context, subject string, handler messaging.MessageHandler[messaging.Message]) (messaging.UnsubscribeFunc, error) {
	consumerCfg := jetstream.ConsumerConfig{
		Durable:       subject + "_durable",
		DeliverPolicy: jetstream.DeliverAllPolicy,
		AckPolicy:     jetstream.AckExplicitPolicy,
	}

	consumer, err := p.js.CreateConsumer(ctx, p.streamName, consumerCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create consumer: %w", err)
	}

	sub, err := consumer.Consume(func(msg jetstream.Msg) {
		m, err := p.deserializer.Deserialize(msg.Data())
		if err != nil {
			p.errorHandler(nil, fmt.Errorf("failed to deserialize message: %w", err))
			return
		}

		if err := handler.Handle(ctx, m); err != nil {
			p.errorHandler(m, fmt.Errorf("failed to handle message: %w", err))
			// TODO: check if its temporary or permanent error to decide ack/nack
			return
		}

		if err := msg.Ack(); err != nil {
			p.errorHandler(m, fmt.Errorf("failed to ack message: %w", err))
			return
		}
	})
	if err != nil {
		return nil, fmt.Errorf("failed to subscribe to subject %s: %w", subject, err)
	}

	unsubscribe := func() {
		sub.Stop()
		p.mu.Lock()
		defer p.mu.Unlock()

		// Remove the handler from the map
		// If there are no more handlers for the subject, delete the entry
		delete(p.handlers, subject)
	}

	// Store the handler
	p.mu.Lock()
	defer p.mu.Unlock()
	p.handlers[subject] = append(p.handlers[subject], handler)

	return unsubscribe, nil
}
