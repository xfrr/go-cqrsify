package cqrs

import (
	"context"
)

type Middleware func(next HandlerFuncAny) HandlerFuncAny

// RecoverMiddleware is a middleware that recovers from panics during request handling.
func RecoverMiddleware(hook func(interface{})) Middleware {
	return func(next HandlerFuncAny) HandlerFuncAny {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			defer func() {
				if r := recover(); r != nil {
					hook(r)
				}
			}()

			return next(ctx, req)
		}
	}
}

// ChainMiddleware chains multiple middlewares into a single middleware.
func ChainMiddleware(middlewares ...Middleware) Middleware {
	return func(next HandlerFuncAny) HandlerFuncAny {
		dispatch := next
		for i := len(middlewares) - 1; i >= 0; i-- {
			dispatch = middlewares[i](dispatch)
		}

		return dispatch
	}
}
