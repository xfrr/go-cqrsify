package messaging

import (
	"context"
	"fmt"
)

type QueryHandler[C Query, R any] = MessageHandlerWithResponse[C, R]
type QueryHandlerFn[C Query, R any] = MessageHandlerWithResponseFn[C, R]

type QueryBus interface {
	QueryDispatcher
	QuerySubscriber
}

// QueryDispatcher is an interface for dispatching querys to a query bus.
type QueryDispatcher interface {
	// Dispatch executes a query. Implementations should provide at-least-once delivery semantics
	// unless otherwise documented.
	Dispatch(ctx context.Context, query Query) error
}

// QuerySubscriber is an interface for subscribing to querys from a query bus.
type QuerySubscriber interface {
	// Subscribe registers a handler for a given logical query name.
	Subscribe(ctx context.Context, subject string, h QueryHandler[Query, any]) (unsubscribe func(), err error)
}

// InMemoryQueryBus is an in-memory implementation of QueryBus.
type InMemoryQueryBus struct {
	*InMemoryMessageBus
}

func NewInMemoryQueryBus(optFns ...MessageBusConfigModifier) *InMemoryQueryBus {
	return &InMemoryQueryBus{
		InMemoryMessageBus: NewInMemoryMessageBus(optFns...),
	}
}

func (b *InMemoryQueryBus) Dispatch(ctx context.Context, qry Query) error {
	return b.Publish(ctx, qry)
}

func (b *InMemoryQueryBus) Subscribe(ctx context.Context, queryName string, h QueryHandler[Query, any]) (func(), error) {
	return b.SubscribeWithReply(ctx, queryName, MessageHandlerWithResponseFn[Message, any](func(ctx context.Context, msg Message) (any, error) {
		qry, ok := msg.(Query)
		if !ok {
			return nil, InvalidMessageTypeError{Expected: fmt.Sprintf("%T", qry), Actual: fmt.Sprintf("%T", msg)}
		}

		return h.Handle(ctx, qry)
	}))
}
