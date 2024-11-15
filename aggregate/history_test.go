package aggregate_test

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"testing"

	"github.com/xfrr/go-cqrsify/aggregate"
	"github.com/xfrr/go-cqrsify/aggregate/event"
)

func TestRestoreStateFromHistory(t *testing.T) {
	type args struct {
		agg    aggregate.Aggregate[string]
		events []aggregate.Event
	}

	type expected struct {
		aggID     string
		aggName   string
		aggVer    int
		aggEvents int
	}

	tests := []struct {
		name     string
		args     args
		expected expected
		err      *aggregate.HistoryIntegrityError
	}{
		{
			name: "should return an error if applying the events to the aggregate fails",
			args: args{
				agg: &mockAggregate{
					Base: aggregate.New("agg-1", "agg-test"),
				},
				events: makeEvents(t, "agg-2", "agg-test", 5),
			},
			err: &aggregate.HistoryIntegrityError{},
		},
		{
			name: "should restore the aggregate state from the given history successfully",
			args: args{
				agg: &mockAggregate{
					Base: aggregate.New("agg-1", "agg-test"),
				},
				events: makeEvents(t, "agg-1", "agg-test", 5),
			},
			expected: expected{
				aggID:     "agg-1",
				aggName:   "agg-test",
				aggVer:    5,
				aggEvents: 0,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := aggregate.RestoreStateFromHistory(tt.args.agg, tt.args.events)
			if tt.err != nil {
				if errors.As(err, &tt.err) {
					return
				}
				t.Fatalf("expected error to be %v, got %v", tt.err, err)
			}

			if tt.args.agg.AggregateVersion() != aggregate.Version(tt.expected.aggVer) {
				t.Fatalf("expected version to be %d, got %d", tt.expected.aggVer, tt.args.agg.AggregateVersion())
			}
		})
	}
}

func makeEvents(t *testing.T, aggID, aggName string, n int) []aggregate.Event {
	events := []aggregate.Event{}
	for i := range make([]int, n) {
		evt, err := event.New(
			randomStr(),
			randomStr(),
			n,
			event.WithAggregate(aggID, aggName, i+1),
		)
		if err != nil {
			t.Fatalf("expected error to be nil, got %v", err)
		}

		events = append(events, evt.Any())
	}
	return events
}

type mockAggregate struct {
	*aggregate.Base[string]
}

func randomStr() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}
