package event

import "context"

// Context represents a event lifecycle context.
type Context[ID comparable, P any] interface {
	context.Context

	Event() Event[ID, P]
}

// BaseContext is the internal implementation of Event Context.
// It called BaseContext to avoid name collision with context.Context.
type BaseContext[ID comparable, P any] struct {
	context.Context

	evt Event[ID, P]
}

// Event returns the underlying event.
func (c BaseContext[ID, P]) Event() Event[ID, P] {
	return c.evt
}

// Any returns the underlying event as any.
func (c BaseContext[ID, P]) Any() Context[any, any] {
	evt := Base[any, any]{
		id:           c.evt.ID(),
		payload:      c.evt.Payload(),
		name:         c.evt.Name(),
		occurredAt:   c.evt.OccurredAt(),
		aggregateRef: c.evt.Aggregate(),
	}
	return WithContext(c, evt)
}

// WithContext returns a context that carries a event.
func WithContext[ID comparable, Payload any](base context.Context, evt Event[ID, Payload]) *BaseContext[ID, Payload] {
	return &BaseContext[ID, Payload]{
		Context: base,
		evt:     evt,
	}
}

// CastContext returns a context that carries a event with given types.
func CastContext[OutID comparable, OutPayload any, InputID comparable, InputPayload any](
	ctx Context[InputID, InputPayload],
) (*BaseContext[OutID, OutPayload], bool) {
	// cast the underlying event
	evt, ok := Cast[OutID, OutPayload](ctx.Event())
	if !ok {
		return nil, false
	}

	return &BaseContext[OutID, OutPayload]{
		Context: ctx,
		evt:     evt,
	}, true
}
