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

func TestInMemoryCommandBus_Dispatch_NoSubscribers(t *testing.T) {
	t.Parallel()

	bus := messaging.NewInMemoryCommandBus()
	cmd := messaging.NewBaseCommand("command.no.subscribers")

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	err := bus.Dispatch(ctx, cmd)
	require.Error(t, err)

	expectedErr := &messaging.NoSubscribersForMessageError{MessageType: cmd.MessageType()}
	require.ErrorAs(t, err, &expectedErr)
	assert.Equal(t, "command.no.subscribers", expectedErr.MessageType)
}

func TestInMemoryCommandBus_Subscribe_ThenHandleSync(t *testing.T) {
	t.Parallel()

	bus := messaging.NewInMemoryCommandBus()
	const topic = "command.sync.topic"
	cmd := messaging.NewBaseCommand(topic)

	seen := make(chan messaging.Command, 1)

	_, err := bus.Subscribe(context.Background(), topic,
		messaging.CommandHandlerFn[messaging.Command](func(_ context.Context, e messaging.Command) error {
			seen <- e
			return nil
		}),
	)
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	require.NoError(t, bus.Dispatch(ctx, cmd))

	select {
	case got := <-seen:
		assert.Equal(t, topic, got.MessageType())
	case <-ctx.Done():
		t.Fatalf("handler was not invoked for %q", topic)
	}
}

func TestInMemoryCommandBus_Unsubscribe_RemovesHandler(t *testing.T) {
	t.Parallel()

	bus := messaging.NewInMemoryCommandBus()
	const topic = "command.unsubscribe"
	cmd := messaging.NewBaseCommand(topic)

	calls := 0
	unsub, err := bus.Subscribe(context.Background(), topic,
		messaging.CommandHandlerFn[messaging.Command](func(_ context.Context, _ messaging.Command) error {
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

	err = bus.Dispatch(ctx, cmd)
	require.Error(t, err)

	expectedErr := &messaging.NoSubscribersForMessageError{MessageType: topic}
	require.ErrorAs(t, err, &expectedErr)

	assert.Equal(t, 0, calls, "handler should not have been called after unsubscribe")
}

func TestInMemoryCommandBus_MiddlewareOrder(t *testing.T) {
	t.Parallel()

	bus := messaging.NewInMemoryCommandBus()
	const topic = "command.mw.order"
	cmd := messaging.NewBaseCommand(topic)

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
		messaging.CommandHandlerFn[messaging.Command](func(_ context.Context, _ messaging.Command) error {
			order = append(order, "H")
			done <- struct{}{}
			return nil
		}),
	)
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	require.NoError(t, bus.Dispatch(ctx, cmd))

	select {
	case <-done:
	case <-ctx.Done():
		t.Fatal("handler not called")
	}

	// wrap() applies middlewares in reverse registration order: A then B -> A> B> H <B <A
	assert.Equal(t, []string{"A>", "B>", "H", "<B", "<A"}, order)
}

func TestInMemoryCommandBus_HandlerError_Propagates_WhenNoErrorHandler(t *testing.T) {
	t.Parallel()

	bus := messaging.NewInMemoryCommandBus()
	const topic = "command.handler.error.nohandler"
	cmd := messaging.NewBaseCommand(topic)

	want := errors.New("boom")

	_, err := bus.Subscribe(context.Background(), topic,
		messaging.CommandHandlerFn[messaging.Command](func(_ context.Context, _ messaging.Command) error {
			return want
		}),
	)
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	err = bus.Dispatch(ctx, cmd)
	require.Error(t, err)
	require.ErrorIs(t, err, want)
}

func TestInMemoryCommandBus_HandlerError_RoutedToErrorHandler(t *testing.T) {
	t.Parallel()

	var gotType string
	var gotErr error

	bus := messaging.NewInMemoryCommandBus(
		messageBusOpt(func(c *messaging.MessageBusConfig) {
			c.ErrorHandler = func(msgType string, err error) {
				gotType = msgType
				gotErr = err
			}
		}),
	)

	const topic = "command.handler.error.withhandler"
	cmd := messaging.NewBaseCommand(topic)
	want := errors.New("kapow")

	_, err := bus.Subscribe(context.Background(), topic,
		messaging.CommandHandlerFn[messaging.Command](func(_ context.Context, _ messaging.Command) error {
			return want
		}),
	)
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	// With ErrorHandler set, Dispatch should not return the handler error.
	err = bus.Dispatch(ctx, cmd)
	require.NoError(t, err)

	assert.Equal(t, topic, gotType)
	require.ErrorIs(t, gotErr, want)
}

func TestInMemoryCommandBus_Dispatch_MultipleCommands_AllDelivered(t *testing.T) {
	t.Parallel()

	bus := messaging.NewInMemoryCommandBus()
	const topic = "command.multi"
	e1 := messaging.NewBaseCommand(topic)
	e2 := messaging.NewBaseCommand(topic)
	e3 := messaging.NewBaseCommand(topic)

	count := 0
	done := make(chan struct{}, 1)

	_, err := bus.Subscribe(context.Background(), topic,
		messaging.CommandHandlerFn[messaging.Command](func(_ context.Context, _ messaging.Command) error {
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

	require.NoError(t, bus.Dispatch(ctx, e1, e2, e3))

	select {
	case <-done:
	case <-ctx.Done():
		t.Fatalf("expected all 3 commands to be handled, got %d", count)
	}
	assert.Equal(t, 3, count)
}

func TestInMemoryCommandBus_AsyncWorkers_ProcessCommands(t *testing.T) {
	t.Parallel()

	// Enable async pipeline with small queue to ensure we exercise the worker path.
	bus := messaging.NewInMemoryCommandBus(
		messageBusOpt(func(c *messaging.MessageBusConfig) {
			c.AsyncWorkers = 2
			c.QueueSize = 2
		}),
	)
	const topic = "command.async"
	cmd := messaging.NewBaseCommand(topic)

	seen := make(chan struct{}, 1)

	_, err := bus.Subscribe(context.Background(), topic,
		messaging.CommandHandlerFn[messaging.Command](func(_ context.Context, _ messaging.Command) error {
			seen <- struct{}{}
			return nil
		}),
	)
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	require.NoError(t, bus.Dispatch(ctx, cmd))

	select {
	case <-seen:
		// ok
	case <-ctx.Done():
		t.Fatal("async worker did not deliver command to handler in time")
	}
}
