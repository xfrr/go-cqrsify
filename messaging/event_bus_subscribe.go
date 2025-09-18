package messaging

import (
	"context"
	"fmt"
)

// SubscribeEvent is a shorthand for handling events.
func SubscribeEvent[E Event](
	ctx context.Context,
	subscriber EventSubscriber,
	eventType string,
	handlerFn EventHandler[E],
) (func(), error) {
	return subscriber.Subscribe(
		ctx,
		eventType,
		EventHandlerFn[Event](func(ctx context.Context, evt Event) error {
			castEvent, ok := evt.(E)
			if !ok {
				return InvalidMessageTypeError{
					Actual:   eventType,
					Expected: fmt.Sprintf("%T", evt),
				}
			}
			return handlerFn.Handle(ctx, castEvent)
		}),
	)
}
