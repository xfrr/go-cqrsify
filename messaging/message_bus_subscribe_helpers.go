package messaging

import (
	"context"
	"fmt"
)

// SubscribeMessage is a shorthand for handling messages.
func SubscribeMessage[E Message](
	ctx context.Context,
	consumer MessageConsumer,
	handlerFn MessageHandler[E],
) (UnsubscribeFunc, error) {
	return consumer.Subscribe(
		ctx,
		MessageHandlerFn[Message](func(ctx context.Context, evt Message) error {
			castMessage, ok := evt.(E)
			if !ok {
				return InvalidMessageTypeError{
					Actual:   fmt.Sprintf("%T", evt),
					Expected: fmt.Sprintf("%T", *new(E)),
				}
			}
			return handlerFn.Handle(ctx, castMessage)
		}),
	)
}
