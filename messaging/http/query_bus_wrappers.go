package messaginghttp

import (
	"context"
	"fmt"

	"github.com/xfrr/go-cqrsify/messaging"
)

var _ messaging.MessagePublisherReplier = (*queryDispatcherWrapper)(nil)

type queryDispatcherWrapper struct {
	dispatcher messaging.QueryDispatcher
}

func (w *queryDispatcherWrapper) PublishRequest(ctx context.Context, msg messaging.Message) (messaging.Message, error) {
	query, ok := msg.(messaging.Query)
	if !ok {
		return nil, fmt.Errorf("expected messaging.Query, got %T", msg)
	}

	return w.dispatcher.Request(ctx, query)
}
