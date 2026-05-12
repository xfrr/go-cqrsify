package messaging

import "context"

// QueryHandlerTypedRouter routes queries to handlers with reply based on query type.
type QueryHandlerTypedRouter struct {
	MessageHandlerWithReplyTypedRouter[Query, QueryReply]
}

// NewQueryHandlerTypedRouter creates a new QueryHandlerTypedRouter.
func NewQueryHandlerTypedRouter() *QueryHandlerTypedRouter {
	return &QueryHandlerTypedRouter{*NewMessageHandlerWithReplyTypedRouter[Query, QueryReply]()}
}

func (r *QueryHandlerTypedRouter) Handle(ctx context.Context, msg Message) (MessageReply, error) {
	qry, ok := msg.(Query)
	if !ok {
		return nil, InvalidMessageTypeError{
			Expected: "Query",
			Actual:   msg.MessageType(),
		}
	}

	msgReply, err := r.MessageHandlerWithReplyTypedRouter.Handle(ctx, qry)
	if err != nil {
		return nil, err
	}

	return msgReply, nil
}
