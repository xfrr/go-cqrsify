package domain

import "context"

// EventHandler is the unit of work that processes an Event.
// Keep it small and side-effect oriented; pure functions can be wrapped if needed.
type EventHandler interface {
	Handle(ctx context.Context, evt Event) error
}

// EventHandlerFunc makes ordinary functions satisfy Handler.
type EventHandlerFunc func(ctx context.Context, evt Event) error

func (f EventHandlerFunc) Handle(ctx context.Context, evt Event) error { return f(ctx, evt) }

// EventHandlerMiddleware composes around Handler to add cross-cutting concerns (logging, tracing, retries, etc).
// Order matters: Use(A,B) -> A(B(handler)).
type EventHandlerMiddleware func(next EventHandler) EventHandler

// EventBus is the abstraction for publishing and subscribing to domain events.
// Note this interface is infra-agnostic (in-memory, NATS, Kafka, etc. can implement it).
type EventBus interface {
	// Publish emits one or more events. Implementations should provide at-least-once delivery semantics
	// unless otherwise documented.
	Publish(ctx context.Context, evts ...Event) error

	// Subscribe registers a handler for a given logical event name.
	// Returns a function to unsubscribe safely.
	Subscribe(eventName string, h EventHandler) (unsubscribe func())

	// Use installs middlewares applied around handlers for *this* bus instance.
	Use(mw ...EventHandlerMiddleware)

	// Close releases resources (workers, connections). Safe to call multiple times.
	Close() error
}
