package messaging

import (
	"context"
	"fmt"
)

type CommandHandler[C Command] = MessageHandler[C]
type CommandHandlerFn[C Command] = MessageHandlerFn[C]

type CommandHandlerWithReply[E Command, R CommandReply] = MessageHandlerWithReply[E, R]
type CommandHandlerWithReplyFn[E Command, R CommandReply] = MessageHandlerWithReplyFn[E, R]

// NewCommandHandlerFn creates a new CommandHandler from the given function.
func NewCommandHandlerFn[C Command](fn func(ctx context.Context, cmd C) error) MessageHandler[C] {
	return MessageHandlerFn[C](fn)
}

// NewCommandHandlerWithReplyFn wraps the given CommandHandlerWithReplyFn into a MessageHandlerWithReplyFn.
func NewCommandHandlerWithReplyFn[E Command, R CommandReply](fn func(ctx context.Context, cmd E) (R, error)) MessageHandlerWithReply[Message, MessageReply] {
	var zeroCmd E
	return MessageHandlerWithReplyFn[Message, MessageReply](func(ctx context.Context, msg Message) (MessageReply, error) {
		castCmd, ok := msg.(E)
		if !ok {
			return nil, InvalidMessageTypeError{
				Actual:   fmt.Sprintf("%T", msg),
				Expected: fmt.Sprintf("%T", zeroCmd),
			}
		}
		return fn(ctx, castCmd)
	})
}

// CommandBus is an interface for dispatching commands and subscribing to command responses.
//
//go:generate moq -pkg messagingmock -out mock/command_bus.go . CommandBus:CommandBus
type CommandBus interface {
	CommandDispatcher
	CommandConsumer
}

// CommandBusReplier is an interface for dispatching commands and subscribing to command responses with reply.
//
//go:generate moq -pkg messagingmock -out mock/command_bus_replier.go . CommandBusReplier:CommandBusReplier
type CommandBusReplier interface {
	CommandDispatcherReplier
	CommandConsumerReplier
}

// CommandDispatcher is an interface for dispatching commands to a command bus.
//
//go:generate moq -pkg messagingmock -out mock/command_dispatcher.go . CommandDispatcher:CommandDispatcher
type CommandDispatcher interface {
	// Dispatch executes a command. Implementations should provide at-least-once delivery semantics
	// unless otherwise documented.
	Dispatch(ctx context.Context, commands ...Command) error
}

// CommandDispatcherReplier is an interface for dispatching commands and waiting for replies.
//
//go:generate moq -pkg messagingmock -out mock/command_dispatcher_replier.go . CommandDispatcherReplier:CommandDispatcherReplier
type CommandDispatcherReplier interface {
	// PublishRequest sends a command and waits for a reply.
	PublishRequest(ctx context.Context, cmd Command) (Message, error)
}

// CommandConsumer is an interface for subscribing to commands from a command bus.
type CommandConsumer interface {
	// Subscribe registers a handler for a given logical command name.
	Subscribe(ctx context.Context, h CommandHandler[Command]) (UnsubscribeFunc, error)
}

// CommandConsumerReplier is an interface for subscribing to commands with reply from a command bus.
//
//go:generate moq -pkg messagingmock -out mock/command_consumer_replier.go . CommandConsumerReplier:CommandConsumerReplier
type CommandConsumerReplier interface {
	SubscribeWithReply(ctx context.Context, handler CommandHandlerWithReply[Command, CommandReply]) (UnsubscribeFunc, error)
}
