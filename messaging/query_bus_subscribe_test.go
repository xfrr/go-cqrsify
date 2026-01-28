package messaging_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xfrr/go-cqrsify/messaging"
)

func TestSubscribeQuery_Success(t *testing.T) {
	t.Parallel()

	var calls int

	queryBus := messaging.NewInMemoryQueryBus(
		messaging.ConfigureInMemoryMessageBusSubjects("test.query", "test.reply"),
	)

	handler := messaging.MessageHandlerWithReplyFn[messaging.Query, messaging.QueryReply](func(_ context.Context, query messaging.Query) (messaging.QueryReply, error) {
		calls++
		return messaging.NewMessage(query.MessageType() + ".reply"), nil
	})

	unsub, err := messaging.RegisterQueryHandler(
		t.Context(),
		queryBus,
		handler,
	)
	require.NoError(t, err)

	testQuery := messaging.NewBaseQuery("test.query")
	res, err := queryBus.Request(t.Context(), testQuery)
	require.NoError(t, err)

	assert.NotNil(t, res)
	assert.Equal(t, 1, calls)
	assert.Equal(t, "test.query.reply", res.MessageType())

	err = unsub()
	require.NoError(t, err)
}
