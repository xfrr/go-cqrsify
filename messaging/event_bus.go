package messaging

import (
	"context"
	"fmt"
)

type EventHandler[E Event] = MessageHandler[E]
type EventHandlerFn[E Event] = MessageHandlerFn[E]

// NewEventHandlerFn wraps the given EventHandlerFn into a MessageHandlerFn.
func NewEventHandlerFn[E Event](fn func(ctx context.Context, evt E) error) MessageHandler[Message] {
	var zero E
	return MessageHandlerFn[Message](func(ctx context.Context, msg Message) error {
		castEvent, ok := msg.(E)
		if !ok {
			return InvalidMessageTypeError{
				Actual:   fmt.Sprintf("%T", msg),
				Expected: fmt.Sprintf("%T", zero),
			}
		}
		return fn(ctx, castEvent)
	})
}

// EventBus is an interface for publishing and subscribing to events.
//
//go:generate moq -pkg messagingmock -out mock/event_bus.go . EventBus:EventBus
type EventBus interface {
	EventPublisher
	EventConsumer
}

// EventPublisher is an interface for publishing events to an event bus.
//
//go:generate moq -pkg messagingmock -out mock/event_publisher.go . EventPublisher:EventPublisher
type EventPublisher interface {
	// Publish emits one or more events. Implementations should provide at-least-once delivery semantics
	// unless otherwise documented.
	Publish(ctx context.Context, events ...Event) error
}

// EventConsumer is an interface for subscribing to events from an event bus.
//
//go:generate moq -pkg messagingmock -out mock/event_consumer.go . EventConsumer:EventConsumer
type EventConsumer interface {
	// Subscribe registers a handler for a given logical event name.
	// It returns an unsubscribe function that can be called to remove the subscription.
	Subscribe(ctx context.Context, h EventHandler[Event]) (UnsubscribeFunc, error)
}
