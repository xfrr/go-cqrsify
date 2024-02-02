package command_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/xfrr/cqrsify/command"
)

func TestInMemoryBus(t *testing.T) {
	var (
		mockTopic = "bus.in_memory_test"
	)

	t.Run("NewInMemoryBus", func(t *testing.T) {
		t.Run("should return a new in-memory bus", func(t *testing.T) {
			bus, err := command.NewInMemoryBus()
			if err != nil {
				t.Errorf("expected err to be nil, got %v", err)
			}
			if bus == nil {
				t.Error("expected bus to not be nil")
			}
		})

		t.Run("should return a new in-memory bus with buffer size", func(t *testing.T) {
			bus, err := command.NewInMemoryBus(
				command.WithBufferSize(100),
			)
			if err != nil {
				t.Errorf("expected err to be nil, got %v", err)
			}
			if bus == nil {
				t.Error("expected bus to not be nil")
			}
		})
	})

	t.Run("Dispatch", func(t *testing.T) {
		t.Run("should return an error when no subscribers are registered", func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			bus, _ := command.NewInMemoryBus()
			err := bus.Dispatch(ctx, mockTopic, command.New("id", "msg").Any())
			if err == nil || !errors.Is(err, command.ErrNoSubscribers) {
				t.Fatalf("expected error to be %v, got %v", command.ErrNoSubscribers, err)
			}
		})

		t.Run("should return an error when context is canceled", func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			bus, _ := command.NewInMemoryBus()
			_, err := bus.Subscribe(ctx, mockTopic)
			if err != nil {
				t.Fatalf("expected error to be nil, got %v", err)
			}

			// force cancel context
			cancel()
			err = bus.Dispatch(ctx, mockTopic, command.New("id", "msg").Any())
			if err == nil || !errors.Is(err, context.Canceled) {
				t.Fatalf("expected error to be %v, got %v", context.Canceled, err)
			}
		})

		t.Run("should dispatch a command to all subscribers", func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			bus, _ := command.NewInMemoryBus()

			ctxch, err := bus.Subscribe(ctx, mockTopic)
			if err != nil {
				t.Fatalf("expected error to be nil, got %v", err)
			}

			err = bus.Dispatch(ctx, mockTopic, command.New("id", "msg").Any())
			if err != nil {
				t.Fatalf("expected error to be nil, got %v", err)
			}

			select {
			case cmd := <-ctxch:
				if cmd.Command().ID() != "id" {
					t.Errorf("expected command id to be %v, got %v", "id", cmd.Command().ID())
				}
				if cmd.Command().Message() != "msg" {
					t.Errorf("expected command message to be %v, got %v", "msg", cmd.Command().Message())
				}
			case <-time.After(1 * time.Second):
				t.Fatal("expected command to be dispatched")
			}
		})
	})

	t.Run("Subscribe", func(t *testing.T) {
		t.Run("should return a channel to receive commands", func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			bus, _ := command.NewInMemoryBus()
			ch, err := bus.Subscribe(ctx, mockTopic)
			if err != nil {
				t.Fatalf("expected error to be nil, got %v", err)
			}
			if ch == nil {
				t.Fatal("expected channel to not be nil")
			}
		})
	})
}
