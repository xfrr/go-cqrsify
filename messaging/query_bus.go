package messaging

import "context"

type QueryHandler[Q Query] = MessageHandler[Q]
type QueryHandlerFn[Q Query] = MessageHandlerFn[Q]

type QueryBus interface {
	QueryDispatcher
	QuerySubscriber
}

// QueryDispatcher is an interface for dispatching querys to a query bus.
type QueryDispatcher interface {
	// DispatchAndWaitReply sends a query and waits for a reply.
	DispatchAndWaitReply(ctx context.Context, qry Query) (Message, error)
}

// QuerySubscriber is an interface for subscribing to querys from a query bus.
type QuerySubscriber interface {
	// Subscribe registers a handler for a given logical query name.
	Subscribe(ctx context.Context, subject string, h QueryHandler[Query]) (UnsubscribeFunc, error)
}
