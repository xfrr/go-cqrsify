package aggregate_test

import (
	"crypto/rand"
	"fmt"
	"testing"

	"github.com/xfrr/cqrsify/aggregate"
	"github.com/xfrr/cqrsify/event"
)

func TestHydrate(t *testing.T) {
	type args struct {
		agg     aggregate.Aggregate[string]
		changes []aggregate.Change
	}

	type expected struct {
		aggID      string
		aggName    string
		aggVer     int
		aggChanges int
	}

	tests := []struct {
		name     string
		args     args
		expected expected
		err      error
	}{
		{
			name: "should return an error if changes cannot be applied",
			args: args{
				agg: &mockAggregate{
					Base: aggregate.New("agg-1", "test"),
				},
				changes: makeEvents("agg-2", "test", 5),
			},
			err: aggregate.ErrInvalidAggregateID,
		},
		{
			name: "should hydrate the aggregate with the given changes",
			args: args{
				agg: &mockAggregate{
					Base: aggregate.New("agg-1", "agg-test"),
				},
				changes: makeEvents("agg-1", "agg-test", 5),
			},
			expected: expected{
				aggID:      "agg-1",
				aggName:    "agg-test",
				aggVer:     5,
				aggChanges: 0,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := aggregate.Hydrate(tt.args.agg, tt.args.changes)
			if err != tt.err {
				t.Fatalf("expected error to be %v, got %v", tt.err, err)
			}

			if tt.args.agg.AggregateVersion() != aggregate.Version(tt.expected.aggVer) {
				t.Fatalf("expected version to be %d, got %d", tt.expected.aggVer, tt.args.agg.AggregateVersion())
			}
		})
	}
}

func makeEvents(aggID, aggName string, n int) []aggregate.Change {
	events := []aggregate.Change{}
	for i := 0; i < n; i++ {
		events = append(events,
			event.New(
				randomStr(),
				randomStr(),
				n,
				event.WithAggregate(aggID, aggName, i+1),
			).Any(),
		)
	}
	return events
}

type mockAggregate struct {
	*aggregate.Base[string]
}

func randomStr() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return fmt.Sprintf("%x", b)
}
