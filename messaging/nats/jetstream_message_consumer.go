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
	conn       *nats.Conn
	js         jetstream.JetStream
	streamName string
	cfg        JetStreamMessageConsumerConfig

	serializer   messaging.MessageSerializer
	deserializer messaging.MessageDeserializer
	errorHandler messaging.ErrorHandler
}

func NewJetStreamMessageConsumer(
	conn *nats.Conn,
	streamName string,
	serializer messaging.MessageSerializer,
	deserializer messaging.MessageDeserializer,
	opts ...MessageConsumerConfiger,
) (*JetStreamMessageConsumer, error) {
	js, err := jetstream.New(conn)
	if err != nil {
		return nil, err
	}

	config := NewJetStreamMessageConsumerConfig(opts...)

	p := &JetStreamMessageConsumer{
		conn:         conn,
		js:           js,
		streamName:   streamName,
		serializer:   serializer,
		deserializer: deserializer,
		errorHandler: config.ErrorHandler,
		cfg:          config,
	}

	return p, nil
}

// Subscribe implements messaging.MessageConsumer.
func (p *JetStreamMessageConsumer) Subscribe(
	ctx context.Context,
	handler messaging.MessageHandler[messaging.Message],
) (messaging.UnsubscribeFunc, error) {
	consumer, err := p.js.CreateOrUpdateConsumer(ctx, p.streamName, p.cfg.ConsumerConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create consumer: %w", err)
	}

	sub, err := consumer.Consume(func(jmsg jetstream.Msg) {
		m, deserializeErr := p.deserializer.Deserialize(jmsg.Data())
		if deserializeErr != nil {
			p.errorHandler.Handle(nil, fmt.Errorf("failed to deserialize message: %w", deserializeErr))
			return
		}

		if m == nil {
			p.errorHandler.Handle(m, fmt.Errorf("no deserializer found for message"))
			return
		}

		if err = handler.Handle(ctx, m); err != nil {
			p.errorHandler.Handle(m, fmt.Errorf("failed to handle message: %w", err))
			// TODO: check if its temporary or permanent error to decide ack/nack
			nakErr := jmsg.Nak()
			if nakErr != nil {
				p.errorHandler.Handle(m, fmt.Errorf("failed to nack message: %w", nakErr))
			}
			return
		}

		if err = jmsg.Ack(); err != nil {
			p.errorHandler.Handle(m, fmt.Errorf("failed to ack message: %w", err))
			return
		}
	})
	if err != nil {
		return nil, fmt.Errorf("failed to subscribe: %w", err)
	}

	return p.unsubscribeFn(sub), nil
}

// SubscribeWithReply implements messaging.MessageConsumerWithReply.
func (p *JetStreamMessageConsumer) SubscribeWithReply(
	ctx context.Context,
	handler messaging.MessageHandlerWithReply[messaging.Message, messaging.MessageReply],
) (messaging.UnsubscribeFunc, error) {
	consumer, err := p.js.CreateOrUpdateConsumer(ctx, p.streamName, p.cfg.ConsumerConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create consumer: %w", err)
	}

	consumerCtx, err := consumer.Consume(func(jmsg jetstream.Msg) {
		m, deserializeErr := p.deserializer.Deserialize(jmsg.Data())
		if deserializeErr != nil {
			p.errorHandler.Handle(nil, fmt.Errorf("failed to deserialize message: %w", deserializeErr))
			return
		}
		if m == nil {
			p.errorHandler.Handle(nil, errors.New("no deserializer found for message"))
			termErr := jmsg.TermWithReason("no_deserializer")
			if termErr != nil {
				p.errorHandler.Handle(nil, fmt.Errorf("failed to term message: %w", termErr))
			}
			return
		}

		replyMsg, handleErr := handler.Handle(ctx, m)
		if handleErr != nil {
			p.errorHandler.Handle(m, fmt.Errorf("failed to handle message: %w", handleErr))
			// TODO: check if its temporary or permanent error to decide ack/nack
			nakErr := jmsg.Term()
			if nakErr != nil {
				p.errorHandler.Handle(m, fmt.Errorf("failed to nack message: %w", nakErr))
			}
			return
		}
		if replyMsg == nil {
			p.errorHandler.Handle(m, errors.New("handler returned nil reply message"))
			termErr := jmsg.TermWithReason("nil_reply")
			if termErr != nil {
				p.errorHandler.Handle(m, fmt.Errorf("failed to term message after nil reply: %w", termErr))
			}
			return
		}

		replyData, serializeErr := p.serializer.Serialize(replyMsg)
		if serializeErr != nil {
			p.errorHandler.Handle(replyMsg, fmt.Errorf("failed to serialize reply message: %w", serializeErr))
			termErr := jmsg.TermWithReason("serialization_failed")
			if termErr != nil {
				p.errorHandler.Handle(m, fmt.Errorf("failed to term message after serialization failure: %w", termErr))
			}
			return
		}

		replySubject := jmsg.Headers().Get(replyHeaderKey)
		if replySubject == "" {
			p.errorHandler.Handle(replyMsg, errors.New("no reply subject found in message headers"))
			termErr := jmsg.TermWithReason("no_reply_subject")
			if termErr != nil {
				p.errorHandler.Handle(m, fmt.Errorf("failed to term message: %w", termErr))
			}
			return
		}

		if err = p.conn.Publish(replySubject, replyData); err != nil {
			p.errorHandler.Handle(replyMsg, fmt.Errorf("failed to send reply message: %w", err))
			nakErr := jmsg.Nak()
			if nakErr != nil {
				p.errorHandler.Handle(m, fmt.Errorf("failed to nack message: %w", nakErr))
			}
			return
		}

		if err = jmsg.Ack(); err != nil {
			p.errorHandler.Handle(m, fmt.Errorf("failed to ack message: %w", err))
			return
		}
	})

	return p.unsubscribeFn(consumerCtx), nil
}

func (p *JetStreamMessageConsumer) unsubscribeFn(
	sub jetstream.ConsumeContext,
) messaging.UnsubscribeFunc {
	return func() error {
		sub.Drain()
		return nil
	}
}
