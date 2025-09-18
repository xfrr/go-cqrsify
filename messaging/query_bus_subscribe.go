package messaging

import (
	"context"
	"fmt"
)

// SubscribeQuery is a shorthand for handling querys.
func SubscribeQuery[Q Query](
	ctx context.Context,
	subscriber QuerySubscriber,
	queryType string,
	handler QueryHandler[Q],
) (UnsubscribeFunc, error) {
	return subscriber.Subscribe(
		ctx,
		queryType,
		QueryHandlerFn[Query](func(ctx context.Context, query Query) error {
			q, ok := query.(Q)
			if !ok {
				return InvalidMessageTypeError{
					Expected: fmt.Sprintf("%T", q),
					Actual:   fmt.Sprintf("%T", query),
				}
			}
			return handler.Handle(ctx, q)
		}),
	)
}
