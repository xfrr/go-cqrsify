package command

import (
	"context"

	"github.com/xfrr/go-cqrsify/message"
)

type Bus = message.Bus

func NewInMemoryBus() *message.InMemoryBus {
	return message.NewInMemoryBus()
}

func Handle[C Command](bus *message.InMemoryBus, handlerFn func(ctx context.Context, cmd C) error) error {
	return message.Handle(bus, handlerFn)
}
