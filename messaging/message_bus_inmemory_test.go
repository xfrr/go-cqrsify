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

// messageBusOpt is a tiny helper to build MessageBusConfigModifier inline.
func messageBusOpt(f func(*messaging.MessageBusConfig)) messaging.MessageBusConfigModifier {
	return func(c *messaging.MessageBusConfig) { f(c) }
}

func TestInMemoryMessageBus_Publish_NoSubscribers(t *testing.T) {
	t.Parallel()

	bus := messaging.NewInMemoryMessageBus()
	msg := messaging.NewBaseQuery("no.subscribers")

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	err := bus.Publish(ctx, msg)
	require.Error(t, err)

	expectedErr := &messaging.NoSubscribersForMessageError{MessageType: msg.MessageType()}
	require.ErrorAs(t, err, &expectedErr)
}

func TestInMemoryMessageBus_Subscribe_ThenHandleSync(t *testing.T) {
	t.Parallel()

	bus := messaging.NewInMemoryMessageBus()
	const topic = "sync.topic"
	msg := messaging.NewBaseQuery(topic)

	seen := make(chan messaging.Message, 1)

	_, err := bus.Subscribe(context.Background(), topic,
		messaging.MessageHandlerFn[messaging.Message](func(_ context.Context, m messaging.Message) error {
			seen <- m
			return nil
		}),
	)
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 150*time.Millisecond)
	defer cancel()

	require.NoError(t, bus.Publish(ctx, msg))

	select {
	case got := <-seen:
		assert.Equal(t, msg.MessageType(), got.MessageType())
	case <-ctx.Done():
		t.Fatal("handler was not invoked")
	}
}

func TestInMemoryMessageBus_Use_MiddlewareOrder(t *testing.T) {
	t.Parallel()

	bus := messaging.NewInMemoryMessageBus()
	const topic = "mw.topic"
	msg := messaging.NewBaseQuery(topic)

	var order []string

	// Middleware A
	mwA := func(next messaging.MessageHandler[messaging.Message]) messaging.MessageHandler[messaging.Message] {
		return messaging.MessageHandlerFn[messaging.Message](func(ctx context.Context, m messaging.Message) error {
			order = append(order, "A>") // enter A
			err := next.Handle(ctx, m)
			order = append(order, "<A") // exit A
			return err
		})
	}
	// Middleware B
	mwB := func(next messaging.MessageHandler[messaging.Message]) messaging.MessageHandler[messaging.Message] {
		return messaging.MessageHandlerFn[messaging.Message](func(ctx context.Context, m messaging.Message) error {
			order = append(order, "B>") // enter B
			err := next.Handle(ctx, m)
			order = append(order, "<B") // exit B
			return err
		})
	}

	bus.Use(mwA, mwB)

	done := make(chan struct{}, 1)
	_, err := bus.Subscribe(context.Background(), topic,
		messaging.MessageHandlerFn[messaging.Message](func(_ context.Context, _ messaging.Message) error {
			order = append(order, "H") // handler
			done <- struct{}{}
			return nil
		}),
	)
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 150*time.Millisecond)
	defer cancel()

	require.NoError(t, bus.Publish(ctx, msg))

	select {
	case <-done:
	case <-ctx.Done():
		t.Fatal("handler not called")
	}

	// wrap() applies middlewares in reverse registration order (last added wraps first).
	// So expected execution: A> B> H <B <A
	assert.Equal(t, []string{"A>", "B>", "H", "<B", "<A"}, order)
}

func TestInMemoryMessageBus_Close_PreventsPublish(t *testing.T) {
	t.Parallel()

	bus := messaging.NewInMemoryMessageBus()
	require.NoError(t, bus.Close())

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	err := bus.Publish(ctx, messaging.NewBaseQuery("anything"))
	require.Error(t, err)
	require.ErrorIs(t, err, messaging.ErrPublishOnClosedBus)
}

func TestInMemoryMessageBus_Subscribe_Unsubscribe_RemovesHandler(t *testing.T) {
	t.Parallel()

	bus := messaging.NewInMemoryMessageBus()
	const topic = "unsub.topic"
	msg := messaging.NewBaseQuery(topic)

	hit := 0

	unsub, err := bus.Subscribe(context.Background(), topic,
		messaging.MessageHandlerFn[messaging.Message](func(_ context.Context, _ messaging.Message) error {
			hit++
			return nil
		}),
	)
	require.NoError(t, err)
	require.NotNil(t, unsub)

	// Unsubscribe immediately.
	unsub()

	// Now there should be no subscribers; expect NoSubscribersForMessageError.
	ctx, cancel := context.WithTimeout(context.Background(), 1000*time.Millisecond)
	defer cancel()

	err = bus.Publish(ctx, msg)
	require.Error(t, err)

	expectedError := &messaging.NoSubscribersForMessageError{MessageType: topic}
	require.ErrorAs(t, err, &expectedError)

	assert.Equal(t, 0, hit, "handler should not have been called after unsubscribe")
}

func TestInMemoryMessageBus_HandlerError_PropagatesWhenNoErrorHandler(t *testing.T) {
	t.Parallel()

	bus := messaging.NewInMemoryMessageBus()
	const topic = "handler.error.nohandler"
	msg := messaging.NewBaseQuery(topic)

	want := errors.New("boom")

	_, err := bus.Subscribe(context.Background(), topic,
		messaging.MessageHandlerFn[messaging.Message](func(_ context.Context, _ messaging.Message) error {
			return want
		}),
	)
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 150*time.Millisecond)
	defer cancel()

	err = bus.Publish(ctx, msg)
	require.Error(t, err)
	require.ErrorIs(t, err, want)
}

func TestInMemoryMessageBus_HandlerError_ReportedViaErrorHandler(t *testing.T) {
	t.Parallel()

	var capturedType string
	var capturedErr error

	bus := messaging.NewInMemoryMessageBus(
		messageBusOpt(func(c *messaging.MessageBusConfig) {
			c.ErrorHandler = func(msgType string, err error) {
				capturedType = msgType
				capturedErr = err
			}
		}),
	)

	const topic = "handler.error.withhandler"
	msg := messaging.NewBaseQuery(topic)
	want := errors.New("boom")

	_, err := bus.Subscribe(context.Background(), topic,
		messaging.MessageHandlerFn[messaging.Message](func(_ context.Context, _ messaging.Message) error {
			return want
		}),
	)
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 150*time.Millisecond)
	defer cancel()

	// When ErrorHandler is set, Publish should not return the handler error.
	err = bus.Publish(ctx, msg)
	require.NoError(t, err)

	assert.Equal(t, topic, capturedType)
	require.ErrorIs(t, capturedErr, want)
}

func TestInMemoryMessageBus_Publish_PropagatesHandlerCtxError(t *testing.T) {
	t.Parallel()

	bus := messaging.NewInMemoryMessageBus()
	const topic = "ctx.cancel.propagates"
	msg := messaging.NewBaseQuery(topic)

	_, err := bus.Subscribe(context.Background(), topic,
		messaging.MessageHandlerFn[messaging.Message](func(ctx context.Context, _ messaging.Message) error {
			// Simulate handler checking ctx and returning its error.
			<-ctx.Done()
			return ctx.Err()
		}),
	)
	require.NoError(t, err)

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel before publish; handler should see context canceled and return error

	err = bus.Publish(ctx, msg)
	require.Error(t, err)
	assert.True(t, errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded))
}
