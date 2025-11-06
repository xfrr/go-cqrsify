package messaging

import (
	"context"
	"fmt"
)

// SubscribeEvent is a shorthand for handling events.
func SubscribeEvent[E Event](
	ctx context.Context,
	consumer EventConsumer,
	handlerFn MessageHandler[E],
) (UnsubscribeFunc, error) {
	var zero E
	return consumer.Subscribe(
		ctx,
		MessageHandlerFn[Event](func(ctx context.Context, evt Event) error {
			castEvent, ok := evt.(E)
			if !ok {
				return InvalidMessageTypeError{
					Actual:   fmt.Sprintf("%T", evt),
					Expected: fmt.Sprintf("%T", zero),
				}
			}
			return handlerFn.Handle(ctx, castEvent)
		}),
	)
}
