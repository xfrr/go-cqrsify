package messaging

import (
	"context"
	"fmt"
)

var _ EventBus = (*InMemoryEventBus)(nil)

// InMemoryEventBus is an in-memory implementation of EventBus.
type InMemoryEventBus struct {
	bus *InMemoryMessageBus
}

func NewInMemoryEventBus(optFns ...MessageBusConfigConfiger) *InMemoryEventBus {
	return &InMemoryEventBus{
		bus: NewInMemoryMessageBus(optFns...),
	}
}

func (b *InMemoryEventBus) Publish(ctx context.Context, events ...Event) error {
	msgs := make([]Message, len(events))
	for i, e := range events {
		msgs[i] = e
	}
	return b.bus.Publish(ctx, msgs...)
}

func (b *InMemoryEventBus) Subscribe(ctx context.Context, h EventHandler[Event]) (UnsubscribeFunc, error) {
	return b.bus.Subscribe(ctx, MessageHandlerFn[Message](func(ctx context.Context, msg Message) error {
		evt, ok := msg.(Event)
		if !ok {
			return InvalidMessageTypeError{Expected: fmt.Sprintf("%T", evt), Actual: fmt.Sprintf("%T", msg)}
		}

		return h.Handle(ctx, evt)
	}))
}

func (b *InMemoryEventBus) Use(mws ...MessageHandlerMiddleware) {
	b.bus.Use(mws...)
}

func (b *InMemoryEventBus) Close() error {
	return b.bus.Close()
}
