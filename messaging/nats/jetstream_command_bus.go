package messagingnats

import (
	"context"

	"github.com/nats-io/nats.go/jetstream"
	"github.com/xfrr/go-cqrsify/messaging"
)

var _ messaging.CommandBus = (*JetStreamCommandBus)(nil)
var _ messaging.CommandBusReplier = (*JetStreamCommandBus)(nil)
var _ messaging.CommandConsumerReplier = (*JetStreamCommandBus)(nil)

type JetStreamCommandBus struct {
	JetStreamMessageBus
}

func NewJetStreamCommandBus(
	publisher *JetstreamMessagePublisher,
	consumer *JetStreamMessageConsumer[jetstream.ConsumerConfig],
) JetStreamCommandBus {
	jmb := NewJetstreamMessageBus(publisher, consumer)
	return JetStreamCommandBus{
		JetStreamMessageBus: jmb,
	}
}

// Dispatch implements messaging.CommandDispatcher.
func (p JetStreamCommandBus) Dispatch(ctx context.Context, commands ...messaging.Command) error {
	msgs := make([]messaging.Message, len(commands))
	for i, e := range commands {
		msgs[i] = e
	}
	return p.Publish(ctx, msgs...)
}

// DispatchRequest implements messaging.CommandBusRequester.
func (p JetStreamCommandBus) DispatchRequest(ctx context.Context, cmd messaging.Command) (messaging.Message, error) {
	return p.PublishRequest(ctx, cmd)
}

// Subscribe implements messaging.CommandConsumer.
func (p JetStreamCommandBus) Subscribe(
	ctx context.Context,
	handler messaging.MessageHandler[messaging.Command],
) (messaging.UnsubscribeFunc, error) {
	wrappedHandler := messaging.MessageHandlerFn[messaging.Message](func(ctx context.Context, msg messaging.Message) error {
		command, ok := msg.(messaging.Command)
		if !ok {
			return messaging.ErrMessageIsNotCommand
		}
		return handler.Handle(ctx, command)
	})
	return p.JetStreamMessageBus.Subscribe(ctx, wrappedHandler)
}

// SubscribeWithReply implements messaging.CommandConsumerWithReply.
func (p JetStreamCommandBus) SubscribeWithReply(
	ctx context.Context,
	handler messaging.MessageHandlerWithReply[messaging.Command, messaging.CommandReply],
) (messaging.UnsubscribeFunc, error) {
	wrappedHandler := messaging.MessageHandlerWithReplyFn[messaging.Message, messaging.MessageReply](func(ctx context.Context, cmd messaging.Message) (messaging.MessageReply, error) {
		command, ok := cmd.(messaging.Command)
		if !ok {
			return nil, messaging.ErrMessageIsNotCommand
		}
		return handler.Handle(ctx, command)
	})
	return p.JetStreamMessageBus.SubscribeWithReply(ctx, wrappedHandler)
}
