package messaging

import (
	"context"
	"fmt"
)

var _ EventBus = (*InMemoryEventBus)(nil)

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
	Publish(ctx context.Context, events ...Event) error
}

// EventSubscriber is an interface for subscribing to events from an event bus.
type EventSubscriber interface {
	// Subscribe registers a handler for a given logical event name.
	// It returns an unsubscribe function that can be called to remove the subscription.
	Subscribe(ctx context.Context, subject string, h EventHandler[Event]) (UnsubscribeFunc, error)
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

func (b *InMemoryEventBus) Publish(ctx context.Context, events ...Event) error {
	msgs := make([]Message, len(events))
	for i, e := range events {
		msgs[i] = e
	}
	return b.InMemoryMessageBus.Publish(ctx, msgs...)
}

func (b *InMemoryEventBus) Subscribe(ctx context.Context, eventName string, h EventHandler[Event]) (UnsubscribeFunc, error) {
	return b.InMemoryMessageBus.Subscribe(ctx, eventName, MessageHandlerFn[Message](func(ctx context.Context, msg Message) error {
		evt, ok := msg.(Event)
		if !ok {
			return InvalidMessageTypeError{Expected: fmt.Sprintf("%T", evt), Actual: fmt.Sprintf("%T", msg)}
		}

		return h.Handle(ctx, evt)
	}))
}
