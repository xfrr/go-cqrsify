package messaging

import (
	"context"
)

// UnsubscribeFunc is a function that can be called to unsubscribe a
// consumer from a message bus.
type UnsubscribeFunc func() error

// MessageHandler is an interface for handling messages.
type MessageHandler[M Message] interface {
	Handle(ctx context.Context, msg M) error
}

// MessageHandlerFn is a function that handles a specific message.
type MessageHandlerFn[M Message] (func(ctx context.Context, msg M) error)

func (f MessageHandlerFn[M]) Handle(ctx context.Context, msg M) error {
	return f(ctx, msg)
}

// MessageHandlerWithReply is an interface for handling messages with reply.
type MessageHandlerWithReply[M Message, R MessageReply] interface {
	Handle(ctx context.Context, msg M) (R, error)
}

// MessageHandlerWithReplyFn is a function that handles a specific message with reply.
type MessageHandlerWithReplyFn[M Message, R MessageReply] (func(ctx context.Context, msg M) (R, error))

func (f MessageHandlerWithReplyFn[M, R]) Handle(ctx context.Context, msg M) (R, error) {
	return f(ctx, msg)
}
