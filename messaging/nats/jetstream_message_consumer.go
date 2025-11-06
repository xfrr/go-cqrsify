package messagingnats

import (
	"context"
	"errors"
	"fmt"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/xfrr/go-cqrsify/messaging"
)

var _ messaging.MessageConsumer = (*JetStreamMessageConsumer[jetstream.ConsumerConfig])(nil)
var _ messaging.MessageConsumerReplier = (*JetStreamMessageConsumer[jetstream.ConsumerConfig])(nil)
var _ messaging.MessageConsumer = (*JetStreamMessageConsumer[jetstream.OrderedConsumerConfig])(nil)

// JetStreamMessageConsumer is a consumer that uses NATS JetStream.
type JetStreamMessageConsumer[T jetStreamConsumerConfig] struct {
	conn     *nats.Conn
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
		m := p.deserializeOrTerm(jmsg)
		if m == nil {
			return
		}

		if err := handler.Handle(ctx, m); err != nil {
			p.handleErr(m, fmt.Errorf("failed to handle message: %w", err))
			if nakErr := jmsg.Nak(); nakErr != nil {
				p.handleErr(m, fmt.Errorf("failed to nak message: %w", nakErr))
			}
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
		m := p.deserializeOrTerm(jmsg)
		if m == nil {
			return
		}

		replyMsg, handleErr := handler.Handle(ctx, m)
		if handleErr != nil {
			p.handleErr(m, fmt.Errorf("failed to handle message: %w", handleErr))
			if termErr := jmsg.Term(); termErr != nil {
				p.handleErr(m, fmt.Errorf("failed to term message: %w", termErr))
			}
			return
		}
		if replyMsg == nil {
			p.handleErr(m, errors.New("handler returned nil reply message"))
			p.termWithReason(jmsg, "nil_reply", m)
			return
		}

		replyData, serializeErr := p.cfg.serializer.Serialize(replyMsg)
		if serializeErr != nil {
			p.handleErr(replyMsg, fmt.Errorf("failed to serialize reply message: %w", serializeErr))
			p.termWithReason(jmsg, "serialization_failed", m)
			return
		}

		replySubject := jmsg.Headers().Get(replyHeaderKey)
		if replySubject == "" {
			p.handleErr(replyMsg, errors.New("no reply subject found in message headers"))
			p.termWithReason(jmsg, "no_reply_subject", m)
			return
		}

		if err := p.conn.Publish(replySubject, replyData); err != nil {
			p.handleErr(replyMsg, fmt.Errorf("failed to send reply message: %w", err))
			if nakErr := jmsg.Nak(); nakErr != nil {
				p.handleErr(m, fmt.Errorf("failed to nak message: %w", nakErr))
			}
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

/*** internal helpers ***/

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

func validateInputs(
	js jetstream.JetStream,
	streamName string,
) error {
	if js == nil {
		return errors.New("jetstream instance cannot be nil")
	}
	if streamName == "" {
		return errors.New("stream name cannot be empty")
	}
	return nil
}

func (p *JetStreamMessageConsumer[T]) handleErr(msg messaging.Message, err error) {
	if p.errorHandler != nil {
		p.errorHandler.Handle(msg, err)
	}
}

func (p *JetStreamMessageConsumer[T]) termWithReason(jmsg jetstream.Msg, reason string, msg messaging.Message) {
	if err := jmsg.TermWithReason(reason); err != nil {
		p.handleErr(msg, fmt.Errorf("failed to term message (%s): %w", reason, err))
	}
}

func (p *JetStreamMessageConsumer[T]) deserializeOrTerm(jmsg jetstream.Msg) messaging.Message {
	m, err := p.cfg.deserializer.Deserialize(jmsg.Data())
	if err != nil {
		p.handleErr(nil, fmt.Errorf("failed to deserialize message: %w", err))
		p.termWithReason(jmsg, "deserialization_failed", nil)
		return nil
	}
	if m == nil {
		p.handleErr(nil, errors.New("nil message after deserialization"))
		p.termWithReason(jmsg, "deserialization_failed", nil)
		return nil
	}
	return m
}
