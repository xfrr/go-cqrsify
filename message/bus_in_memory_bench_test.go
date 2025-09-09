package message_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/xfrr/go-cqrsify/message"
)

func BenchmarkInMemoryBus_Dispatch(b *testing.B) {
	bus := message.NewInMemoryBus()
	handler := &handlerWrapper{
		fn: func(ctx context.Context, msg message.Message) (any, error) {
			return nil, nil
		},
	}

	err := bus.RegisterHandler("com.org.test_message", handler)
	require.NoError(b, err)

	msg := TestMessage{message.NewBase()}
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = bus.Dispatch(ctx, "com.org.test_message", msg)
	}
}

func BenchmarkInMemoryBus_DispatchWithMiddleware(b *testing.B) {
	bus := message.NewInMemoryBus()
	handler := &handlerWrapper{
		fn: func(ctx context.Context, msg message.Message) (any, error) {
			return nil, nil
		},
	}

	// Add a simple middleware
	middleware := func(h message.Handler[message.Message, any]) message.Handler[message.Message, any] {
		return &handlerWrapper{
			fn: func(ctx context.Context, msg message.Message) (any, error) {
				return h.Handle(ctx, msg)
			},
		}
	}

	err := bus.RegisterHandler("com.org.test_message", handler)
	require.NoError(b, err)

	bus.Use(middleware)

	msg := TestMessage{message.NewBase()}
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = bus.Dispatch(ctx, "com.org.test_message", msg)
	}
}
