package event_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/xfrr/go-cqrsify/aggregate/event"
)

func TestBus(t *testing.T) {
	var (
		mockSubject = "bus.test"
	)

	t.Run("NewBus", func(t *testing.T) {
		t.Run("should return a new bus", func(t *testing.T) {
			bus, err := event.NewInMemoryBus()
			if err != nil {
				t.Errorf("expected err to be nil, got %v", err)
			}
			if bus == nil {
				t.Error("expected bus to not be nil")
			}
		})

		t.Run("should return a new bus with buffer size", func(t *testing.T) {
			bus, err := event.NewInMemoryBus(
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
			bus, _ := event.NewInMemoryBus()
			evt, err := event.New("id", mockSubject, "payload")
			if err != nil {
				t.Fatal(err)
			}

			err = bus.Publish(ctx, evt.Any())
			if err == nil || !errors.Is(err, event.ErrBusHasNoSubscribers) {
				t.Fatalf("expected error to be %v, got %v", event.ErrBusHasNoSubscribers, err)
			}
		})

		t.Run("should publish a event to all subscribers", func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			bus, _ := event.NewInMemoryBus()

			ctxch, err := bus.Subscribe(ctx, mockSubject)
			if err != nil {
				t.Fatalf("expected error to be nil, got %v", err)
			}

			evt, err := event.New("id", mockSubject, "payload")
			if err != nil {
				t.Fatal(err)
			}

			err = bus.Publish(ctx, evt.Any())
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
			bus, _ := event.NewInMemoryBus()
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
