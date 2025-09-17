package messagingnats

import (
	"context"

	"github.com/nats-io/nats.go"
	"github.com/xfrr/go-cqrsify/messaging"
)

// Ensure PubSubMessageBus implements the MessageBus interface.
var _ messaging.QueryBus = (*PubSubQueryBus)(nil)

// PubSubMessageBus is a NATS-based implementation of the MessageBus interface.
// It provides methods for publishing and subscribing to messages using NATS as the underlying message bus.
type PubSubQueryBus struct {
	*PubSubMessageBus
}

func (p *PubSubQueryBus) DispatchAndWaitReply(ctx context.Context, query messaging.Query) (messaging.Message, error) {
	res, err := p.PubSubMessageBus.PublishRequest(ctx, query)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (p *PubSubQueryBus) Subscribe(ctx context.Context, subject string, h messaging.QueryHandler[messaging.Query]) (messaging.UnsubscribeFunc, error) {
	wrappedHandler := messaging.MessageHandlerFn[messaging.Message](func(ctx context.Context, msg messaging.Message) error {
		query, ok := msg.(messaging.Query)
		if !ok {
			return messaging.ErrMessageIsNotQuery
		}
		return h.Handle(ctx, query)
	})
	return p.PubSubMessageBus.Subscribe(ctx, subject, wrappedHandler)
}

func NewPubSubQueryBus(
	conn *nats.Conn,
	serializer messaging.MessageSerializer,
	deserializer messaging.MessageDeserializer,
	opts ...PubSubMessageBusOption,
) *PubSubQueryBus {
	return &PubSubQueryBus{
		PubSubMessageBus: NewPubSubMessageBus(conn, serializer, deserializer, opts...),
	}
}
