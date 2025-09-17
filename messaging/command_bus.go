package messaging

import (
	"context"
	"fmt"
)

var _ CommandBus = (*InMemoryCommandBus)(nil)

type CommandHandler[C Command] = MessageHandler[C]
type CommandHandlerFn[C Command] = MessageHandlerFn[C]

type CommandBus interface {
	CommandDispatcher
	CommandSubscriber
}

// CommandDispatcher is an interface for dispatching commands to a command bus.
type CommandDispatcher interface {
	// Dispatch executes a command. Implementations should provide at-least-once delivery semantics
	// unless otherwise documented.
	Dispatch(ctx context.Context, commands ...Command) error
}

// CommandSubscriber is an interface for subscribing to commands from a command bus.
type CommandSubscriber interface {
	// Subscribe registers a handler for a given logical command name.
	Subscribe(ctx context.Context, subject string, h CommandHandler[Command]) (UnsubscribeFunc, error)
}

// InMemoryCommandBus is an in-memory implementation of CommandBus.
type InMemoryCommandBus struct {
	*InMemoryMessageBus
}

func NewInMemoryCommandBus(optFns ...MessageBusConfigModifier) *InMemoryCommandBus {
	return &InMemoryCommandBus{
		InMemoryMessageBus: NewInMemoryMessageBus(optFns...),
	}
}

func (b *InMemoryCommandBus) Dispatch(ctx context.Context, commands ...Command) error {
	msgs := make([]Message, len(commands))
	for i, e := range commands {
		msgs[i] = e
	}
	return b.InMemoryMessageBus.Publish(ctx, msgs...)
}

func (b *InMemoryCommandBus) Subscribe(ctx context.Context, commandName string, h CommandHandler[Command]) (UnsubscribeFunc, error) {
	return b.InMemoryMessageBus.Subscribe(ctx, commandName, MessageHandlerFn[Message](func(ctx context.Context, msg Message) error {
		cmd, ok := msg.(Command)
		if !ok {
			return InvalidMessageTypeError{Expected: fmt.Sprintf("%T", cmd), Actual: fmt.Sprintf("%T", msg)}
		}
		return h.Handle(ctx, cmd)
	}))
}
