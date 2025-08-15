package message

import (
	"context"
)

var _ Bus = (*InMemoryBus)(nil)

type Bus interface {
	Dispatch(ctx context.Context, cmd Message) error
}
