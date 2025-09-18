package messaging

import "context"

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
