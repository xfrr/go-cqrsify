package messaging

import (
	"context"
	"fmt"
)

// SubscribeQuery is a shorthand for handling querys.
func SubscribeQuery[Q Query, R any](
	ctx context.Context,
	subscriber QuerySubscriber,
	queryType string,
	handler QueryHandler[Q, R],
) (func(), error) {
	return subscriber.Subscribe(
		ctx,
		queryType,
		MessageHandlerWithResponseFn[Query, any](func(ctx context.Context, qry Query) (any, error) {
			req, ok := qry.(Q)
			if !ok {
				var zero R
				return zero, InvalidMessageTypeError{
					Actual:   queryType,
					Expected: fmt.Sprintf("%T", qry),
				}
			}

			res, err := handler.Handle(ctx, req)
			if err != nil {
				var zero R
				return zero, err
			}

			return any(res), nil
		}),
	)
}
