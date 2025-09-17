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

func TestInMemoryEventBus_Publish_NoSubscribers(t *testing.T) {
	t.Parallel()

	bus := messaging.NewInMemoryEventBus()
	evt := messaging.NewBaseEvent("event.no.subscribers")

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	err := bus.Publish(ctx, evt)
	require.Error(t, err)

	expectedErr := &messaging.NoSubscribersForMessageError{MessageType: evt.MessageType()}
	require.ErrorAs(t, err, &expectedErr)
	assert.Equal(t, "event.no.subscribers", expectedErr.MessageType)
}

func TestInMemoryEventBus_Subscribe_ThenHandleSync(t *testing.T) {
	t.Parallel()

	bus := messaging.NewInMemoryEventBus()
	const topic = "event.sync.topic"
	evt := messaging.NewBaseEvent(topic)

	seen := make(chan messaging.Event, 1)

	_, err := bus.Subscribe(context.Background(), topic,
		messaging.EventHandlerFn[messaging.Event](func(ctx context.Context, e messaging.Event) error {
			seen <- e
			return nil
		}),
	)
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	require.NoError(t, bus.Publish(ctx, evt))

	select {
	case got := <-seen:
		assert.Equal(t, topic, got.MessageType())
	case <-ctx.Done():
		t.Fatalf("handler was not invoked for %q", topic)
	}
}

func TestInMemoryEventBus_Unsubscribe_RemovesHandler(t *testing.T) {
	t.Parallel()

	bus := messaging.NewInMemoryEventBus()
	const topic = "event.unsubscribe"
	evt := messaging.NewBaseEvent(topic)

	calls := 0
	unsub, err := bus.Subscribe(context.Background(), topic,
		messaging.EventHandlerFn[messaging.Event](func(ctx context.Context, e messaging.Event) error {
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

	err = bus.Publish(ctx, evt)
	require.Error(t, err)

	expectedErr := &messaging.NoSubscribersForMessageError{MessageType: topic}
	require.ErrorAs(t, err, &expectedErr)

	assert.Equal(t, 0, calls, "handler should not have been called after unsubscribe")
}

func TestInMemoryEventBus_MiddlewareOrder(t *testing.T) {
	t.Parallel()

	bus := messaging.NewInMemoryEventBus()
	const topic = "event.mw.order"
	evt := messaging.NewBaseEvent(topic)

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
		messaging.EventHandlerFn[messaging.Event](func(ctx context.Context, e messaging.Event) error {
			order = append(order, "H")
			done <- struct{}{}
			return nil
		}),
	)
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	require.NoError(t, bus.Publish(ctx, evt))

	select {
	case <-done:
	case <-ctx.Done():
		t.Fatal("handler not called")
	}

	// wrap() applies middlewares in reverse registration order: A then B -> A> B> H <B <A
	assert.Equal(t, []string{"A>", "B>", "H", "<B", "<A"}, order)
}

func TestInMemoryEventBus_HandlerError_Propagates_WhenNoErrorHandler(t *testing.T) {
	t.Parallel()

	bus := messaging.NewInMemoryEventBus()
	const topic = "event.handler.error.nohandler"
	evt := messaging.NewBaseEvent(topic)

	want := errors.New("boom")

	_, err := bus.Subscribe(context.Background(), topic,
		messaging.EventHandlerFn[messaging.Event](func(ctx context.Context, e messaging.Event) error {
			return want
		}),
	)
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	err = bus.Publish(ctx, evt)
	require.Error(t, err)
	require.ErrorIs(t, err, want)
}

func TestInMemoryEventBus_HandlerError_RoutedToErrorHandler(t *testing.T) {
	t.Parallel()

	var gotType string
	var gotErr error

	bus := messaging.NewInMemoryEventBus(
		messageBusOpt(func(c *messaging.MessageBusConfig) {
			c.ErrorHandler = func(msgType string, err error) {
				gotType = msgType
				gotErr = err
			}
		}),
	)

	const topic = "event.handler.error.withhandler"
	evt := messaging.NewBaseEvent(topic)
	want := errors.New("kapow")

	_, err := bus.Subscribe(context.Background(), topic,
		messaging.EventHandlerFn[messaging.Event](func(ctx context.Context, e messaging.Event) error {
			return want
		}),
	)
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	// With ErrorHandler set, Publish should not return the handler error.
	err = bus.Publish(ctx, evt)
	require.NoError(t, err)

	assert.Equal(t, topic, gotType)
	require.ErrorIs(t, gotErr, want)
}

func TestInMemoryEventBus_Publish_MultipleEvents_AllDelivered(t *testing.T) {
	t.Parallel()

	bus := messaging.NewInMemoryEventBus()
	const topic = "event.multi"
	e1 := messaging.NewBaseEvent(topic)
	e2 := messaging.NewBaseEvent(topic)
	e3 := messaging.NewBaseEvent(topic)

	count := 0
	done := make(chan struct{}, 1)

	_, err := bus.Subscribe(context.Background(), topic,
		messaging.EventHandlerFn[messaging.Event](func(ctx context.Context, e messaging.Event) error {
			count++
			if count == 3 {
				done <- struct{}{}
			}
			return nil
		}),
	)
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	require.NoError(t, bus.Publish(ctx, e1, e2, e3))

	select {
	case <-done:
	case <-ctx.Done():
		t.Fatalf("expected all 3 events to be handled, got %d", count)
	}
	assert.Equal(t, 3, count)
}

func TestInMemoryEventBus_AsyncWorkers_ProcessEvents(t *testing.T) {
	t.Parallel()

	// Enable async pipeline with small queue to ensure we exercise the worker path.
	bus := messaging.NewInMemoryEventBus(
		messageBusOpt(func(c *messaging.MessageBusConfig) {
			c.AsyncWorkers = 2
			c.QueueSize = 2
		}),
	)
	const topic = "event.async"
	evt := messaging.NewBaseEvent(topic)

	seen := make(chan struct{}, 1)

	_, err := bus.Subscribe(context.Background(), topic,
		messaging.EventHandlerFn[messaging.Event](func(ctx context.Context, e messaging.Event) error {
			seen <- struct{}{}
			return nil
		}),
	)
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	require.NoError(t, bus.Publish(ctx, evt))

	select {
	case <-seen:
		// ok
	case <-ctx.Done():
		t.Fatal("async worker did not deliver event to handler in time")
	}
}

func TestInMemoryEventBus_Close_UnsubscribesAll(t *testing.T) {
	t.Parallel()

	bus := messaging.NewInMemoryEventBus()
	const subject = "event.close.unsub"
	event := messaging.NewBaseEvent(subject)

	_, err := bus.Subscribe(
		context.Background(),
		subject,
		messaging.EventHandlerFn[messaging.Event](func(_ context.Context, _ messaging.Event) error {
			return nil
		}),
	)
	require.NoError(t, err)

	// Close the bus, which should unsubscribe all handlers.
	err = bus.Close()
	require.NoError(t, err)

	// Attempting to dispatch should now fail with no subscribers.
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	err = bus.Publish(ctx, event)
	require.Error(t, err)
	require.ErrorIs(t, err, messaging.ErrPublishOnClosedBus)
}
