package messaging

import (
	"context"
)

// DispatchCommand is a shorthand for sending commands.
func DispatchCommand[C Command](ctx context.Context, dispatcher CommandDispatcher, cmd C) error {
	err := dispatcher.Dispatch(ctx, cmd)
	if err != nil {
		return err
	}
	return nil
}
