package messaging

import (
	"context"
	"fmt"
)

// DispatchQuery sends a query and waits for a reply.
func DispatchQuery[Q Query, R any](
	ctx context.Context,
	dispatcher QueryDispatcher,
	query Q,
) (R, error) {
	var zero R

	// send the query and wait for a reply
	replyMsg, err := dispatcher.DispatchAndWaitReply(ctx, query)
	if err != nil {
		return zero, err
	}

	// cast the reply to the expected type
	reply, ok := replyMsg.(R)
	if !ok {
		return zero, InvalidMessageTypeError{
			Actual:   fmt.Sprintf("%T", replyMsg),
			Expected: fmt.Sprintf("%T", zero),
		}
	}

	return reply, nil
}
