package command

import "context"

// anyContext is a command context with any payload.
type anyContext = Context[any]

// Context represents a command lifecycle context.
type Context[P any] interface {
	context.Context

	Command() Command[P]
}

// cmdContext is the internal implementation of Command Context.
// It called cmdContext to avoid name collision with context.Context.
type cmdContext[P any] struct {
	context.Context

	cmd Command[P]
}

// Command returns the underlying command.
func (c cmdContext[P]) Command() Command[P] {
	return c.cmd
}

// NewContext returns a context that carries a command.
func NewContext[M any](base context.Context, cmd Command[M]) Context[M] {
	return cmdContext[M]{
		Context: base,
		cmd:     cmd,
	}
}

// CastContext returns a context that carries a command with a different message type.
func CastContext[To, From any](ctx Context[From]) (Context[To], bool) {
	cmd, ok := Cast[To, From](ctx.Command())
	if !ok {
		return nil, false
	}

	return NewContext[To](ctx, cmd), true
}
