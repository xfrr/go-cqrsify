package messaging

import (
	"context"
	"fmt"
)

// RegisterQueryHandler is a shorthand for handling queries of a specific type.
func RegisterQueryHandler[Q Query, R QueryReply](
	ctx context.Context,
	consumer QueryConsumer,
	handler MessageHandlerWithReply[Q, R],
) (UnsubscribeFunc, error) {
	return consumer.Subscribe(
		ctx,
		MessageHandlerWithReplyFn[Query, QueryReply](func(ctx context.Context, query Query) (QueryReply, error) {
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
