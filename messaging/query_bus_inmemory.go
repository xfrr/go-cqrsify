package messaging

import (
	"context"
	"fmt"
)

var _ QueryBus = (*InMemoryQueryBus)(nil)

// InMemoryQueryBus is an in-memory implementation of QueryBus.
type InMemoryQueryBus struct {
	bus *InMemoryMessageBus
}

func NewInMemoryQueryBus(optFns ...MessageBusConfigConfiger) *InMemoryQueryBus {
	return &InMemoryQueryBus{
		bus: NewInMemoryMessageBus(optFns...),
	}
}

func (b *InMemoryQueryBus) Request(ctx context.Context, query Query) (Message, error) {
	return b.bus.PublishRequest(ctx, query)
}

func (b *InMemoryQueryBus) Subscribe(ctx context.Context, h MessageHandlerWithReply[Query, QueryReply]) (UnsubscribeFunc, error) {
	return b.bus.SubscribeWithReply(ctx, MessageHandlerWithReplyFn[Message, MessageReply](func(ctx context.Context, msg Message) (MessageReply, error) {
		q, ok := msg.(Query)
		if !ok {
			return nil, InvalidMessageTypeError{
				Expected: fmt.Sprintf("%T", q),
				Actual:   fmt.Sprintf("%T", msg),
			}
		}
		return h.Handle(ctx, q)
	}))
}

func (b *InMemoryQueryBus) Use(mws ...MessageHandlerMiddleware) {
	b.bus.Use(mws...)
}

func (b *InMemoryQueryBus) Close() error {
	return b.bus.Close()
}
