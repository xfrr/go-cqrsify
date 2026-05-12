package messaging

import "context"

// CommandHandlerTypedRouter routes commands to handlers based on command type.
type CommandHandlerTypedRouter struct {
	MessageHandlerTypedRouter[Command]
}

// CommandHandlerWithReplyTypedRouter routes commands to handlers with reply based on command type.
type CommandHandlerWithReplyTypedRouter struct {
	MessageHandlerWithReplyTypedRouter[Command, CommandReply]
}

// NewCommandHandlerTypedRouter creates a new CommandHandlerTypedRouter.
func NewCommandHandlerTypedRouter() *CommandHandlerTypedRouter {
	return &CommandHandlerTypedRouter{*NewMessageHandlerTypedRouter[Command]()}
}

// NewCommandHandlerWithReplyTypedRouter creates a new CommandHandlerWithReplyTypedRouter.
func NewCommandHandlerWithReplyTypedRouter() *CommandHandlerWithReplyTypedRouter {
	return &CommandHandlerWithReplyTypedRouter{*NewMessageHandlerWithReplyTypedRouter[Command, CommandReply]()}
}

func (r *CommandHandlerTypedRouter) Handle(ctx context.Context, cmd Message) error {
	command, ok := cmd.(Command)
	if !ok {
		return InvalidMessageTypeError{
			Expected: "Command",
			Actual:   cmd.MessageType(),
		}
	}

	return r.MessageHandlerTypedRouter.Handle(ctx, command)
}

func (r *CommandHandlerWithReplyTypedRouter) Handle(ctx context.Context, msg Message) (MessageReply, error) {
	cmd, ok := msg.(Command)
	if !ok {
		return nil, InvalidMessageTypeError{
			Expected: "Command",
			Actual:   msg.MessageType(),
		}
	}

	msgReply, err := r.MessageHandlerWithReplyTypedRouter.Handle(ctx, cmd)
	if err != nil {
		return nil, err
	}

	return msgReply, nil
}
