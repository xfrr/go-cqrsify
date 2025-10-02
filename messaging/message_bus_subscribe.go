package messaging

import (
	"context"
	"fmt"
)

// SubscribeMessage is a shorthand for handling messages.
func SubscribeMessage[E Message](
	ctx context.Context,
	subscriber MessageSubscriber,
	messageType string,
	handlerFn MessageHandler[E],
) (func(), error) {
	return subscriber.Subscribe(
		ctx,
		messageType,
		MessageHandlerFn[Message](func(ctx context.Context, evt Message) error {
			castMessage, ok := evt.(E)
			if !ok {
				return InvalidMessageTypeError{
					Actual:   messageType,
					Expected: fmt.Sprintf("%T", evt),
				}
			}
			return handlerFn.Handle(ctx, castMessage)
		}),
	)
}
