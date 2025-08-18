package message

import (
	"context"
)

var _ Bus = (*InMemoryBus)(nil)

type Bus interface {
	Dispatch(ctx context.Context, msg Message) (resp any, err error)
}
