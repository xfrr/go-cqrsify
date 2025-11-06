package messaging_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xfrr/go-cqrsify/messaging"
)

func TestInMemoryEventBus_Subscribe_ThenHandleSync(t *testing.T) {
	t.Parallel()

	const subject = "event.sync.subject"
	bus := messaging.NewInMemoryEventBus(messaging.ConfigureInMemoryMessageBusSubjects(subject))
	evt := messaging.NewBaseEvent(subject)

	seen := make(chan messaging.Event, 1)

	_, err := bus.Subscribe(context.Background(),
		messaging.MessageHandlerFn[messaging.Event](func(_ context.Context, e messaging.Event) error {
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
		assert.Equal(t, subject, got.MessageType())
	case <-ctx.Done():
		t.Fatalf("handler was not invoked for %q", subject)
	}
}
