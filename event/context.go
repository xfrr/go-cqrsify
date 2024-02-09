package event

import "context"

// anyContext is a event context with any payload.
type anyContext = Context[any]

// Context represents a event lifecycle context.
type Context[P any] interface {
	context.Context

	Event() Event[P]
}

// BaseContext is the internal implementation of Event Context.
// It called BaseContext to avoid name collision with context.Context.
type BaseContext[P any] struct {
	context.Context

	evt Event[P]
}

// Event returns the underlying event.
func (c BaseContext[P]) Event() Event[P] {
	return c.evt
}

// WithContext returns a context that carries a event.
func WithContext[M any](base context.Context, evt Event[M]) Context[M] {
	return BaseContext[M]{
		Context: base,
		evt:     evt,
	}
}

// CastContext returns a context that carries a event with a different message type.
func CastContext[To, From any](ctx Context[From]) (Context[To], bool) {
	evt, ok := Cast[To, From](ctx.Event())
	if !ok {
		return nil, false
	}

	return WithContext[To](ctx, evt), true
}
