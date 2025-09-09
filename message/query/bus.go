package query

import (
	"context"

	"github.com/xfrr/go-cqrsify/message"
)

type Bus = message.Bus

func NewInMemoryBus() *message.InMemoryBus {
	return message.NewInMemoryBus()
}

func Handle[Q Query, R any](bus *message.InMemoryBus, topic string, handlerFn func(ctx context.Context, qry Q) (R, error)) error {
	return message.Handle(bus, topic, handlerFn)
}
