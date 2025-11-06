package messagingnats

import (
	"context"
	"errors"
	"fmt"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/xfrr/go-cqrsify/messaging"
)

var _ messaging.MessageConsumer = (*JetStreamMessageConsumer)(nil)
var _ messaging.MessageConsumerReplier = (*JetStreamMessageConsumer)(nil)

// JetStreamMessageConsumer is a consumer that uses NATS JetStream.
type JetStreamMessageConsumer struct {
	conn     *nats.Conn
	js       jetstream.JetStream
	consumer jetstream.Consumer

	streamName string
	cfg        JetStreamMessageConsumerConfig

	serializer   messaging.MessageSerializer
	deserializer messaging.MessageDeserializer
	errorHandler messaging.ErrorHandler
}

// NewJetStreamMessageConsumer creates a standard JetStream consumer.
func NewJetStreamMessageConsumer(
	ctx context.Context,
	conn *nats.Conn,
	streamName string,
	serializer messaging.MessageSerializer,
	deserializer messaging.MessageDeserializer,
	opts ...JetStreamMessageConsumerConfiger,
) (*JetStreamMessageConsumer, error) {
	if err := validateInputs(conn, streamName, serializer, deserializer); err != nil {
		return nil, err
	}

	js, err := jetstream.New(conn)
	if err != nil {
		return nil, err
	}

	config := NewJetStreamMessageConsumerConfig(opts...)
	consumer, err := js.CreateOrUpdateConsumer(ctx, streamName, config.ConsumerConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create consumer: %w", err)
	}

	return newConsumer(conn, js, consumer, streamName, serializer, deserializer, config), nil
}

// NewJetStreamOrderedMessageConsumer creates an ordered JetStream consumer.
func NewJetStreamOrderedMessageConsumer(
	ctx context.Context,
	conn *nats.Conn,
	streamName string,
	serializer messaging.MessageSerializer,
	deserializer messaging.MessageDeserializer,
	opts ...JetStreamMessageConsumerConfiger,
) (*JetStreamMessageConsumer, error) {
	if err := validateInputs(conn, streamName, serializer, deserializer); err != nil {
		return nil, err
	}

	js, err := jetstream.New(conn)
	if err != nil {
		return nil, err
	}

	config := NewJetStreamOrderedMessageConsumerConfig(opts...)
	consumer, err := js.OrderedConsumer(ctx, streamName, config.OrderedConsumerConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create ordered consumer: %w", err)
	}

	return newConsumer(conn, js, consumer, streamName, serializer, deserializer, config), nil
}

// Subscribe implements messaging.MessageConsumer.
func (p *JetStreamMessageConsumer) Subscribe(
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
func (p *JetStreamMessageConsumer) SubscribeWithReply(
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

		replyData, serializeErr := p.serializer.Serialize(replyMsg)
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

func (p *JetStreamMessageConsumer) unsubscribeFn(sub jetstream.ConsumeContext) messaging.UnsubscribeFunc {
	return func() error {
		sub.Drain()
		return nil
	}
}

/*** internal helpers ***/

func newConsumer(
	conn *nats.Conn,
	js jetstream.JetStream,
	consumer jetstream.Consumer,
	streamName string,
	serializer messaging.MessageSerializer,
	deserializer messaging.MessageDeserializer,
	cfg JetStreamMessageConsumerConfig,
) *JetStreamMessageConsumer {
	p := &JetStreamMessageConsumer{
		conn:         conn,
		js:           js,
		consumer:     consumer,
		streamName:   streamName,
		serializer:   serializer,
		deserializer: deserializer,
		errorHandler: cfg.ErrorHandler,
		cfg:          cfg,
	}
	if p.errorHandler == nil {
		p.errorHandler = messaging.DefaultErrorHandler
	}
	return p
}

func validateInputs(
	conn *nats.Conn,
	streamName string,
	serializer messaging.MessageSerializer,
	deserializer messaging.MessageDeserializer,
) error {
	if conn == nil {
		return errors.New("nats connection cannot be nil")
	}
	if streamName == "" {
		return errors.New("stream name cannot be empty")
	}
	if serializer == nil {
		return errors.New("serializer cannot be nil")
	}
	if deserializer == nil {
		return errors.New("deserializer cannot be nil")
	}
	return nil
}

func (p *JetStreamMessageConsumer) handleErr(msg messaging.Message, err error) {
	if p.errorHandler != nil {
		p.errorHandler.Handle(msg, err)
	}
}

func (p *JetStreamMessageConsumer) termWithReason(jmsg jetstream.Msg, reason string, msg messaging.Message) {
	if err := jmsg.TermWithReason(reason); err != nil {
		p.handleErr(msg, fmt.Errorf("failed to term message (%s): %w", reason, err))
	}
}

func (p *JetStreamMessageConsumer) deserializeOrTerm(jmsg jetstream.Msg) messaging.Message {
	m, err := p.deserializer.Deserialize(jmsg.Data())
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
