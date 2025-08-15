package event

import (
	"context"

	"github.com/xfrr/go-cqrsify/message"
)

type Bus interface {
	Publish(ctx context.Context, evt Event) error
}

type InMemoryBus struct {
	bus *message.InMemoryBus
}

func (b *InMemoryBus) Publish(ctx context.Context, evt Event) error {
	return b.bus.Dispatch(ctx, evt)
}

func NewInMemoryBus() *InMemoryBus {
	return &InMemoryBus{
		bus: message.NewInMemoryBus(),
	}
}

func Handle[C Event](bus *InMemoryBus, handlerFn func(ctx context.Context, evt C) error) error {
	return message.Handle(bus.bus, handlerFn)
}
