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

func NewInMemoryQueryBus(optFns ...MessageBusConfigModifier) *InMemoryQueryBus {
	return &InMemoryQueryBus{
		bus: NewInMemoryMessageBus(optFns...),
	}
}

func (b *InMemoryQueryBus) DispatchAndWaitReply(ctx context.Context, query Query) (Message, error) {
	if err := b.bus.Publish(ctx, query); err != nil {
		return nil, fmt.Errorf("query_bus: failed to dispatch query: %w", err)
	}

	// Ensure the query supports replies
	replayable, ok := query.(ReplyableMessage)
	if !ok {
		return nil, fmt.Errorf("query_bus: query does not support replies: %T", query)
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
