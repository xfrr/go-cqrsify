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
	Subscribe(ctx context.Context, subject string, h QueryHandler[Query]) (unsubscribe func(), err error)
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

func (b *InMemoryQueryBus) DispatchAndWaitReply(ctx context.Context, qry Query) (Message, error) {
	// 	1. subscribe a one-time handler for the reply
	replyCh := make(chan Message, 1)

	// wrap the query to include the reply subject
	wrappedQuery := baseMessageWithReply{
		Message: qry,
		replyCh: replyCh,
	}

	// 2. dispatch the query
	if err := b.Publish(ctx, wrappedQuery); err != nil {
		return nil, fmt.Errorf("query_bus: failed to dispatch query: %w", err)
	}

	// 3. wait for the reply or context cancellation
	select {
	case reply := <-replyCh:
		return reply, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func (b *InMemoryQueryBus) Subscribe(ctx context.Context, queryName string, h QueryHandler[Query]) (func(), error) {
	return b.InMemoryMessageBus.Subscribe(ctx, queryName, MessageHandlerFn[Message](func(ctx context.Context, msg Message) error {
		cmd, ok := msg.(Query)
		if !ok {
			return InvalidMessageTypeError{Expected: fmt.Sprintf("%T", cmd), Actual: fmt.Sprintf("%T", msg)}
		}
		return h.Handle(ctx, cmd)
	}))
}

type baseMessageWithReply struct {
	Message
	replyCh chan<- Message
}

func (m baseMessageWithReply) Reply(ctx context.Context, response Message) error {
	select {
	case m.replyCh <- response:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
