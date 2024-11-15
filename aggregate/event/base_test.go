package event_test

import (
	"testing"

	"github.com/xfrr/go-cqrsify/aggregate/event"
)

func TestNew(t *testing.T) {
	e, err := event.New("test", "test", &struct{}{},
		event.WithAggregate("aggregate-id", "aggregate-name", 1),
	)
	if err != nil {
		t.Fatalf("event.New() should not return an error: %v", err)
	}
	if e == nil {
		t.Fatal("event.New() should return a valid event")
	}

	if e.ID() == "" {
		t.Error("event.New() should return an event with a valid ID")
	}

	if e.Name() != "test" {
		t.Error("event.New() should return an event with a valid name")
	}

	if e.Payload() == nil {
		t.Error("event.New() should return an event with a valid payload")
	}

	if e.OccurredAt().IsZero() {
		t.Error("event.New() should return an event with a valid timestamp")
	}

	if e.Aggregate().ID != "aggregate-id" {
		t.Error("event.New() should return an event with a valid aggregate reference")
	}

	if e.Aggregate().Name != "aggregate-name" {
		t.Error("event.New() should return an event with a valid aggregate reference")
	}

	if e.Aggregate().Version != 1 {
		t.Error("event.New() should return an event with a valid aggregate reference")
	}
}
