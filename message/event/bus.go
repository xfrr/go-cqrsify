package event

import (
	"context"

	"github.com/xfrr/go-cqrsify/message"
)

var _ Bus = (*InMemoryBus)(nil)

type Bus interface {
	Publish(ctx context.Context, topic string, evt Event) error
}

type InMemoryBus struct {
	bus *message.InMemoryBus
}

func NewInMemoryBus() *InMemoryBus {
	return &InMemoryBus{
		bus: message.NewInMemoryBus(),
	}
}

func (b *InMemoryBus) Publish(ctx context.Context, topic string, evt Event) error {
	_, err := b.bus.Dispatch(ctx, topic, evt)
	return err
}

func Handle[C Event](bus *InMemoryBus, topic string, handlerFn func(ctx context.Context, evt C) error) error {
	return message.Handle(bus.bus, topic, message.HandlerFn[C, any](func(ctx context.Context, evt C) (any, error) {
		return nil, handlerFn(ctx, evt)
	}))
}
