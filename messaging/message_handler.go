package messaging

import (
	"context"
)

// MessageHandler is an interface for handling messages.
type MessageHandler[M Message] interface {
	Handle(ctx context.Context, msg M) error
}

// MessageHandlerFn is a function that handles a specific message.
type MessageHandlerFn[M Message] (func(ctx context.Context, msg M) error)

func (f MessageHandlerFn[M]) Handle(ctx context.Context, msg M) error {
	return f(ctx, msg)
}

type MessageHandlerWithResponse[M Message, R any] interface {
	Handle(ctx context.Context, msg M) (R, error)
}

type MessageHandlerWithResponseFn[M Message, R any] func(ctx context.Context, msg M) (R, error)

func (f MessageHandlerWithResponseFn[M, R]) Handle(ctx context.Context, msg M) (R, error) {
	return f(ctx, msg)
}
