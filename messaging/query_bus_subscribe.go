package messaging

import (
	"context"
	"fmt"
)

// SubscribeQuery is a shorthand for handling querys.
func SubscribeQuery[Q Query, R QueryReply](
	ctx context.Context,
	consumer QueryConsumer,
	handler QueryHandler[Q, R],
) (UnsubscribeFunc, error) {
	return consumer.Subscribe(
		ctx,
		QueryHandlerFn[Query, QueryReply](func(ctx context.Context, query Query) (QueryReply, error) {
			q, ok := query.(Q)
			if !ok {
				return nil, InvalidMessageTypeError{
					Expected: fmt.Sprintf("%T", q),
					Actual:   fmt.Sprintf("%T", query),
				}
			}
			return handler.Handle(ctx, q)
		}),
	)
}
