package messaging

import (
	"context"
	"fmt"
)

// QueryBus is an interface for a query bus that can dispatch and consume queries.
//
//go:generate moq -pkg messagingmock -out mock/query_bus.go . QueryBus:QueryBus
type QueryBus interface {
	QueryDispatcher
	QueryConsumer
}

// QueryDispatcher is an interface for dispatching querys to a query bus.
//
//go:generate moq -pkg messagingmock -out mock/query_dispatcher.go . QueryDispatcher:QueryDispatcher
type QueryDispatcher interface {
	// Request sends a query and waits for a reply.
	Request(ctx context.Context, qry Query) (Message, error)
}

// QueryConsumer is an interface for subscribing to querys from a query bus.
//
//go:generate moq -pkg messagingmock -out mock/query_consumer.go . QueryConsumer:QueryConsumer
type QueryConsumer interface {
	// Subscribe registers a handler for a given logical query name.
	Subscribe(ctx context.Context, h MessageHandlerWithReply[Query, QueryReply]) (UnsubscribeFunc, error)
}

// QueryHandlerFn wraps the given QueryHandlerFn into a MessageHandlerWithReplyFn.
func QueryHandlerFn[Q Query, R QueryReply](fn func(ctx context.Context, qry Q) (R, error)) MessageHandlerWithReply[Message, MessageReply] {
	var zeroQry Q
	return MessageHandlerWithReplyFn[Message, MessageReply](func(ctx context.Context, msg Message) (MessageReply, error) {
		castQry, ok := msg.(Q)
		if !ok {
			return nil, InvalidMessageTypeError{
				Actual:   fmt.Sprintf("%T", msg),
				Expected: fmt.Sprintf("%T", zeroQry),
			}
		}
		return fn(ctx, castQry)
	})
}
