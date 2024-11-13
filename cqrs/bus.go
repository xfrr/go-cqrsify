package cqrs

import (
	"context"
	"sync"
)

const HeaderKey = "cqrs_header"

type InMemoryBus struct {
	handlers    sync.Map
	middlewares []Middleware
}

func NewInMemoryBus() *InMemoryBus {
	return &InMemoryBus{
		handlers:    sync.Map{},
		middlewares: []Middleware{},
	}
}

func (b *InMemoryBus) Close() {
	b.handlers.Range(func(key, _ interface{}) bool {
		b.handlers.Delete(key)
		return true
	})
}

func (b *InMemoryBus) Exists(name string) bool {
	_, ok := b.handlers.Load(name)
	return ok
}

// Use adds a middleware to all handlers registered in the bus.
func (b *InMemoryBus) Use(middleware Middleware) {
	b.middlewares = append(b.middlewares, middleware)
}

// Dispatch dispatches a request to the bus.
func (b *InMemoryBus) Dispatch(
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

func (b *InMemoryBus) RegisterHandler(
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

func (b *InMemoryBus) UnregisterHandler(ctx context.Context, name string) {
	b.handlers.Delete(name)
}
