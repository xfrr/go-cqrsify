package cqrs

import "context"

type Header map[string]any

type contextKey string

func WithHeader(key string, value any) DispatchOption {
	return DispatchOption(func(ctx context.Context, request interface{}) context.Context {
		header, ok := HeaderFromContext(ctx)
		if !ok {
			header = Header{}
		}

		header[key] = value
		return context.WithValue(ctx, contextKey(HeaderKey), header)
	})
}

func HeaderFromContext(ctx context.Context) (Header, bool) {
	header, ok := ctx.Value(contextKey(HeaderKey)).(Header)
	return header, ok
}
