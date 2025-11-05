package messagingnats

import (
	"context"

	"github.com/xfrr/go-cqrsify/messaging"
)

// Ensure PubSubMessageBus implements the MessageBus interface.
var _ messaging.QueryBus = (*PubSubQueryBus)(nil)

// PubSubMessageBus is a NATS-based implementation of the MessageBus interface.
// It provides methods for publishing and subscribing to messages using NATS as the underlying message bus.
type PubSubQueryBus struct {
	*PubSubMessageBus
}

func NewPubSubQueryBus(
	pubSubPublisher *PubSubMessagePublisher,
	pubSubConsumer *PubSubMessageConsumer,
) *PubSubQueryBus {
	return &PubSubQueryBus{
		PubSubMessageBus: NewPubSubMessageBus(pubSubPublisher, pubSubConsumer),
	}
}

func (p *PubSubQueryBus) Request(ctx context.Context, query messaging.Query) (messaging.Message, error) {
	res, err := p.PublishRequest(ctx, query)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (p *PubSubQueryBus) Subscribe(ctx context.Context, h messaging.QueryHandler[messaging.Query, messaging.QueryReply]) (messaging.UnsubscribeFunc, error) {
	wrappedHandler := messaging.MessageHandlerWithReplyFn[messaging.Message, messaging.QueryReply](func(ctx context.Context, msg messaging.Message) (messaging.QueryReply, error) {
		query, ok := msg.(messaging.Query)
		if !ok {
			return nil, messaging.ErrMessageIsNotQuery
		}
		return h.Handle(ctx, query)
	})
	return p.PubSubMessageBus.SubscribeWithReply(ctx, wrappedHandler)
}
