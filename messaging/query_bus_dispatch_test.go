package messaging_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xfrr/go-cqrsify/messaging"
)

func TestDispatchQuery_Success(t *testing.T) {
	t.Parallel()

	var calls int
	queryBus := messaging.NewInMemoryQueryBus(messaging.ConfigureInMemoryMessageBusSubjects("test.query"))
	handler := messaging.MessageHandlerWithReplyFn[messaging.Query, messaging.QueryReply](func(_ context.Context, _ messaging.Query) (messaging.QueryReply, error) {
		calls++

		return messaging.NewMessage("test.query.reply"), nil
	})

	unsub, err := queryBus.Subscribe(
		context.Background(),
		handler,
	)
	require.NoError(t, err)
	defer func() {
		unsubErr := unsub()
		require.NoError(t, unsubErr)
	}()

	testQuery := messaging.NewBaseQuery("test.query")
	res, err := messaging.DispatchQuery[messaging.Query, messaging.Message](context.Background(), queryBus, testQuery)
	require.NoError(t, err)

	assert.NotNil(t, res)
	assert.Equal(t, 1, calls)
	assert.Equal(t, "test.query.reply", res.MessageType())
}
