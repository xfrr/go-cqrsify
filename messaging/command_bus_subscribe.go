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
	var cmde E
	return subscriber.Subscribe(
		ctx,
		commandType,
		CommandHandlerFn[Command](func(ctx context.Context, evt Command) error {
			castCommand, ok := evt.(E)
			if !ok {
				return InvalidMessageTypeError{
					Actual:   fmt.Sprintf("%T", evt),
					Expected: fmt.Sprintf("%T", cmde),
				}
			}
			return handler.Handle(ctx, castCommand)
		}),
	)
}
