package messaging_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xfrr/go-cqrsify/messaging"
)

// TestQueryHandlerTypedRouter_RegisterAndHandle tests successful registration and dispatch.
func TestQueryHandlerTypedRouter_RegisterAndHandle(t *testing.T) {
	t.Parallel()

	const queryType = "test.query"
	router := messaging.NewQueryHandlerTypedRouter()

	handlerCalled := false
	handler := messaging.MessageHandlerWithReplyFn[messaging.Query, messaging.QueryReply](
		func(_ context.Context, _ messaging.Query) (messaging.QueryReply, error) {
			handlerCalled = true
			return messaging.NewMessage(queryType + ".reply"), nil
		},
	)

	err := messaging.RegisterQueryHandlerTypedRouter[messaging.Query, messaging.QueryReply](
		&router.MessageHandlerWithReplyTypedRouter,
		queryType,
		handler,
	)
	require.NoError(t, err)

	query := messaging.NewBaseQuery(queryType)
	reply, err := router.Handle(context.Background(), query)
	require.NoError(t, err)
	assert.True(t, handlerCalled)
	assert.NotNil(t, reply)
	assert.Equal(t, queryType+".reply", reply.MessageType())
}

// TestQueryHandlerTypedRouter_DirectRegister tests direct Register method access.
func TestQueryHandlerTypedRouter_DirectRegister(t *testing.T) {
	t.Parallel()

	const queryType = "test.query.direct"
	router := messaging.NewQueryHandlerTypedRouter()

	handlerCalled := false
	handler := messaging.MessageHandlerWithReplyFn[messaging.Query, messaging.QueryReply](
		func(_ context.Context, _ messaging.Query) (messaging.QueryReply, error) {
			handlerCalled = true
			return messaging.NewMessage(queryType + ".reply"), nil
		},
	)

	err := router.Register(queryType, handler)
	require.NoError(t, err)

	query := messaging.NewBaseQuery(queryType)
	reply, err := router.Handle(context.Background(), query)
	require.NoError(t, err)
	assert.True(t, handlerCalled)
	assert.Equal(t, queryType+".reply", reply.MessageType())
}

// TestQueryHandlerTypedRouter_DuplicateRegistration tests that duplicate registrations are rejected.
func TestQueryHandlerTypedRouter_DuplicateRegistration(t *testing.T) {
	t.Parallel()

	const queryType = "test.query.dup"
	router := messaging.NewQueryHandlerTypedRouter()

	handler := messaging.MessageHandlerWithReplyFn[messaging.Query, messaging.QueryReply](
		func(_ context.Context, _ messaging.Query) (messaging.QueryReply, error) {
			return messaging.NewMessage(queryType + ".reply"), nil
		},
	)

	err := router.Register(queryType, handler)
	require.NoError(t, err)

	// Attempt to register the same query type again
	err = router.Register(queryType, handler)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "handler already exists")
}

// TestQueryHandlerTypedRouter_MissingHandler tests that missing handlers return ErrHandlerNotFound.
func TestQueryHandlerTypedRouter_MissingHandler(t *testing.T) {
	t.Parallel()

	router := messaging.NewQueryHandlerTypedRouter()

	query := messaging.NewBaseQuery("test.query.missing")
	_, err := router.Handle(context.Background(), query)
	require.Error(t, err)
	assert.ErrorIs(t, err, messaging.ErrHandlerNotFound)
}

// TestQueryHandlerTypedRouter_TypedCasting tests handler with a concrete query type.
func TestQueryHandlerTypedRouter_TypedCasting(t *testing.T) {
	t.Parallel()

	const queryType = "test.query.typed"
	router := messaging.NewQueryHandlerTypedRouter()

	type customQuery struct {
		messaging.BaseQuery
		Field string
	}

	type customReply struct {
		messaging.BaseQueryReply
		Result string
	}

	handlerCalled := false
	typedHandler := messaging.MessageHandlerWithReplyFn[customQuery, customReply](
		func(_ context.Context, qry customQuery) (customReply, error) {
			handlerCalled = true
			return customReply{
				BaseQueryReply: messaging.NewMessage(queryType + ".reply"),
				Result:         "received: " + qry.Field,
			}, nil
		},
	)

	err := messaging.RegisterQueryHandlerTypedRouter[customQuery, customReply](
		&router.MessageHandlerWithReplyTypedRouter,
		queryType,
		typedHandler,
	)
	require.NoError(t, err)

	// Dispatch the typed query
	query := customQuery{
		BaseQuery: messaging.NewBaseQuery(queryType),
		Field:     "test-value",
	}
	reply, err := router.Handle(context.Background(), query)
	require.NoError(t, err)
	assert.True(t, handlerCalled)
	assert.NotNil(t, reply)
	assert.Equal(t, queryType+".reply", reply.MessageType())
}

// TestQueryHandlerTypedRouter_HandlerError tests error propagation from handler.
func TestQueryHandlerTypedRouter_HandlerError(t *testing.T) {
	t.Parallel()

	const queryType = "test.query.error"
	router := messaging.NewQueryHandlerTypedRouter()

	expectedErr := errors.New("handler error")
	handler := messaging.MessageHandlerWithReplyFn[messaging.Query, messaging.QueryReply](
		func(_ context.Context, _ messaging.Query) (messaging.QueryReply, error) {
			return nil, expectedErr
		},
	)

	err := router.Register(queryType, handler)
	require.NoError(t, err)

	query := messaging.NewBaseQuery(queryType)
	_, err = router.Handle(context.Background(), query)
	require.Error(t, err)
	assert.Equal(t, expectedErr, err)
}

// TestQueryHandlerTypedRouter_ContextCancellation tests context cancellation propagation.
func TestQueryHandlerTypedRouter_ContextCancellation(t *testing.T) {
	t.Parallel()

	const queryType = "test.query.cancel"
	router := messaging.NewQueryHandlerTypedRouter()

	handler := messaging.MessageHandlerWithReplyFn[messaging.Query, messaging.QueryReply](
		func(ctx context.Context, _ messaging.Query) (messaging.QueryReply, error) {
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			default:
				return messaging.NewMessage(queryType + ".reply"), nil
			}
		},
	)

	err := router.Register(queryType, handler)
	require.NoError(t, err)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	query := messaging.NewBaseQuery(queryType)
	_, err = router.Handle(ctx, query)
	require.Error(t, err)
	assert.ErrorIs(t, err, context.Canceled)
}
