package aggregate_test

import (
	"errors"
	"testing"

	"github.com/xfrr/go-cqrsify/aggregate"
	"github.com/xfrr/go-cqrsify/aggregate/event"
)

func TestVerifyHistoryIntegrity(t *testing.T) {
	tests := []struct {
		name   string
		agg    aggregate.Aggregate[string]
		events func() []aggregate.Event
		err    *aggregate.HistoryIntegrityError
	}{
		{
			name: "valid history",
			agg:  aggregate.New("agg-test-id", "agg-test"),
			events: func() []aggregate.Event {
				return makeEvents(t, "agg-test-id", "agg-test", 5)
			},
		},
		{
			name: "event has no aggregate",
			agg:  aggregate.New("agg-test-id", "agg-test"),
			events: func() []aggregate.Event {
				evt, err := event.New(
					"event-1",
					"event-name",
					"event-payload",
				)
				if err != nil {
					t.Fatal(err)
				}
				return []aggregate.Event{evt.Any()}
			},
			err: aggregate.NewHistoryIntegrityError("event has no aggregate"),
		},
		{
			name:   "event has different aggregate ID",
			agg:    aggregate.New("agg-test-id", "agg-test"),
			events: func() []aggregate.Event { return makeEvents(t, "agg-2", "agg-test", 5) },
			err:    aggregate.NewHistoryIntegrityError("event has different aggregate ID"),
		},
		{
			name:   "event has different aggregate name",
			agg:    aggregate.New("agg-test-id", "agg-test-other"),
			events: func() []aggregate.Event { return makeEvents(t, "agg-test-id", "agg-test", 5) },
			err:    aggregate.NewHistoryIntegrityError("event has different aggregate name"),
		},
		{
			name: "event has unexpected version",
			agg:  aggregate.New("agg-test-id", "agg-test-name"),
			events: func() []aggregate.Event {
				evt, err := event.New(
					"event-test-id",
					"event-test-name",
					"event-test-payload",
					event.WithAggregate("agg-test-id", "agg-test-name", 2),
				)
				if err != nil {
					t.Fatal(err)
				}

				return []aggregate.Event{evt.Any()}
			},
			err: aggregate.NewHistoryIntegrityError("event has unexpected version"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := aggregate.VerifyHistoryIntegrity(tt.agg, tt.events())
			if tt.err != nil {
				if !errors.As(err, &tt.err) {
					t.Fatalf("expected error %T, got %T", tt.err, err)
				}

				if tt.err.Error() != err.Error() {
					t.Fatalf("expected error message %s, got %s", tt.err.Error(), err.Error())
				}
			}

			if tt.err == nil && err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
		})
	}
}
