package domain_test

import (
	"crypto/rand"
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/xfrr/go-cqrsify/domain"
)

func TestRestoreStateFromHistory_ValidHistory(t *testing.T) {
	agg := domain.NewAggregate("agg-1", "agg-test")
	events := makeEvents(agg.AggregateID(), agg.AggregateName(), 5)

	err := domain.RestoreAggregateFromHistory(agg, events)
	require.NoError(t, err)
}

func makeEvents(aggID, aggName string, n int) []domain.Event {
	events := []domain.Event{}
	agg := domain.NewAggregate(aggID, aggName)
	defer agg.ClearEvents()

	for range make([]int, n) {
		evt := domain.NewEvent(
			randomStr(),
			domain.CreateEventAggregateRef(agg),
		)

		agg.RecordEvent(evt)
		events = append(events, evt)
	}

	return events
}

type mockAggregate struct {
	*domain.BaseAggregate[string]
}

func randomStr() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}
