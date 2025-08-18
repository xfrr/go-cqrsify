package event

import (
	"context"

	"github.com/xfrr/go-cqrsify/message"
)

var _ Bus = (*InMemoryBus)(nil)

type Bus interface {
	Publish(ctx context.Context, evt Event) error
}

type InMemoryBus struct {
	bus *message.InMemoryBus
}

func NewInMemoryBus() *InMemoryBus {
	return &InMemoryBus{
		bus: message.NewInMemoryBus(),
	}
}

func (b *InMemoryBus) Publish(ctx context.Context, evt Event) error {
	_, err := b.bus.Dispatch(ctx, evt)
	return err
}

func Handle[C Event](bus *InMemoryBus, handlerFn func(ctx context.Context, evt C) error) error {
	return message.Handle(bus.bus, message.HandlerFn[C, any](func(ctx context.Context, evt C) (any, error) {
		return nil, handlerFn(ctx, evt)
	}))
}
