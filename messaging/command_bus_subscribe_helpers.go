package messaging

import (
	"context"
	"fmt"
)

// SubscribeCommand is a shorthand for handling commands.
func SubscribeCommand[E Command](
	ctx context.Context,
	consumer CommandConsumer,
	handler CommandHandler[E],
) (UnsubscribeFunc, error) {
	var zero E
	return consumer.Subscribe(
		ctx,
		CommandHandlerFn[Command](func(ctx context.Context, evt Command) error {
			castCommand, ok := evt.(E)
			if !ok {
				return InvalidMessageTypeError{
					Actual:   fmt.Sprintf("%T", evt),
					Expected: fmt.Sprintf("%T", zero),
				}
			}
			return handler.Handle(ctx, castCommand)
		}),
	)
}

// SubscribeCommandWithReply is a shorthand for handling commands with reply.
func SubscribeCommandWithReply[E Command, R CommandReply](
	ctx context.Context,
	consumer CommandConsumerReplier,
	handler CommandHandlerWithReply[E, R],
) (UnsubscribeFunc, error) {
	var zero E
	return consumer.SubscribeWithReply(
		ctx,
		CommandHandlerWithReplyFn[Command, CommandReply](func(ctx context.Context, cmd Command) (CommandReply, error) {
			castCmd, ok := cmd.(E)
			if !ok {
				return nil, InvalidMessageTypeError{
					Actual:   fmt.Sprintf("%T", cmd),
					Expected: fmt.Sprintf("%T", zero),
				}
			}
			return handler.Handle(ctx, castCmd)
		}),
	)
}
