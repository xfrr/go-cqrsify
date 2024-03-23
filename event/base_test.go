package event_test

import (
	"testing"

	"github.com/xfrr/cqrsify/event"
)

func TestNew(t *testing.T) {
	e := event.New("test", "test", &struct{}{},
		event.WithAggregate("aggregate-id", "aggregate-name", 1),
	)
	if e == nil {
		t.Fatal("event.New() should return a valid event")
	}

	if e.ID() == "" {
		t.Error("event.New() should return an event with a valid ID")
	}

	if e.Reason() != "test" {
		t.Error("event.New() should return an event with a valid reason")
	}

	if e.Payload() == nil {
		t.Error("event.New() should return an event with a valid payload")
	}

	if e.Time().IsZero() {
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
