package messagingnats

import (
	"context"

	"github.com/xfrr/go-cqrsify/messaging"
)

// Ensure PubSubMessageBus implements the MessageBus interface.
var _ messaging.CommandBus = (*PubSubCommandBus)(nil)
var _ messaging.CommandBusReplier = (*PubSubCommandBus)(nil)
var _ messaging.CommandConsumerReplier = (*PubSubCommandBus)(nil)

// PubSubMessageBus is a NATS-based implementation of the MessageBus interface.
// It provides methods for publishing and subscribing to messages using NATS as the underlying message bus.
type PubSubCommandBus struct {
	PubSubMessageBus
}

func NewPubSubCommandBus(
	publisher *PubSubMessagePublisher,
	consumer *PubSubMessageConsumer,
) PubSubCommandBus {
	return PubSubCommandBus{
		PubSubMessageBus: NewPubSubMessageBus(publisher, consumer),
	}
}

func (p PubSubCommandBus) Dispatch(ctx context.Context, commands ...messaging.Command) error {
	msgs := make([]messaging.Message, len(commands))
	for i, e := range commands {
		msgs[i] = e
	}
	return p.Publish(ctx, msgs...)
}

func (p PubSubCommandBus) DispatchRequest(ctx context.Context, command messaging.Command) (messaging.Message, error) {
	reply, err := p.PubSubMessageBus.PublishRequest(ctx, command)
	if err != nil {
		return nil, err
	}
	return reply, nil
}

func (p PubSubCommandBus) Subscribe(ctx context.Context, handler messaging.MessageHandler[messaging.Command]) (messaging.UnsubscribeFunc, error) {
	wrappedHandler := messaging.MessageHandlerFn[messaging.Message](func(ctx context.Context, msg messaging.Message) error {
		command, ok := msg.(messaging.Command)
		if !ok {
			return messaging.ErrMessageIsNotCommand
		}
		return handler.Handle(ctx, command)
	})
	return p.PubSubMessageBus.Subscribe(ctx, wrappedHandler)
}

func (p PubSubCommandBus) SubscribeWithReply(
	ctx context.Context,
	handler messaging.MessageHandlerWithReply[messaging.Command, messaging.CommandReply],
) (messaging.UnsubscribeFunc, error) {
	wrappedHandler := messaging.MessageHandlerWithReplyFn[messaging.Message, messaging.MessageReply](func(ctx context.Context, msg messaging.Message) (messaging.MessageReply, error) {
		command, ok := msg.(messaging.Command)
		if !ok {
			return nil, messaging.ErrMessageIsNotCommand
		}
		return handler.Handle(ctx, command)
	})
	return p.PubSubMessageBus.SubscribeWithReply(ctx, wrappedHandler)
}
