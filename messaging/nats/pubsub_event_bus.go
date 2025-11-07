package messagingnats

import (
	"context"

	"github.com/xfrr/go-cqrsify/messaging"
)

// Ensure PubSubMessageBus implements the MessageBus interface.
var _ messaging.MessageBus = (*PubSubMessageBus)(nil)

// PubSubMessageBus is a NATS-based implementation of the MessageBus interface.
// It provides methods for publishing and subscribing to messages using NATS as the underlying message bus.
type PubSubEventBus struct {
	PubSubMessageBus
}

func NewPubSubEventBus(
	pubSubPublisher *PubSubMessagePublisher,
	pubSubConsumer *PubSubMessageConsumer,
) PubSubEventBus {
	return PubSubEventBus{
		PubSubMessageBus: NewPubSubMessageBus(pubSubPublisher, pubSubConsumer),
	}
}

// Publish implements messaging.MessageBus.
func (p PubSubEventBus) Publish(ctx context.Context, events ...messaging.Event) error {
	msgs := make([]messaging.Message, len(events))
	for i, e := range events {
		msgs[i] = e
	}
	return p.PubSubMessageBus.Publish(ctx, msgs...)
}

// Subscribe implements messaging.MessageBus.
func (p PubSubEventBus) Subscribe(ctx context.Context, handler messaging.MessageHandler[messaging.Event]) (messaging.UnsubscribeFunc, error) {
	wrappedHandler := messaging.MessageHandlerFn[messaging.Message](func(ctx context.Context, msg messaging.Message) error {
		event, ok := msg.(messaging.Event)
		if !ok {
			return messaging.ErrMessageIsNotEvent
		}
		return handler.Handle(ctx, event)
	})
	return p.PubSubMessageBus.Subscribe(ctx, wrappedHandler)
}
