package event_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/xfrr/cqrsify/event"
)

func TestBus(t *testing.T) {
	var (
		mockSubject = "bus.test"
	)

	t.Run("NewBus", func(t *testing.T) {
		t.Run("should return a new bus", func(t *testing.T) {
			bus, err := event.NewBus()
			if err != nil {
				t.Errorf("expected err to be nil, got %v", err)
			}
			if bus == nil {
				t.Error("expected bus to not be nil")
			}
		})

		t.Run("should return a new bus with buffer size", func(t *testing.T) {
			bus, err := event.NewBus(
				event.WithBufferSize(100),
			)
			if err != nil {
				t.Errorf("expected err to be nil, got %v", err)
			}
			if bus == nil {
				t.Error("expected bus to not be nil")
			}
		})
	})

	t.Run("Publish", func(t *testing.T) {
		t.Run("should return an error when no subscribers are registered", func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			bus, _ := event.NewBus()
			evt := event.New("id", "reason", "payload").Any()
			err := bus.Publish(ctx, mockSubject, evt)
			if err == nil || !errors.Is(err, event.ErrNoSubscribers) {
				t.Fatalf("expected error to be %v, got %v", event.ErrNoSubscribers, err)
			}
		})

		t.Run("should return an error when context is canceled", func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			bus, _ := event.NewBus()
			_, err := bus.Subscribe(ctx, mockSubject)
			if err != nil {
				t.Fatalf("expected error to be nil, got %v", err)
			}

			// force cancel context
			cancel()
			err = bus.Publish(ctx, mockSubject, event.New("id", "reason", "payload").Any())
			if err == nil || !errors.Is(err, context.Canceled) {
				t.Fatalf("expected error to be %v, got %v", context.Canceled, err)
			}
		})

		t.Run("should execute the fallback function when publish times out", func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			done := make(chan struct{})
			fallback := func(ctx context.Context, subject string, evt event.Event[any, any]) {
				done <- struct{}{}
			}

			bus, _ := event.NewBus(
				event.WithPublishTimeout(1*time.Microsecond),
				event.WithPublishTimeoutFallback(fallback),
			)
			ch, err := bus.Subscribe(ctx, mockSubject)
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
						bus.Publish(ctx, mockSubject, event.New("id", "reason", "payload").Any())
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

		t.Run("should publish a event to all subscribers", func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			bus, _ := event.NewBus()

			ctxch, err := bus.Subscribe(ctx, mockSubject)
			if err != nil {
				t.Fatalf("expected error to be nil, got %v", err)
			}

			err = bus.Publish(ctx, mockSubject, event.New("id", "reason", "payload").Any())
			if err != nil {
				t.Fatalf("expected error to be nil, got %v", err)
			}

			select {
			case evt := <-ctxch:
				if evt.Event().ID() != "id" {
					t.Errorf("expected event id to be %v, got %v", "id", evt.Event().ID())
				}
				if evt.Event().Payload() != "payload" {
					t.Errorf("expected event payload to be %v, got %v", "payload", evt.Event().Payload())
				}
			case <-time.After(1 * time.Second):
				t.Fatal("expected event to be published")
			}
		})
	})

	t.Run("Subscribe", func(t *testing.T) {
		t.Run("should return a channel to receive events", func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			bus, _ := event.NewBus()
			ch, err := bus.Subscribe(ctx, mockSubject)
			if err != nil {
				t.Fatalf("expected error to be nil, got %v", err)
			}
			if ch == nil {
				t.Fatal("expected channel to not be nil")
			}
		})
	})
}
