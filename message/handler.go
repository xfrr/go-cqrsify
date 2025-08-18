package message

import (
	"context"
)

// Handler is an interface for handling messages.
type Handler[M Message, R any] interface {
	Handle(ctx context.Context, msg M) (R, error)
}

// HandlerFn is a function that handles a specific message.
type HandlerFn[M Message, R any] func(ctx context.Context, msg M) (R, error)

func (f HandlerFn[M, R]) Handle(ctx context.Context, msg M) (R, error) {
	return f(ctx, msg)
}

type handlerRegistrar interface {
	RegisterHandler(msgType string, handler Handler[Message, any]) error
}

// Handle is a shorthand for handling messages.
func Handle[M Message, R any](b handlerRegistrar, handlerFn HandlerFn[M, R]) error {
	msgName := NameOf[M]()
	return b.RegisterHandler(msgName, HandlerFn[Message, any](func(ctx context.Context, msg Message) (any, error) {
		castMessage, ok := msg.(M)
		if !ok {
			return nil, InvalidMessageTypeError{
				Actual:   msgName,
				Expected: NameOf[M](),
			}
		}
		return handlerFn.Handle(ctx, castMessage)
	}))
}

type HandlerWithResponse[M Message, R any] interface {
	Handle(ctx context.Context, msg M) (R, error)
}

type HandlerWithResponseFn[M Message, R any] func(ctx context.Context, msg M) (R, error)

func (f HandlerWithResponseFn[M, R]) Handle(ctx context.Context, msg M) (R, error) {
	return f(ctx, msg)
}

type handlerWithResponseRegistrar interface {
	RegisterHandler(msgType string, handler HandlerWithResponse[Message, any]) error
}

// HandleWithResponse is a shorthand for handling messages with a response.
func HandleWithResponse[M Message, R any](b handlerWithResponseRegistrar, handlerFn HandlerWithResponseFn[M, R]) error {
	msgName := NameOf[M]()
	return b.RegisterHandler(msgName, HandlerWithResponseFn[Message, any](func(ctx context.Context, msg Message) (any, error) {
		castedMessage, ok := msg.(M)
		if !ok {
			return *new(R), InvalidMessageTypeError{
				Actual:   msgName,
				Expected: NameOf[M](),
			}
		}
		return handlerFn.Handle(ctx, castedMessage)
	}))
}
