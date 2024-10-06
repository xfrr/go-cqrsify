package cqrs

import (
	"context"
)

type Middleware func(func(context.Context, interface{}) (interface{}, error)) func(context.Context, interface{}) (interface{}, error)

// RecoverMiddleware is a middleware that recovers from panics during request handling.
func RecoverMiddleware(hook func(interface{})) Middleware {
	return func(next func(context.Context, interface{}) (interface{}, error)) func(context.Context, interface{}) (interface{}, error) {
		return func(ctx context.Context, cmd interface{}) (interface{}, error) {
			defer func() {
				if r := recover(); r != nil {
					hook(r)
				}
			}()

			return next(ctx, cmd)
		}
	}
}

// ChainMiddleware chains multiple middlewares into a single middleware.
func ChainMiddleware(middlewares ...Middleware) Middleware {
	return func(next func(context.Context, interface{}) (interface{}, error)) func(context.Context, interface{}) (interface{}, error) {
		return func(ctx context.Context, cmd interface{}) (interface{}, error) {
			dispatch := next
			for i := len(middlewares) - 1; i >= 0; i-- {
				dispatch = middlewares[i](dispatch)
			}

			return dispatch(ctx, cmd)
		}
	}
}
