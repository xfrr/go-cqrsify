package message

import (
	"context"
	"fmt"
	"reflect"
	"sync"
)

var _ Bus = (*InMemoryBus)(nil)

type handlerFn func(ctx context.Context, msg Message) (any, error)

// InMemoryBus is an in-memory implementation of the message bus.
type InMemoryBus struct {
	mu              sync.RWMutex
	handlers        map[string]Handler[Message, any] // Store original handlers
	wrappedHandlers map[string]handlerFn             // Store wrapped handlers
	middlewares     []HandlerMiddleware
}

func NewInMemoryBus() *InMemoryBus {
	return &InMemoryBus{
		handlers:        make(map[string]Handler[Message, any]),
		wrappedHandlers: make(map[string]handlerFn),
		middlewares:     make([]HandlerMiddleware, 0),
	}
}

func (b *InMemoryBus) Dispatch(ctx context.Context, msg Message) (res any, err error) {
	t := reflect.TypeOf(msg).Name()
	b.mu.RLock()
	handler, ok := b.wrappedHandlers[t]
	b.mu.RUnlock()
	if !ok {
		return nil, fmt.Errorf("no handler registered for message type %s", t)
	}

	return handler(ctx, msg)
}

func (b *InMemoryBus) RegisterHandler(msgType string, handler Handler[Message, any]) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	if _, exists := b.handlers[msgType]; exists {
		return fmt.Errorf("handler already registered for message type %s", msgType)
	}

	// Store the original handler
	b.handlers[msgType] = handler

	// Create the wrapped handler with current middleware chain
	wrapped := b.applyMiddlewares(handler)
	b.wrappedHandlers[msgType] = func(ctx context.Context, msg Message) (any, error) {
		return wrapped.Handle(ctx, msg)
	}

	return nil
}

func (b *InMemoryBus) Use(middleware HandlerMiddleware) {
	b.mu.Lock()
	defer b.mu.Unlock()

	// Add the middleware to the chain
	b.middlewares = append(b.middlewares, middleware)

	// Re-wrap all existing handlers with the new middleware chain
	for msgType, originalHandler := range b.handlers {
		wrapped := b.applyMiddlewares(originalHandler)
		b.wrappedHandlers[msgType] = func(ctx context.Context, msg Message) (any, error) {
			return wrapped.Handle(ctx, msg)
		}
	}
}

func (b *InMemoryBus) applyMiddlewares(handler Handler[Message, any]) Handler[Message, any] {
	// To achieve "last added, first executed" behavior:
	// If middlewares = [middleware1, middleware2]
	// We want: middleware2 → middleware1 → handler
	// So we apply them in forward order, which makes the last one the outermost wrapper
	wrapped := handler
	for _, middleware := range b.middlewares {
		wrapped = middleware(wrapped)
	}
	return wrapped
}
