package messagingnats

import (
	"context"
	"errors"
	"fmt"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/xfrr/go-cqrsify/messaging"
	"go.opentelemetry.io/otel/propagation"
)

var _ messaging.MessageConsumer = (*JetStreamMessageConsumer[jetstream.ConsumerConfig])(nil)
var _ messaging.MessageConsumerReplier = (*JetStreamMessageConsumer[jetstream.ConsumerConfig])(nil)
var _ messaging.MessageConsumer = (*JetStreamMessageConsumer[jetstream.OrderedConsumerConfig])(nil)

// JetStreamMessageConsumer is a consumer that uses NATS JetStream.
type JetStreamMessageConsumer[T jetStreamConsumerConfig] struct {
	js       jetstream.JetStream
	consumer jetstream.Consumer

	streamName   string
	cfg          JetStreamMessageConsumerConfig[T]
	errorHandler messaging.ErrorHandler
}

// NewJetStreamMessageConsumer creates a standard JetStream consumer.
func NewJetStreamMessageConsumer(
	ctx context.Context,
	js jetstream.JetStream,
	streamName string,
	opts ...JetStreamMessageConsumerConfiger[jetstream.ConsumerConfig],
) (*JetStreamMessageConsumer[jetstream.ConsumerConfig], error) {
	if err := validateInputs(js, streamName); err != nil {
		return nil, err
	}

	config := NewJetStreamMessageConsumerConfig(opts...)
	consumer, err := js.CreateOrUpdateConsumer(ctx, streamName, config.ConsumerConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create consumer: %w", err)
	}

	return newConsumer(js, consumer, streamName, config), nil
}

// NewJetStreamOrderedMessageConsumer creates an ordered JetStream consumer.
func NewJetStreamOrderedMessageConsumer(
	ctx context.Context,
	js jetstream.JetStream,
	streamName string,
	opts ...JetStreamMessageConsumerConfiger[jetstream.OrderedConsumerConfig],
) (*JetStreamMessageConsumer[jetstream.OrderedConsumerConfig], error) {
	if err := validateInputs(js, streamName); err != nil {
		return nil, err
	}

	config := NewJetStreamOrderedMessageConsumerConfig(opts...)
	consumer, err := js.OrderedConsumer(ctx, streamName, config.ConsumerConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create ordered consumer: %w", err)
	}

	return newConsumer(js, consumer, streamName, config), nil
}

// Subscribe implements messaging.MessageConsumer.
func (p *JetStreamMessageConsumer[T]) Subscribe(
	ctx context.Context,
	handler messaging.MessageHandler[messaging.Message],
) (messaging.UnsubscribeFunc, error) {
	if handler == nil {
		return nil, errors.New("handler cannot be nil")
	}
	if p.consumer == nil {
		return nil, errors.New("consumer is not initialized")
	}

	cc, err := p.consumer.Consume(func(jmsg jetstream.Msg) {
		m := p.deserializeMessage(jmsg)
		if m == nil {
			return
		}

		// Extract tracing context from message headers
		// and create a new context for handling the message
		msgCtx := p.propagateTracingContext(ctx, jmsg)

		if err := handler.Handle(msgCtx, m); err != nil {
			// TODO: decide whether to Nak or Term based on error type
			p.errAndNak(jmsg, m, fmt.Errorf("failed to handle message: %w", err))
			return
		}

		if err := jmsg.Ack(); err != nil {
			p.handleErr(m, fmt.Errorf("failed to ack message: %w", err))
			return
		}
	})
	if err != nil {
		return nil, fmt.Errorf("failed to subscribe: %w", err)
	}

	return p.unsubscribeFn(cc), nil
}

