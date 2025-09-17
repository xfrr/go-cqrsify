package messaging

import (
	"context"
	"fmt"
)

var _ QueryBus = (*InMemoryQueryBus)(nil)

type QueryHandler[Q Query] = MessageHandler[Q]
type QueryHandlerFn[Q Query] = MessageHandlerFn[Q]

type QueryBus interface {
	QueryDispatcher
	QuerySubscriber
}

// QueryDispatcher is an interface for dispatching querys to a query bus.
type QueryDispatcher interface {
	// DispatchAndWaitReply sends a query and waits for a reply.
	DispatchAndWaitReply(ctx context.Context, qry Query) (Message, error)
}

// QuerySubscriber is an interface for subscribing to querys from a query bus.
type QuerySubscriber interface {
	// Subscribe registers a handler for a given logical query name.
	Subscribe(ctx context.Context, subject string, h QueryHandler[Query]) (UnsubscribeFunc, error)
}

// InMemoryQueryBus is an in-memory implementation of QueryBus.
type InMemoryQueryBus struct {
	bus *InMemoryMessageBus
}

func NewInMemoryQueryBus(optFns ...MessageBusConfigModifier) *InMemoryQueryBus {
	return &InMemoryQueryBus{
		bus: NewInMemoryMessageBus(optFns...),
	}
}

func (b *InMemoryQueryBus) DispatchAndWaitReply(ctx context.Context, qry Query) (Message, error) {
	// Dispatch the query
	if err := b.bus.Publish(ctx, qry); err != nil {
		return nil, fmt.Errorf("query_bus: failed to dispatch query: %w", err)
	}

	// Ensure the query supports replies
	replayable, ok := qry.(ReplyableMessage)
	if !ok {
		return nil, fmt.Errorf("query_bus: query does not support replies: %T", qry)
	}

	// Get the reply message
	replyMsg, err := replayable.GetReply(ctx)
	if err != nil {
		return nil, fmt.Errorf("query_bus: failed to get reply: %w", err)
	}

	return replyMsg, nil
}

func (b *InMemoryQueryBus) Subscribe(ctx context.Context, queryName string, h QueryHandler[Query]) (UnsubscribeFunc, error) {
	return b.bus.Subscribe(ctx, queryName, MessageHandlerFn[Message](func(ctx context.Context, msg Message) error {
		cmd, ok := msg.(Query)
		if !ok {
			return InvalidMessageTypeError{Expected: fmt.Sprintf("%T", cmd), Actual: fmt.Sprintf("%T", msg)}
		}
		return h.Handle(ctx, cmd)
	}))
}

func (b *InMemoryQueryBus) Use(mws ...MessageHandlerMiddleware) {
	b.bus.Use(mws...)
}

func (b *InMemoryQueryBus) Close() error {
	return b.bus.Close()
}
