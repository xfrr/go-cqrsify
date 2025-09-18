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

	type testQueryReply struct {
		messaging.BaseQueryReply
	}

	var calls int
	queryBus := messaging.NewInMemoryQueryBus()
	handler := messaging.QueryHandlerFn[messaging.Query](func(ctx context.Context, query messaging.Query) error {
		calls++

		err := query.Reply(ctx, testQueryReply{
			BaseQueryReply: messaging.NewBaseQueryReply(query),
		})
		require.NoError(t, err)
		return nil
	})

	unsub, err := queryBus.Subscribe(
		context.Background(),
		"test.query",
		handler,
	)
	require.NoError(t, err)
	defer unsub()

	testQuery := messaging.NewBaseQuery("test.query")
	res, err := messaging.DispatchQuery[messaging.Query, messaging.Message](context.Background(), queryBus, testQuery)
	require.NoError(t, err)

	assert.NotNil(t, res)
	assert.Equal(t, 1, calls)
	assert.Equal(t, "test.query.reply", res.MessageType())
}
