package messaging

import (
	"context"
	"fmt"
)

// RegisterCommandHandler is a shorthand for handling commands.
func RegisterCommandHandler[E Command](
	ctx context.Context,
	consumer CommandConsumer,
	handler MessageHandler[E],
) (UnsubscribeFunc, error) {
	var zero E
	return consumer.Subscribe(
		ctx,
		MessageHandlerFn[Command](func(ctx context.Context, evt Command) error {
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

// RegisterCommandHandlerWithReply is a shorthand for handling commands with reply.
func RegisterCommandHandlerWithReply[E Command, R CommandReply](
	ctx context.Context,
	consumer CommandConsumerReplier,
	handler MessageHandlerWithReply[E, R],
) (UnsubscribeFunc, error) {
	var zero E
	return consumer.SubscribeWithReply(
		ctx,
		MessageHandlerWithReplyFn[Command, CommandReply](func(ctx context.Context, cmd Command) (CommandReply, error) {
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
