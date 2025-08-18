package message

import "context"

type HandlerMiddleware func(next Handler[Message, any]) Handler[Message, any]

// HandlerPanicRecoveryMiddleware is a middleware that recovers from panics during request handling.
func HandlerPanicRecoveryMiddleware(hook func(any)) HandlerMiddleware {
	return func(next Handler[Message, any]) Handler[Message, any] {
		return HandlerFn[Message, any](func(ctx context.Context, msg Message) (any, error) {
			defer func() {
				if r := recover(); r != nil {
					hook(r)
				}
			}()

			return next.Handle(ctx, msg)
		})
	}
}

// ChainHandlerMiddlewares is a middleware that composes multiple middlewares.
func ChainHandlerMiddlewares(middlewares ...HandlerMiddleware) HandlerMiddleware {
	return func(next Handler[Message, any]) Handler[Message, any] {
		for i := len(middlewares) - 1; i >= 0; i-- {
			next = middlewares[i](next)
		}
		return next
	}
}
