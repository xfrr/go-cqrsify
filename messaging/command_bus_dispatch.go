package messaging

import (
	"context"
)

// DispatchCommand is a shorthand for sending commands.
func DispatchCommand(ctx context.Context, dispatcher CommandDispatcher, cmd Command) error {
	err := dispatcher.Dispatch(ctx, cmd)
	if err != nil {
		return err
	}
	return nil
}
