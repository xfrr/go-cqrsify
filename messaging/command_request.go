package messaging

import (
	"context"
	"fmt"
)

// Request is a short-hand function to send a command and wait for a reply.
func Request[R Message](ctx context.Context, bus CommandBusReplier, cmd Command) (R, error) {
	res, err := bus.DispatchRequest(ctx, cmd)
	if err != nil {
		var zero R
		return zero, err
	}

	r, ok := res.(R)
	if !ok {
		var zero R
		return zero, fmt.Errorf("expected response of type %T but got %T", zero, res)
	}

	return r, nil
}
