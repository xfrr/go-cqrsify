package messaging

import (
	"context"
	"fmt"
)

// SubscribeCommand is a shorthand for handling commands.
func SubscribeCommand[E Command](
	ctx context.Context,
	subscriber CommandSubscriber,
	commandType string,
	handler CommandHandler[E],
) (func(), error) {
	return subscriber.Subscribe(
		ctx,
		commandType,
		MessageHandlerFn[Command](func(ctx context.Context, evt Command) error {
			castCommand, ok := evt.(E)
			if !ok {
				return InvalidMessageTypeError{
					Actual:   commandType,
					Expected: fmt.Sprintf("%T", evt),
				}
			}
			return handler.Handle(ctx, castCommand)
		}),
	)
}
