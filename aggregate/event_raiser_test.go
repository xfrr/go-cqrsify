package aggregate_test

import (
	"errors"
	"testing"

	"github.com/xfrr/go-cqrsify/aggregate"
	"github.com/xfrr/go-cqrsify/aggregate/event"
)

func TestRaiseEvent(t *testing.T) {
	const (
		aggID     = "agg-id"
		aggName   = "TestAggregate"
		eventID   = "evt-1"
		eventName = "test"
	)

	type mockAggregate struct {
		*aggregate.Base[string]
	}

	agg := &mockAggregate{
		Base: aggregate.New(aggID, aggName),
	}

	handlerCalls := 0
	agg.HandleEvent(eventName, func(_ event.Event[any, any]) {
		handlerCalls++
	})

	t.Run("should return error if aggregate is nil", func(t *testing.T) {
		err := aggregate.RaiseEvent[string, string, any](
			nil, eventID,
			eventName,
			&struct{}{},
		)
		if err == nil {
			t.Fatalf("expected error, got nil")
		}
		if !errors.Is(err, aggregate.NewRaiseEventError("aggregate is nil")) {
			t.Fatalf("expected error to be 'aggregate is nil', got %v", err)
		}
	})

	t.Run("should return error if event creation fails", func(t *testing.T) {
		err := aggregate.RaiseEvent[string, string, any](
			agg, eventID,
			eventName,
			nil,
		)
		if err == nil {
			t.Fatalf("expected error, got nil")
		}
		if !errors.Is(err, aggregate.NewRaiseEventErrorWithCause("failed to create event", event.NewValidationError("event payload is nil"))) {
			t.Fatalf("expected error to be 'failed to create event', got %v", err)
		}
	})

	t.Run("should raise events successfully", func(t *testing.T) {
		err := aggregate.RaiseEvent(agg, eventID, eventName, &struct{}{})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(agg.AggregateEvents()) != 1 {
			t.Fatalf("expected aggregate events to have 1 event, got %d", len(agg.AggregateEvents()))
		}
		if handlerCalls != 1 {
			t.Fatalf("expected handler to be called 1 time, got %d", handlerCalls)
		}

		// apply second event
		err = aggregate.RaiseEvent(agg, "evt-2", "test", &struct{}{})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(agg.AggregateEvents()) != 2 {
			t.Fatalf("expected aggregate events to have 2 events, got %d", len(agg.AggregateEvents()))
		}
		if handlerCalls != 2 {
			t.Fatalf("expected handler to be called 2 times, got %d", handlerCalls)
		}
	})
}
