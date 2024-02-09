package command_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/xfrr/cqrsify/command"
)

func TestBus(t *testing.T) {
	var (
		mockTopic = "bus.test"
	)

	t.Run("NewBus", func(t *testing.T) {
		t.Run("should return a new bus", func(t *testing.T) {
			bus, err := command.NewBus()
			if err != nil {
				t.Errorf("expected err to be nil, got %v", err)
			}
			if bus == nil {
				t.Error("expected bus to not be nil")
			}
		})

		t.Run("should return a new bus with buffer size", func(t *testing.T) {
			bus, err := command.NewBus(
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
			bus, _ := command.NewBus()
			err := bus.Dispatch(ctx, mockTopic, command.New("id", "msg").Any())
			if err == nil || !errors.Is(err, command.ErrNoSubscribers) {
				t.Fatalf("expected error to be %v, got %v", command.ErrNoSubscribers, err)
			}
		})

		t.Run("should return an error when context is canceled", func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			bus, _ := command.NewBus()
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

		t.Run("should execute the fallback function when dispatch times out", func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			done := make(chan struct{})
			fallback := func(ctx context.Context, topic string, cmd command.Command[any]) {
				done <- struct{}{}
			}

			bus, _ := command.NewBus(
				command.WithDispatchTimeout(1*time.Microsecond),
				command.WithDispatchTimeoutFallback(fallback),
			)
			ch, err := bus.Subscribe(ctx, mockTopic)
			if err != nil {
				t.Fatalf("expected error to be nil, got %v", err)
			}

			// simulate slow subscriber
			go func() {
				for {
					select {
					case <-ctx.Done():
						return
					default:
						time.Sleep(1 * time.Second)

						<-ch
					}
				}
			}()

			go func() {
				for {
					select {
					case <-ctx.Done():
						return
					default:
						bus.Dispatch(ctx, mockTopic, command.New("id", "msg").Any())
					}
				}
			}()

			for {
				select {
				case <-done:
					return
				case <-time.After(2 * time.Second):
					t.Fatal("expected fallback function to be executed")
				}
			}
		})

		t.Run("should dispatch a command to all subscribers", func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			bus, _ := command.NewBus()

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
			bus, _ := command.NewBus()
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
