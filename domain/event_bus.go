package domain

import (
	"context"
	"fmt"
)

type EventTypeMismatchError struct {
	Expected any
	Actual   any
}

func (e EventTypeMismatchError) Error() string {
	return fmt.Sprintf("eventbus: event type mismatch; expected %T, got %T", e.Expected, e.Actual)
}

// EventHandler is the unit of work that processes an Event.
// Keep it small and side-effect oriented; pure functions can be wrapped if needed.
type EventHandler[T Event] interface {
	Handle(ctx context.Context, evt T) error
}

// EventHandlerFunc makes ordinary functions satisfy Handler.
type EventHandlerFunc[T Event] func(ctx context.Context, evt T) error

func (f EventHandlerFunc[T]) Handle(ctx context.Context, evt T) error { return f(ctx, evt) }

// EventHandlerMiddleware composes around Handler to add cross-cutting concerns (logging, tracing, retries, etc).
type EventHandlerMiddleware[T Event] func(next EventHandler[T]) EventHandler[T]

// EventBus is the abstraction for publishing and subscribing to domain events.
// Note this interface is infra-agnostic (in-memory, NATS, Kafka, etc. can implement it).
type EventBus interface {
	// Publish emits one or more events. Implementations should provide at-least-once delivery semantics
	// unless otherwise documented.
	Publish(ctx context.Context, evts ...Event) error

	// Subscribe registers a handler for a given logical event name.
	// Returns a function to unsubscribe safely.
	Subscribe(eventName string, h EventHandler[Event]) (unsubscribe func())

	// Use installs middlewares applied around handlers for *this* bus instance.
	Use(mw ...EventHandlerMiddleware[Event])

	// Close releases resources (workers, connections). Safe to call multiple times.
	Close() error
}

// SubscribeEvent is a helper to register typed event handlers with type assertion at runtime.
func SubscribeEvent[T Event](bus EventBus, eventName string, h EventHandler[T]) func() {
	wrapped := EventHandlerFunc[Event](func(ctx context.Context, evt Event) error {
		e, ok := evt.(T)
		if !ok {
			return EventTypeMismatchError{Expected: new(T), Actual: evt}
		}
		return h.Handle(ctx, e)
	})
	return bus.Subscribe(eventName, wrapped)
}

// HandleFunc is a helper to register ordinary functions as event handlers.
func HandleFunc[T Event](bus EventBus, eventName string, fn func(ctx context.Context, evt T) error) (unsubscribe func()) {
	return SubscribeEvent(bus, eventName, EventHandlerFunc[T](fn))
}
