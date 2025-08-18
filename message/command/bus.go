package command

import (
	"context"

	"github.com/xfrr/go-cqrsify/message"
)

type Bus = message.Bus

func NewInMemoryBus() *message.InMemoryBus {
	return message.NewInMemoryBus()
}

func Handle[C Command, R any](bus *message.InMemoryBus, handlerFn func(ctx context.Context, cmd C) (R, error)) error {
	return message.Handle(bus, handlerFn)
}
