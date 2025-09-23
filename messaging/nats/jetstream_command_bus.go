package messagingnats

import (
	"context"

	"github.com/nats-io/nats.go"
	"github.com/xfrr/go-cqrsify/messaging"
)

var _ messaging.CommandBus = (*JetStreamCommandBus)(nil)

type JetStreamCommandBus struct {
	*JetStreamMessageBus
}

func NewJetStreamCommandBus(
	ctx context.Context,
	conn *nats.Conn,
	streamName string,
	serializer messaging.MessageSerializer,
	deserializer messaging.MessageDeserializer,
	opts ...JetStreamMessageBusOption,
) (*JetStreamCommandBus, error) {
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
	return &JetStreamCommandBus{
		JetStreamMessageBus: jmb,
	}, nil
}

// Publish implements messaging.MessageBus.
func (p *JetStreamCommandBus) Dispatch(ctx context.Context, commands ...messaging.Command) error {
	msgs := make([]messaging.Message, len(commands))
	for i, e := range commands {
		msgs[i] = e
	}
	return p.Publish(ctx, msgs...)
}

// Subscribe implements messaging.MessageBus.
func (p *JetStreamCommandBus) Subscribe(ctx context.Context, commandType string, handler messaging.CommandHandler[messaging.Command]) (messaging.UnsubscribeFunc, error) {
	wrappedHandler := messaging.MessageHandlerFn[messaging.Message](func(ctx context.Context, msg messaging.Message) error {
		command, ok := msg.(messaging.Command)
		if !ok {
			return messaging.ErrMessageIsNotCommand
		}
		return handler.Handle(ctx, command)
	})
	return p.JetStreamMessageBus.Subscribe(ctx, commandType, wrappedHandler)
}
