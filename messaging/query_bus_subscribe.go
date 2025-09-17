package messaging

import (
	"context"
	"fmt"
)

// SubscribeQuery is a shorthand for handling querys.
func SubscribeQuery[E Query](
	ctx context.Context,
	subscriber QuerySubscriber,
	queryType string,
	handler QueryHandler[E],
) (UnsubscribeFunc, error) {
	return subscriber.Subscribe(
		ctx,
		queryType,
		MessageHandlerFn[Query](func(ctx context.Context, evt Query) error {
			castQuery, ok := evt.(E)
			if !ok {
				return InvalidMessageTypeError{
					Actual:   queryType,
					Expected: fmt.Sprintf("%T", evt),
				}
			}
			return handler.Handle(ctx, castQuery)
		}),
	)
}
