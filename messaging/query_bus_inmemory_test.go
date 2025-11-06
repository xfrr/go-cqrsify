package messaging_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xfrr/go-cqrsify/messaging"
)

func TestInMemoryQueryBus_Subscribe_ThenHandleSync(t *testing.T) {
	t.Parallel()

	const subject = "query.sync.topic"
	bus := messaging.NewInMemoryQueryBus(messaging.ConfigureInMemoryMessageBusSubjects(subject))
	queryReplyMsg := messaging.NewMessage(subject + ".reply")

	seen := make(chan messaging.Query, 1)

	_, err := bus.Subscribe(context.Background(),
		messaging.MessageHandlerWithReplyFn[messaging.Query, messaging.QueryReply](func(_ context.Context, query messaging.Query) (messaging.QueryReply, error) {
			seen <- query
			return queryReplyMsg, nil
		}),
	)
	require.NoError(t, err)

	res, err := bus.Request(t.Context(), messaging.NewBaseQuery(subject))
	require.NoError(t, err)
	assert.Equal(t, queryReplyMsg, res)

	select {
	case got := <-seen:
		assert.Equal(t, subject, got.MessageType())
	case <-t.Context().Done():
		t.Fatalf("handler was not invoked for %q", subject)
	}
}

func TestInMemoryQueryBus_Dispatch_NoHandlers(t *testing.T) {
	t.Parallel()

	bus := messaging.NewInMemoryQueryBus()
	query := messaging.NewBaseQuery("query.no.handlers")

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	res, err := bus.Request(ctx, query)
	require.Error(t, err)
	assert.Nil(t, res)

	expectedErr := messaging.NoHandlersForMessageError{MessageType: query.MessageType()}
	require.ErrorIs(t, err, expectedErr)
	assert.Equal(t, "query.no.handlers", expectedErr.MessageType)
}
