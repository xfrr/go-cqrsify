package aggregate_test

import (
	"testing"

	"github.com/xfrr/go-cqrsify/aggregate"
	"github.com/xfrr/go-cqrsify/event"
)

func TestNextChange(t *testing.T) {
	const (
		aggID       = "agg-id"
		aggName     = "TestAggregate"
		eventID     = "evt-1"
		eventReason = "test"
	)

	type mockAggregate struct {
		*aggregate.Base[string]
	}

	agg := &mockAggregate{
		Base: aggregate.New(aggID, aggName),
	}

	handlerCalls := 0
	agg.When(eventReason, func(change event.Event[any, any]) {
		handlerCalls++
	})

	aggregate.NextChange(agg, eventID, eventReason, &struct{}{})
	if len(agg.AggregateChanges()) != 1 {
		t.Fatalf("expected aggregate events to have 1 event, got %d", len(agg.AggregateChanges()))
	}
	if handlerCalls != 1 {
		t.Fatalf("expected handler to be called 1 time, got %d", handlerCalls)
	}

	// apply second event
	aggregate.NextChange(agg, "evt-2", "test", &struct{}{})
	if len(agg.AggregateChanges()) != 2 {
		t.Fatalf("expected aggregate events to have 2 events, got %d", len(agg.AggregateChanges()))
	}
	if handlerCalls != 2 {
		t.Fatalf("expected handler to be called 2 times, got %d", handlerCalls)
	}
}
