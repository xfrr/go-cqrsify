package aggregate_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/xfrr/go-cqrsify/domain/aggregate"
)

func TestVerifyHistoryIntegrity_ValidHistory(t *testing.T) {
	agg := aggregate.New("agg-test-id", "agg-test-name")
	events := makeEvents(t, agg.AggregateID(), agg.AggregateName(), 5)
	agg.ClearEvents()

	err := aggregate.VerifyHistoryIntegrity(agg, events)
	require.NoError(t, err)
}

func TestVerifyHistoryIntegrity_InvalidHistory_InvalidAggregateID(t *testing.T) {
	agg := aggregate.New("agg-test-id", "agg-test-name")
	events := makeEvents(t, agg.AggregateID(), agg.AggregateName(), 5)
	agg.ClearEvents()

	// Modify the events to have an invalid aggregate ID
	for i := range events {
		events[i] = aggregate.NewEvent("event-1", aggregate.CreateEventAggregateRef(aggregate.New("invalid-agg-id", "agg-test-name")))
	}

	err := aggregate.VerifyHistoryIntegrity(agg, events)
	require.Error(t, err)
}

func TestVerifyHistoryIntegrity_InvalidHistory_InvalidAggregateName(t *testing.T) {
	agg := aggregate.New("agg-test-id", "agg-test-name")
	events := makeEvents(t, agg.AggregateID(), agg.AggregateName(), 5)
	agg.ClearEvents()

	// Modify the events to have an invalid aggregate name
	for i := range events {
		events[i] = aggregate.NewEvent("event-1", aggregate.CreateEventAggregateRef(aggregate.New("agg-test-id", "invalid-agg-name")))
	}

	err := aggregate.VerifyHistoryIntegrity(agg, events)
	require.Error(t, err)
}

func TestVerifyHistoryIntegrity_InvalidHistory_UnexpectedVersion(t *testing.T) {
	agg := aggregate.New("agg-test-id", "agg-test-name")
	agg.RecordEvent(aggregate.NewEvent("event-1", aggregate.CreateEventAggregateRef(agg)))
	events := makeEvents(t, agg.AggregateID(), agg.AggregateName(), 5)

	expectedError := aggregate.NewHistoryIntegrityError("history integrity error")
	err := aggregate.VerifyHistoryIntegrity(agg, events)
	require.ErrorAs(t, err, &expectedError)
}
