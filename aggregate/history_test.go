package aggregate_test

import (
	"crypto/rand"
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/xfrr/go-cqrsify/aggregate"
)

func TestRestoreStateFromHistory_ValidHistory(t *testing.T) {
	agg := aggregate.New("agg-1", "agg-test")
	events := makeEvents(t, agg.AggregateID(), agg.AggregateName(), 5)

	err := aggregate.RestoreFromHistory(agg, events)
	require.NoError(t, err)
}

func makeEvents(t *testing.T, aggID, aggName string, n int) []aggregate.Event {
	events := []aggregate.Event{}
	agg := aggregate.New(aggID, aggName)
	defer agg.ClearEvents()

	for range make([]int, n) {
		evt := aggregate.NewEvent(
			randomStr(),
			aggregate.CreateEventAggregateRef(agg),
		)

		agg.RecordEvent(evt)
		events = append(events, evt)
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
