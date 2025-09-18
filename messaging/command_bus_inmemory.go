package messaging

import (
	"context"
	"fmt"
)

var _ CommandBus = (*InMemoryCommandBus)(nil)

// InMemoryCommandBus is an in-memory implementation of CommandBus.
type InMemoryCommandBus struct {
	bus *InMemoryMessageBus
}

func NewInMemoryCommandBus(optFns ...MessageBusConfigModifier) *InMemoryCommandBus {
	return &InMemoryCommandBus{
		bus: NewInMemoryMessageBus(optFns...),
	}
}

func (b *InMemoryCommandBus) Dispatch(ctx context.Context, commands ...Command) error {
	msgs := make([]Message, len(commands))
	for i, e := range commands {
		msgs[i] = e
	}
	return b.bus.Publish(ctx, msgs...)
}

func (b *InMemoryCommandBus) Subscribe(ctx context.Context, commandName string, h CommandHandler[Command]) (UnsubscribeFunc, error) {
	return b.bus.Subscribe(ctx, commandName, MessageHandlerFn[Message](func(ctx context.Context, msg Message) error {
		cmd, ok := msg.(Command)
		if !ok {
			return InvalidMessageTypeError{Expected: fmt.Sprintf("%T", cmd), Actual: fmt.Sprintf("%T", msg)}
		}
		return h.Handle(ctx, cmd)
	}))
}

func (b *InMemoryCommandBus) Use(mws ...MessageHandlerMiddleware) {
	b.bus.Use(mws...)
}

func (b *InMemoryCommandBus) Close() error {
	return b.bus.Close()
}