// SubscribeWithReply implements messaging.MessageConsumerWithReply.
func (p *JetStreamMessageConsumer[T]) SubscribeWithReply(
	ctx context.Context,
	handler messaging.MessageHandlerWithReply[messaging.Message, messaging.MessageReply],
) (messaging.UnsubscribeFunc, error) {
	if handler == nil {
		return nil, errors.New("handler cannot be nil")
	}
	if p.consumer == nil {
		return nil, errors.New("consumer is not initialized")
	}

	cc, err := p.consumer.Consume(func(jmsg jetstream.Msg) {
		m := p.deserializeMessage(jmsg)
		if m == nil {
			return
		}

		// Extract tracing context from message headers
		// and create a new context for handling the message
		msgCtx := p.propagateTracingContext(ctx, jmsg)

		replyMsg, handleErr := handler.Handle(msgCtx, m)
		if handleErr != nil {
			// TODO: decide whether to Term or Nak based on error type
			p.errAndTerm(jmsg, m, "message_handling_failed", fmt.Errorf("failed to handle message: %w", handleErr))
			return
		}
		if replyMsg == nil {
			p.errAndTerm(jmsg, m, "nil_reply", errors.New("handler returned nil reply"))
			return
		}

		replyData, serializeErr := p.cfg.Serializer.Serialize(replyMsg)
		if serializeErr != nil {
			p.errAndTerm(jmsg, m, "reply_serialization_failed", fmt.Errorf("failed to serialize reply message: %w", serializeErr))
			return
		}

		replySubject := jmsg.Headers().Get(replyHeaderKey)
		if replySubject == "" {
			p.errAndTerm(jmsg, m, "missing_reply_subject", errors.New("missing reply subject in message headers"))
			return
		}

		headers := nats.Header{}
		// Inject tracing context into reply message headers
		p.cfg.OTELPropagator.Inject(msgCtx, propagation.HeaderCarrier(headers))

		// Create the reply message
		natsReplyMsg := &nats.Msg{
			Subject: replySubject,
			Data:    replyData,
			Header:  headers,
		}

		if _, err := p.js.PublishMsg(msgCtx, natsReplyMsg); err != nil {
			p.errAndTerm(jmsg, m, "reply_publication_failed", fmt.Errorf("failed to publish reply message: %w", err))
			return
		}

		if err := jmsg.Ack(); err != nil {
			p.handleErr(m, fmt.Errorf("failed to ack message: %w", err))
			return
		}
	})
	if err != nil {
		return nil, fmt.Errorf("failed to subscribe with reply: %w", err)
	}

	return p.unsubscribeFn(cc), nil
}

func (p *JetStreamMessageConsumer[T]) unsubscribeFn(sub jetstream.ConsumeContext) messaging.UnsubscribeFunc {
	return func() error {
		sub.Drain()
		return nil
	}
}

func (p *JetStreamMessageConsumer[T]) deserializeMessage(jmsg jetstream.Msg) messaging.Message {
	m, err := p.cfg.Deserializer.Deserialize(jmsg.Data())
	if err != nil {
		p.errAndTerm(
			jmsg,
			nil,
			"deserialization_failed",
			fmt.Errorf("failed to deserialize message: %w", err),
		)
		return nil
	}
	if m == nil {
		p.errAndTerm(
			jmsg,
			nil,
			"nil_message",
			errors.New("nil message after deserialization"),
		)
		return nil
	}
	return m
}

func (p *JetStreamMessageConsumer[T]) propagateTracingContext(parent context.Context, jmsg jetstream.Msg) context.Context {
	return p.cfg.OTELPropagator.Extract(parent, propagation.HeaderCarrier(jmsg.Headers()))
}

func (p *JetStreamMessageConsumer[T]) errAndNak(jmsg jetstream.Msg, msg messaging.Message, err error) {
	p.handleErr(msg, err)
	if nakErr := jmsg.Nak(); nakErr != nil {
		p.handleErr(msg, fmt.Errorf("failed to nak message: %w", nakErr))
	}
}

func (p *JetStreamMessageConsumer[T]) errAndTerm(jmsg jetstream.Msg, msg messaging.Message, reason string, err error) {
	p.handleErr(msg, err)
	if termErr := jmsg.TermWithReason(reason); termErr != nil {
		p.handleErr(msg, fmt.Errorf("failed to term message: %w", termErr))
	}
}

func (p *JetStreamMessageConsumer[T]) handleErr(msg messaging.Message, err error) {
	if p.errorHandler != nil {
		p.errorHandler.Handle(msg, err)
	}
}

func newConsumer[T jetStreamConsumerConfig](
	js jetstream.JetStream,
	consumer jetstream.Consumer,
	streamName string,
	cfg JetStreamMessageConsumerConfig[T],
) *JetStreamMessageConsumer[T] {
	p := &JetStreamMessageConsumer[T]{
		js:           js,
		consumer:     consumer,
		streamName:   streamName,
		errorHandler: cfg.ErrorHandler,
		cfg:          cfg,
	}
	if p.errorHandler == nil {
		p.errorHandler = messaging.DefaultErrorHandler
	}
	return p
}

func validateInputs(js jetstream.JetStream, streamName string) error {
	if js == nil {
		return errors.New("jetstream context cannot be nil")
	}
	if streamName == "" {
		return errors.New("stream name cannot be empty")
	}
	return nil
}
