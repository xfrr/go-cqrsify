package messaging

import (
	"context"
	"fmt"
)

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
	// DispatchRequest sends a command and waits for a reply.
	DispatchRequest(ctx context.Context, cmd Command) (Message, error)
}

// CommandConsumer is an interface for subscribing to commands from a command bus.
type CommandConsumer interface {
	// Subscribe registers a handler for a given logical command name.
	Subscribe(ctx context.Context, h MessageHandler[Command]) (UnsubscribeFunc, error)
}

// CommandConsumerReplier is an interface for subscribing to commands with reply from a command bus.
//
//go:generate moq -pkg messagingmock -out mock/command_consumer_replier.go . CommandConsumerReplier:CommandConsumerReplier
type CommandConsumerReplier interface {
	SubscribeWithReply(ctx context.Context, handler MessageHandlerWithReply[Command, CommandReply]) (UnsubscribeFunc, error)
}

// NewCommandHandlerFn creates a new CommandHandler from the given function.
func CommandHandlerFn[C Command](fn func(ctx context.Context, cmd C) error) MessageHandler[Message] {
	var zeroCmd C
	return MessageHandlerFn[Message](func(ctx context.Context, msg Message) error {
		castCmd, ok := msg.(C)
		if !ok {
			return InvalidMessageTypeError{
				Actual:   fmt.Sprintf("%T", msg),
				Expected: fmt.Sprintf("%T", zeroCmd),
			}
		}
		return fn(ctx, castCmd)
	})
}

// NewCommandHandlerWithReplyFn wraps the given CommandHandlerWithReplyFn into a MessageHandlerWithReplyFn.
func CommandHandlerWithReplyFn[E Command, R CommandReply](fn func(ctx context.Context, cmd E) (R, error)) MessageHandlerWithReply[Message, MessageReply] {
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
