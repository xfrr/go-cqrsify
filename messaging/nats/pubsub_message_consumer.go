package messagingnats

import (
	"context"
	"fmt"

	"github.com/nats-io/nats.go"
	"github.com/xfrr/go-cqrsify/messaging"
)

var _ messaging.MessageConsumer = (*PubSubMessageConsumer)(nil)

// PubSubMessageConsumer is a consumer that uses NATS JetStream.
type PubSubMessageConsumer struct {
	conn *nats.Conn

	serializer   messaging.MessageSerializer
	deserializer messaging.MessageDeserializer
	cfg          PubSubMessageConsumerConfig
}

func NewPubSubMessageConsumer(
	conn *nats.Conn,
	serializer messaging.MessageSerializer,
	deserializer messaging.MessageDeserializer,
	opts ...PubSubMessageConsumerConfiger,
) (*PubSubMessageConsumer, error) {
	cfg := NewPubSubMessageConsumerConfig(opts...)

	p := &PubSubMessageConsumer{
		conn:         conn,
		serializer:   serializer,
		deserializer: deserializer,
		cfg:          cfg,
	}

	return p, nil
}

// Subscribe implements messaging.MessageBus.
func (p *PubSubMessageConsumer) Subscribe(
	ctx context.Context,
	handler messaging.MessageHandler[messaging.Message],
) (messaging.UnsubscribeFunc, error) {
	sub, err := p.conn.Subscribe(p.cfg.Subject, p.natsMsgHandler(ctx, handler))
	if err != nil {
		return nil, fmt.Errorf("failed to subscribe to subject %s: %w", p.cfg.Subject, err)
	}

	return p.unsubscribeFn(p.cfg.Subject, sub), nil
}

// SubscribeWithReply implements messaging.MessageBus.
func (p *PubSubMessageConsumer) SubscribeWithReply(
	ctx context.Context,
	handler messaging.MessageHandlerWithReply[messaging.Message, messaging.Message],
) (messaging.UnsubscribeFunc, error) {
	sub, err := p.conn.Subscribe(p.cfg.Subject, p.natsMsgHandlerWithReply(ctx, handler))
	if err != nil {
		return nil, fmt.Errorf("failed to subscribe to subject %s: %w", p.cfg.Subject, err)
	}

	return p.unsubscribeFn(p.cfg.Subject, sub), nil
}

func (p *PubSubMessageConsumer) natsMsgHandler(
	ctx context.Context,
	handler messaging.MessageHandler[messaging.Message],
) nats.MsgHandler {
	return func(natsMsg *nats.Msg) {
		m, err := p.deserializer.Deserialize(natsMsg.Data)
		if err != nil {
			p.cfg.ErrorHandler.Handle(nil, fmt.Errorf("failed to deserialize message: %w", err))
			return
		}

		handleErr := handler.Handle(ctx, m)
		if handleErr != nil {
			p.cfg.ErrorHandler.Handle(m, fmt.Errorf("failed to handle message: %w", handleErr))
			return
		}
	}
}

func (p *PubSubMessageConsumer) natsMsgHandlerWithReply(
	ctx context.Context,
	handler messaging.MessageHandlerWithReply[messaging.Message, messaging.Message],
) nats.MsgHandler {
	return func(natsMsg *nats.Msg) {
		m, err := p.deserializer.Deserialize(natsMsg.Data)
		if err != nil {
			p.cfg.ErrorHandler.Handle(nil, fmt.Errorf("failed to deserialize message: %w", err))
			return
		}

		replyMsg, handleErr := handler.Handle(ctx, m)
		if handleErr != nil {
			p.cfg.ErrorHandler.Handle(m, fmt.Errorf("failed to handle message: %w", handleErr))
			return
		}

		replyData, serializeErr := p.serializer.Serialize(replyMsg)
		if serializeErr != nil {
			p.cfg.ErrorHandler.Handle(replyMsg, fmt.Errorf("failed to serialize reply message: %w", serializeErr))
			return
		}

		err = natsMsg.Respond(replyData)
		if err != nil {
			p.cfg.ErrorHandler.Handle(replyMsg, fmt.Errorf("failed to send reply message: %w", err))
			return
		}
	}
}

func (p *PubSubMessageConsumer) unsubscribeFn(
	subject string,
	sub *nats.Subscription,
) messaging.UnsubscribeFunc {
	return messaging.UnsubscribeFunc(func() error {
		err := sub.Drain()
		if err != nil {
			return fmt.Errorf("failed to unsubscribe from subject %s: %w", subject, err)
		}
		return nil
	})
}
