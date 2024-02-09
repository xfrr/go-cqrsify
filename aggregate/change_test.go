package aggregate_test

import (
	"testing"

	"github.com/xfrr/cqrsify/aggregate"
)

func TestApplyEvent(t *testing.T) {
	const (
		aggID       = "agg-id"
		aggName     = "TestAggregate"
		eventID     = "evt-1"
		eventReason = "test"
	)

	type mockAggregate struct {
		*aggregate.Base
	}

	agg := &mockAggregate{
		Base: aggregate.New(aggID, aggName),
	}

	evt := aggregate.ApplyChange(agg, eventID, eventReason, &struct{}{})

	if evt.ID() != eventID {
		t.Errorf("expected event id to be evt-1, got %s", evt.ID())
	}

	if evt.Reason() != eventReason {
		t.Errorf("expected event reason to be test, got %s", evt.Reason())
	}

	if evt.Aggregate().ID != aggID {
		t.Errorf("expected aggregate id to be 1, got %s", evt.Aggregate().ID)
	}

	if evt.Aggregate().Name != aggName {
		t.Errorf("expected aggregate name to be TestAggregate, got %s", evt.Aggregate().Name)
	}

	if evt.Aggregate().Version != 1 {
		t.Errorf("expected aggregate version to be 1, got %d", evt.Aggregate().Version)
	}

	if len(agg.AggregateChanges()) != 1 {
		t.Fatalf("expected aggregate events to have 1 event, got %d", len(agg.AggregateChanges()))
	}

	// apply second event
	evt2 := aggregate.ApplyChange(agg, "evt-2", "test", &struct{}{})
	if evt2.Aggregate().Version != 2 {
		t.Errorf("expected aggregate version to be 2, got %d", evt2.Aggregate().Version)
	}

	if len(agg.AggregateChanges()) != 2 {
		t.Fatalf("expected aggregate events to have 2 events, got %d", len(agg.AggregateChanges()))
	}
}
