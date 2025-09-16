package messaging

import (
	"context"
	"fmt"
)

// DispatchAndWaitReply sends a query and waits for a reply.
func DispatchAndWaitReply[Q Query, R any](
	ctx context.Context,
	dispatcher QueryDispatcher,
	qry Q,
) (R, error) {
	var zero R

	// send the query and wait for a reply
	replyMsg, err := dispatcher.DispatchAndWaitReply(ctx, qry)
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
