package messagingnats

import (
	"context"

	"github.com/nats-io/nats.go"
	"github.com/xfrr/go-cqrsify/messaging"
)

var _ messaging.MessageBus = (*JetStreamMessageBus)(nil)

type JetstreamEventBus struct {
	*JetStreamMessageBus
}

func NewJetstreamEventBus(
	ctx context.Context,
	conn *nats.Conn,
	streamName string,
	serializer messaging.MessageSerializer,
	deserializer messaging.MessageDeserializer,
	opts ...JetStreamMessageBusOption,
) (*JetstreamEventBus, error) {
	jmb, err := NewJetstreamMessageBus(
		ctx,
		conn,
		streamName,
		serializer,
		deserializer,
		opts...,
	)
	if err != nil {
		return nil, err
	}
	return &JetstreamEventBus{
		JetStreamMessageBus: jmb,
	}, nil
}

// Publish implements messaging.MessageBus.
func (p *JetstreamEventBus) Publish(ctx context.Context, events ...messaging.Event) error {
	msgs := make([]messaging.Message, len(events))
	for i, e := range events {
		msgs[i] = e
	}
	return p.JetStreamMessageBus.Publish(ctx, msgs...)
}

// Subscribe implements messaging.MessageBus.
func (p *JetstreamEventBus) Subscribe(ctx context.Context, eventType string, handler messaging.EventHandler[messaging.Event]) (messaging.UnsubscribeFunc, error) {
	wrappedHandler := messaging.MessageHandlerFn[messaging.Message](func(ctx context.Context, msg messaging.Message) error {
		event, ok := msg.(messaging.Event)
		if !ok {
			return messaging.ErrMessageIsNotEvent
		}
		return handler.Handle(ctx, event)
	})
	return p.JetStreamMessageBus.Subscribe(ctx, eventType, wrappedHandler)
}
