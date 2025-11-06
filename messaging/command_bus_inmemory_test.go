package messaging_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xfrr/go-cqrsify/messaging"
)

func TestInMemoryCommandBus_Subscribe_ThenHandleSync(t *testing.T) {
	t.Parallel()

	const subject = "command.sync.subject"
	bus := messaging.NewInMemoryCommandBus(messaging.ConfigureInMemoryMessageBusSubjects(subject))
	cmd := messaging.NewBaseCommand(subject)

	seen := make(chan messaging.Command, 1)

	_, err := bus.Subscribe(
		context.Background(),
		messaging.MessageHandlerFn[messaging.Command](func(_ context.Context, e messaging.Command) error {
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
		assert.Equal(t, subject, got.MessageType())
	case <-ctx.Done():
		t.Fatalf("handler was not invoked for %q", subject)
	}
}
