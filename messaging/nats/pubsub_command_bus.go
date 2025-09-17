package messagingnats

import (
	"context"

	"github.com/nats-io/nats.go"
	"github.com/xfrr/go-cqrsify/messaging"
)

// Ensure PubSubMessageBus implements the MessageBus interface.
var _ messaging.CommandBus = (*PubSubCommandBus)(nil)

// PubSubMessageBus is a NATS-based implementation of the MessageBus interface.
// It provides methods for publishing and subscribing to messages using NATS as the underlying message bus.
type PubSubCommandBus struct {
	*PubSubMessageBus
}

func NewPubSubCommandBus(
	conn *nats.Conn,
	serializer messaging.MessageSerializer,
	deserializer messaging.MessageDeserializer,
	opts ...PubSubMessageBusOption,
) *PubSubCommandBus {
	return &PubSubCommandBus{
		PubSubMessageBus: NewPubSubMessageBus(conn, serializer, deserializer, opts...),
	}
}

// Publish implements messaging.MessageBus.
func (p *PubSubCommandBus) Dispatch(ctx context.Context, commands ...messaging.Command) error {
	msgs := make([]messaging.Message, len(commands))
	for i, e := range commands {
		msgs[i] = e
	}
	return p.PubSubMessageBus.Publish(ctx, msgs...)
}

// Subscribe implements messaging.MessageBus.
func (p *PubSubCommandBus) Subscribe(ctx context.Context, commandType string, handler messaging.MessageHandler[messaging.Command]) (messaging.UnsubscribeFunc, error) {
	wrappedHandler := messaging.MessageHandlerFn[messaging.Message](func(ctx context.Context, msg messaging.Message) error {
		command, ok := msg.(messaging.Command)
		if !ok {
			return messaging.ErrMessageIsNotCommand
		}
		return handler.Handle(ctx, command)
	})
	return p.PubSubMessageBus.Subscribe(ctx, commandType, wrappedHandler)
}
