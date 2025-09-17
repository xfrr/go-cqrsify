package messagingnats

import (
	"context"

	"github.com/nats-io/nats.go"
	"github.com/xfrr/go-cqrsify/messaging"
)

// Ensure JetstreamMessageBus implements the MessageBus interface.
var _ messaging.QueryBus = (*JetstreamQueryBus)(nil)

// JetstreamMessageBus is a NATS-based implementation of the MessageBus interface.
// It provides methods for publishing and subscribing to messages using NATS as the underlying message bus.
type JetstreamQueryBus struct {
	*JetstreamMessageBus
}

func (p *JetstreamQueryBus) DispatchAndWaitReply(ctx context.Context, query messaging.Query) (messaging.Message, error) {
	res, err := p.JetstreamMessageBus.PublishRequest(ctx, query)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (p *JetstreamQueryBus) Subscribe(ctx context.Context, subject string, h messaging.QueryHandler[messaging.Query]) (unsubscribe func(), err error) {
	wrappedHandler := messaging.MessageHandlerFn[messaging.Message](func(ctx context.Context, msg messaging.Message) error {
		query, ok := msg.(messaging.Query)
		if !ok {
			return messaging.ErrMessageIsNotQuery
		}
		return h.Handle(ctx, query)
	})
	return p.JetstreamMessageBus.Subscribe(ctx, subject, wrappedHandler)
}

func NewJetstreamQueryBus(
	ctx context.Context,
	conn *nats.Conn,
	streamName string,
	serializer messaging.MessageSerializer,
	deserializer messaging.MessageDeserializer,
	opts ...JetStreamMessageBusOption,
) (*JetstreamQueryBus, error) {
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

	return &JetstreamQueryBus{
		JetstreamMessageBus: jmb,
	}, nil
}
