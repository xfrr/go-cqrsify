package messaging_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xfrr/go-cqrsify/messaging"
)

func TestInMemoryQueryBus_DispatchAndWaitReply_Success(t *testing.T) {
	t.Parallel()

	bus := messaging.NewInMemoryQueryBus()

	// Prepare a query with a subject/type the bus will route by.
	const subject = "test.query.success"
	q := messaging.NewBaseQuery(subject)

	// Subscribe a handler that replies immediately.
	_, err := bus.Subscribe(context.Background(), subject, messaging.QueryHandlerFn[messaging.Query](func(ctx context.Context, qry messaging.Query) error {
		// reply with nil (any Message is acceptable; nil keeps the test decoupled from concrete types)
		return qry.Reply(ctx, nil)
	}))
	require.NoError(t, err, "subscribe should succeed")

	// Dispatch and await the reply.
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	reply, err := bus.DispatchAndWaitReply(ctx, q)
	require.NoError(t, err, "expected a successful reply without error")
	assert.Nil(t, reply, "we sent a nil Message; expect nil back")
}

func TestInMemoryQueryBus_DispatchAndWaitReply_ErrorWhenNoSubscriber(t *testing.T) {
	t.Parallel()

	bus := messaging.NewInMemoryQueryBus()

	const subject = "test.query.timeout.no.subscriber"
	q := messaging.NewBaseQuery(subject)

	// No subscription for subject → DispatchAndWaitReply should block and then time out.
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
	defer cancel()

	time.Sleep(5 * time.Millisecond) // ensure the context timeout has taken effect

	reply, err := bus.DispatchAndWaitReply(ctx, q)
	require.Error(t, err, "expected context timeout when no handler replies")

	expectedErr := &messaging.NoSubscribersForMessageError{MessageType: subject}
	require.ErrorAs(t, err, &expectedErr)
	assert.Nil(t, reply)
}

func TestInMemoryQueryBus_HandlerCanDoWorkBeforeReply(t *testing.T) {
	t.Parallel()

	bus := messaging.NewInMemoryQueryBus()

	const subject = "test.query.handler.work"
	q := messaging.NewBaseQuery(subject)

	_, err := bus.Subscribe(context.Background(), subject, messaging.QueryHandlerFn[messaging.Query](func(ctx context.Context, qry messaging.Query) error {
		// Simulate some work but stay within the caller's timeout.
		select {
		case <-time.After(15 * time.Millisecond):
			return qry.Reply(ctx, nil)
		case <-ctx.Done():
			return ctx.Err()
		}
	}))
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	reply, err := bus.DispatchAndWaitReply(ctx, q)
	require.NoError(t, err)
	assert.Nil(t, reply)
}

func TestInMemoryQueryBus_Close_UnsubscribesAll(t *testing.T) {
	t.Parallel()

	bus := messaging.NewInMemoryQueryBus()

	const subject = "test.query.close"
	q := messaging.NewBaseQuery(subject)

	_, err := bus.Subscribe(context.Background(), subject, messaging.QueryHandlerFn[messaging.Query](func(ctx context.Context, qry messaging.Query) error {
		return qry.Reply(ctx, nil)
	}))
	require.NoError(t, err)

	// Close the bus, which should unsubscribe all handlers.
	err = bus.Close()
	require.NoError(t, err)

	// Attempting to dispatch should now fail with no subscribers.
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	reply, err := bus.DispatchAndWaitReply(ctx, q)
	require.Error(t, err)
	require.ErrorIs(t, err, messaging.ErrPublishOnClosedBus)
	assert.Nil(t, reply)
}
