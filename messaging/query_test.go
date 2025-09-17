package messaging_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xfrr/go-cqrsify/messaging"
)

func TestInMemoryQueryBus_Dispatch_NoSubscribers(t *testing.T) {
	t.Parallel()

	bus := messaging.NewInMemoryQueryBus()
	qry := messaging.NewBaseQuery("query.no.subscribers")

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	res, err := bus.DispatchAndWaitReply(ctx, qry)
	require.Error(t, err)
	assert.Nil(t, res)

	expectedErr := &messaging.NoSubscribersForMessageError{MessageType: qry.MessageType()}
	require.ErrorAs(t, err, &expectedErr)
	assert.Equal(t, "query.no.subscribers", expectedErr.MessageType)
}

func TestInMemoryQueryBus_Subscribe_ThenHandleSync(t *testing.T) {
	t.Parallel()

	bus := messaging.NewInMemoryQueryBus()
	const topic = "query.sync.topic"
	qry := messaging.NewBaseQuery(topic)

	seen := make(chan messaging.Query, 1)

	_, err := bus.Subscribe(context.Background(), topic,
		messaging.QueryHandlerFn[messaging.Query](func(_ context.Context, query messaging.Query) error {
			seen <- query
			// Reply to the query to unblock DispatchAndWaitReply
			query.Reply(context.Background(), query) // Echo reply
			return nil
		}),
	)
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	res, err := bus.DispatchAndWaitReply(ctx, qry)
	require.NoError(t, err)
	assert.Equal(t, qry, res)

	select {
	case got := <-seen:
		assert.Equal(t, topic, got.MessageType())
	case <-ctx.Done():
		t.Fatalf("handler was not invoked for %q", topic)
	}
}

func TestInMemoryQueryBus_Unsubscribe_RemovesHandler(t *testing.T) {
	t.Parallel()

	bus := messaging.NewInMemoryQueryBus()
	const topic = "query.unsubscribe"
	qry := messaging.NewBaseQuery(topic)

	calls := 0
	unsub, err := bus.Subscribe(context.Background(), topic,
		messaging.QueryHandlerFn[messaging.Query](func(_ context.Context, _ messaging.Query) error {
			calls++
			return nil
		}),
	)
	require.NoError(t, err)
	require.NotNil(t, unsub)

	// Unsubscribe immediately.
	unsub()
	require.NoError(t, err)

	// Now there should be no subscribers → expect NoSubscribersForMessageError.
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	res, err := bus.DispatchAndWaitReply(ctx, qry)
	require.Error(t, err)
	assert.Nil(t, res)

	expectedErr := &messaging.NoSubscribersForMessageError{MessageType: topic}
	require.ErrorAs(t, err, &expectedErr)

	assert.Equal(t, 0, calls, "handler should not have been called after unsubscribe")
}

func TestInMemoryQueryBus_MiddlewareOrder(t *testing.T) {
	t.Parallel()

	bus := messaging.NewInMemoryQueryBus()
	const topic = "query.mw.order"
	query := messaging.NewBaseQuery(topic)
	queryReply := messaging.NewBaseQuery(query.MessageType() + ".reply")

	var order []string

	// Middleware A (outermost when executed)
	mwA := func(next messaging.MessageHandler[messaging.Message]) messaging.MessageHandler[messaging.Message] {
		return messaging.MessageHandlerFn[messaging.Message](func(ctx context.Context, m messaging.Message) error {
			order = append(order, "A>")
			err := next.Handle(ctx, m)
			order = append(order, "<A")
			return err
		})
	}
	// Middleware B (innermost when executed if added after A)
	mwB := func(next messaging.MessageHandler[messaging.Message]) messaging.MessageHandler[messaging.Message] {
		return messaging.MessageHandlerFn[messaging.Message](func(ctx context.Context, m messaging.Message) error {
			order = append(order, "B>")
			err := next.Handle(ctx, m)
			order = append(order, "<B")
			return err
		})
	}

	bus.Use(mwA, mwB)

	done := make(chan struct{}, 1)
	_, err := bus.Subscribe(context.Background(), topic,
		messaging.QueryHandlerFn[messaging.Query](func(_ context.Context, query messaging.Query) error {
			order = append(order, "H")
			done <- struct{}{}
			// Reply to the query to unblock DispatchAndWaitReply
			query.Reply(context.Background(), queryReply) // Echo reply
			return nil
		}),
	)
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	res, err := bus.DispatchAndWaitReply(ctx, query)
	require.NoError(t, err)
	assert.Equal(t, queryReply, res)

	select {
	case <-done:
	case <-ctx.Done():
		t.Fatal("handler not called")
	}

	// wrap() applies middlewares in reverse registration order: A then B -> A> B> H <B <A
	assert.Equal(t, []string{"A>", "B>", "H", "<B", "<A"}, order)
}

func TestInMemoryQueryBus_HandlerError_Propagates_WhenNoErrorHandler(t *testing.T) {
	t.Parallel()

	bus := messaging.NewInMemoryQueryBus()
	const topic = "query.handler.error.nohandler"
	qry := messaging.NewBaseQuery(topic)

	want := errors.New("boom")

	_, err := bus.Subscribe(context.Background(), topic,
		messaging.QueryHandlerFn[messaging.Query](func(_ context.Context, _ messaging.Query) error {
			return want
		}),
	)
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	res, err := bus.DispatchAndWaitReply(ctx, qry)
	require.Error(t, err)
	require.ErrorIs(t, err, want)
	assert.Nil(t, res)
}

func TestInMemoryQueryBus_HandlerError_RoutedToErrorHandler_Timeout(t *testing.T) {
	t.Parallel()

	var gotType string
	var gotErr error

	bus := messaging.NewInMemoryQueryBus(
		messageBusOpt(func(c *messaging.MessageBusConfig) {
			c.ErrorHandler = func(msgType string, err error) {
				gotType = msgType
				gotErr = err
			}
		}),
	)

	const topic = "query.handler.error.withhandler"
	qry := messaging.NewBaseQuery(topic)
	want := errors.New("kapow")

	_, err := bus.Subscribe(context.Background(), topic,
		messaging.QueryHandlerFn[messaging.Query](func(_ context.Context, query messaging.Query) error {
			return want
		}),
	)
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 1000*time.Millisecond)
	defer cancel()

	// With ErrorHandler set, Dispatch should not return the handler error.
	res, err := bus.DispatchAndWaitReply(ctx, qry)
	require.Error(t, err)
	require.ErrorIs(t, err, context.DeadlineExceeded)
	assert.Nil(t, res)
	assert.Equal(t, topic, gotType)
	assert.Equal(t, want, gotErr)
}
