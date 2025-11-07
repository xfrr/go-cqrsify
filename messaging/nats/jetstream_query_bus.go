package messagingnats

import (
	"context"

	"github.com/nats-io/nats.go/jetstream"
	"github.com/xfrr/go-cqrsify/messaging"
)

var _ messaging.QueryBus = (*JetstreamQueryBus)(nil)

type JetstreamQueryBus struct {
	JetStreamMessageBus
}

func NewJetstreamQueryBus(
	publisher *JetstreamMessagePublisher,
	consumer *JetStreamMessageConsumer[jetstream.ConsumerConfig],
) JetstreamQueryBus {
	jmb := NewJetstreamMessageBus(publisher, consumer)
	return JetstreamQueryBus{
		JetStreamMessageBus: jmb,
	}
}

func (p JetstreamQueryBus) Request(ctx context.Context, query messaging.Query) (messaging.Message, error) {
	return p.PublishRequest(ctx, query)
}

func (p JetstreamQueryBus) Subscribe(ctx context.Context, h messaging.MessageHandlerWithReply[messaging.Query, messaging.QueryReply]) (messaging.UnsubscribeFunc, error) {
	wrappedHandler := messaging.MessageHandlerWithReplyFn[messaging.Message, messaging.QueryReply](func(ctx context.Context, msg messaging.Message) (messaging.QueryReply, error) {
		query, ok := msg.(messaging.Query)
		if !ok {
			return nil, messaging.ErrMessageIsNotQuery
		}
		return h.Handle(ctx, query)
	})
	return p.JetStreamMessageBus.SubscribeWithReply(ctx, wrappedHandler)
}
