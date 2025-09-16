package messaging

import (
	"context"
)

type EventHandler[E Event] = MessageHandler[E]
type EventHandlerFn[E Event] = MessageHandlerFn[E]

type EventBus interface {
	EventPublisher
	EventSubscriber
}

// EventPublisher is an interface for publishing events to an event bus.
type EventPublisher interface {
	// Publish emits one or more events. Implementations should provide at-least-once delivery semantics
	// unless otherwise documented.
	Publish(ctx context.Context, evt Event) error
}

// EventSubscriber is an interface for subscribing to events from an event bus.
type EventSubscriber interface {
	// Subscribe registers a handler for a given logical event name.
	// It returns an unsubscribe function that can be called to remove the subscription.
	Subscribe(ctx context.Context, subject string, h EventHandler[Event]) (func(), error)
}

// InMemoryEventBus is an in-memory implementation of EventBus.
type InMemoryEventBus struct {
	*InMemoryMessageBus
}

func NewInMemoryEventBus(optFns ...MessageBusConfigModifier) *InMemoryEventBus {
	return &InMemoryEventBus{
		InMemoryMessageBus: NewInMemoryMessageBus(optFns...),
	}
}

func (b *InMemoryEventBus) Publish(ctx context.Context, evt Event) error {
	return b.InMemoryMessageBus.Publish(ctx, evt)
}

func (b *InMemoryEventBus) Subscribe(ctx context.Context, eventName string, h EventHandler[Event]) (func(), error) {
	return b.InMemoryMessageBus.Subscribe(ctx, eventName, MessageHandlerFn[Message](func(ctx context.Context, msg Message) error {
		evt, ok := msg.(Event)
		if !ok {
			return nil
		}
		return h.Handle(ctx, evt)
	}))
}
