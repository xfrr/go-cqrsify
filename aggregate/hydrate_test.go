package aggregate_test

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/xfrr/cqrsify/aggregate"
	"github.com/xfrr/cqrsify/event"
)

func TestApplyHistory(t *testing.T) {
	type mockAggregate struct {
		*aggregate.Base
	}

	type args struct {
		changes []aggregate.Change
	}

	tests := []struct {
		name    string
		args    args
		aggID   string
		aggName string
		aggVer  int
		err     error
	}{
		{
			name: "it should apply all changes to the aggregate",
			args: args{
				changes: makeEvents("mock-id", "mock-name", 5),
			},
			aggID:   "mock-id",
			aggName: "mock-name",
			aggVer:  5,
			err:     nil,
		},
		{
			name:    "it should return error when the change payload is nil",
			aggID:   "mock-id",
			aggName: "mock-name",
			aggVer:  0,
			args: args{
				changes: []aggregate.Change{
					event.New[any]("mock-id", "mock-name", nil,
						event.WithAggregate(event.Aggregate{
							ID:      "mock-id",
							Name:    "mock-name",
							Version: 1,
						}),
					).Any(),
				},
			},
			err: aggregate.ErrInvalidChangePayload,
		},
		{
			name:    "it should return an error if the aggregate ID does not match",
			aggID:   "mock-id",
			aggName: "mock-name",
			aggVer:  0,
			args: args{
				changes: []aggregate.Change{
					event.New("mock-id", "mock-name", 1).Any(),
					event.New("mock-id-2", "mock-name", 2).Any(),
				},
			},
			err: aggregate.ErrInvalidAggregateID,
		},
		{
			name:    "it should return an error if the aggregate name does not match",
			aggID:   "mock-id",
			aggName: "mock-name",
			aggVer:  0,
			args: args{
				changes: []aggregate.Change{
					makeEvents("mock-id", "mock-name", 1)[0],
					makeEvents("mock-id", "mock-name-2", 1)[0],
				},
			},
			err: aggregate.ErrInvalidAggregateName,
		},
		{
			name:    "it should return an error if the version is not consecutive",
			aggID:   "mock-id",
			aggName: "mock-name",
			aggVer:  0,
			args: args{
				changes: []aggregate.Change{
					makeEvents("mock-id", "mock-name", 1)[0],
					makeEvents("mock-id", "mock-name", 1)[0],
				},
			},
			err: aggregate.ErrInvalidVersion,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			agg := &mockAggregate{
				Base: aggregate.New(tt.aggID, tt.aggName),
			}

			err := aggregate.Hydrate(agg, tt.args.changes)
			if err != tt.err {
				t.Errorf("expected error to be %v, got %v", tt.err, err)
			}

			if agg.AggregateVersion() != aggregate.Version(tt.aggVer) {
				t.Errorf("expected version to be %d, got %d", tt.aggVer, agg.AggregateVersion())
			}
		})
	}
}

func randomStr() string {
	rnd := rand.Intn(1000)
	return fmt.Sprint(rnd)
}

func makeEvents(aggID, aggName string, n int) []aggregate.Change {
	events := []aggregate.Change{}
	for i := 0; i < n; i++ {
		events = append(events,
			event.New(
				randomStr(),
				randomStr(),
				n,
				event.WithAggregate(event.Aggregate{
					ID:      aggID,
					Name:    aggName,
					Version: i + 1,
				}),
			).Any(),
		)
	}
	return events
}
