package cqrs

import (
	"context"
	"sync"
)

const HeaderKey = "cqrs_header"

type bus struct {
	handlers    sync.Map
	middlewares []Middleware
}

func NewBus() *bus {
	return &bus{
		handlers: sync.Map{},
	}
}

func (b *bus) Close() {
	b.handlers.Range(func(key, value interface{}) bool {
		b.handlers.Delete(key)
		return true
	})
}

func (b *bus) Exists(name string) bool {
	_, ok := b.handlers.Load(name)
	return ok
}

// Use adds a middleware to all handlers registered in the bus.
func (b *bus) Use(middleware func(func(context.Context, interface{}) (interface{}, error)) func(context.Context, interface{}) (interface{}, error)) {
	b.middlewares = append(b.middlewares, middleware)
}

// Dispatch dispatches a request to the bus.
func (b *bus) Dispatch(
	ctx context.Context,
	name string,
	request interface{},
	opts ...DispatchOption,
) (interface{}, error) {
	handler, ok := b.handlers.Load(name)
	if !ok {
		return nil, ErrHandlerNotFound
	}

	for _, opt := range opts {
		ctx = opt(ctx, request)
	}

	h := handler.(HandlerFuncAny)
	return h(ctx, request)
}

func (b *bus) RegisterHandler(
	ctx context.Context,
	name string,
	handler HandlerFuncAny,
) error {
	// Check if handler already exists.
	if _, ok := b.handlers.Load(name); ok {
		return ErrHandlerAlreadyRegistered
	}

	// apply middlewares to handler
	for _, middleware := range b.middlewares {
		handler = middleware(handler)
	}

	// Register handler.
	b.handlers.Store(name, handler)
	return nil
}

func (b *bus) UnregisterHandler(ctx context.Context, name string) {
	b.handlers.Delete(name)
}
