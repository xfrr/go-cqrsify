package message_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xfrr/go-cqrsify/message"
)

func TestPanicRecoveryHandlerMiddleware(t *testing.T) {
	var recovered any
	hook := func(r any) {
		recovered = r
	}

	middleware := message.HandlerPanicRecoveryMiddleware(hook)
	handler := middleware(message.HandlerFn[message.Message](func(ctx context.Context, msg message.Message) error {
		panic("test panic")
	}))

	err := handler.Handle(context.Background(), nil)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if recovered != "test panic" {
		t.Fatalf("expected recovered to be 'test panic', got %v", recovered)
	}
}

func TestComposeHandlerMiddlewares(t *testing.T) {
	var inc int
	middleware1 := message.HandlerMiddleware(func(next message.Handler[message.Message]) message.Handler[message.Message] {
		return message.HandlerFn[message.Message](func(ctx context.Context, msg message.Message) error {
			inc++
			return next.Handle(ctx, msg)
		})
	})

	middleware2 := message.HandlerMiddleware(func(next message.Handler[message.Message]) message.Handler[message.Message] {
		return message.HandlerFn[message.Message](func(ctx context.Context, msg message.Message) error {
			inc++
			return next.Handle(ctx, msg)
		})
	})

	handler := message.HandlerFn[message.Message](func(ctx context.Context, msg message.Message) error {
		return nil
	})

	chain := message.ChainHandlerMiddlewares(middleware1, middleware2)
	chainedHandler := chain(handler)

	err := chainedHandler.Handle(context.Background(), &messageMock{})
	require.NoError(t, err)

	assert.Equal(t, 2, inc, "expected middleware to be called twice")
}
