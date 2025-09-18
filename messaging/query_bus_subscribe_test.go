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

	queryBus := messaging.NewInMemoryQueryBus()
	handler := messaging.QueryHandlerFn[messaging.Query](func(ctx context.Context, query messaging.Query) error {
		calls++
		err := query.Reply(ctx, messaging.NewBaseQuery("test.reply"))
		require.NoError(t, err)
		return nil
	})

	unsub, err := messaging.SubscribeQuery(
		context.Background(),
		queryBus,
		"test.query",
		handler,
	)
	require.NoError(t, err)
	defer unsub()

	testQuery := messaging.NewBaseQuery("test.query")
	res, err := queryBus.DispatchAndWaitReply(context.Background(), testQuery)
	require.NoError(t, err)

	assert.NotNil(t, res)
	assert.Equal(t, 1, calls)
	assert.Equal(t, "test.reply", res.MessageType())
}
