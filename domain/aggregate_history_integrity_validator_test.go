package domain_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/xfrr/go-cqrsify/domain"
)

func TestVerifyHistoryIntegrity_ValidHistory(t *testing.T) {
	agg := domain.NewAggregate("agg-test-id", "agg-test-name")
	events := makeEvents(t, agg.AggregateID(), agg.AggregateName(), 5)
	agg.ClearEvents()

	err := domain.VerifyHistoryIntegrity(agg, events)
	require.NoError(t, err)
}

func TestVerifyHistoryIntegrity_InvalidHistory_InvalidAggregateID(t *testing.T) {
	agg := domain.NewAggregate("agg-test-id", "agg-test-name")
	events := makeEvents(t, agg.AggregateID(), agg.AggregateName(), 5)
	agg.ClearEvents()

	// Modify the events to have an invalid aggregate ID
	for i := range events {
		events[i] = domain.NewEvent("event-1", domain.CreateEventAggregateRef(domain.NewAggregate("invalid-agg-id", "agg-test-name")))
	}

	err := domain.VerifyHistoryIntegrity(agg, events)
	require.Error(t, err)
}

func TestVerifyHistoryIntegrity_InvalidHistory_InvalidAggregateName(t *testing.T) {
	agg := domain.NewAggregate("agg-test-id", "agg-test-name")
	events := makeEvents(t, agg.AggregateID(), agg.AggregateName(), 5)
	agg.ClearEvents()

	// Modify the events to have an invalid aggregate name
	for i := range events {
		events[i] = domain.NewEvent("event-1", domain.CreateEventAggregateRef(domain.NewAggregate("agg-test-id", "invalid-agg-name")))
	}

	err := domain.VerifyHistoryIntegrity(agg, events)
	require.Error(t, err)
}

func TestVerifyHistoryIntegrity_InvalidHistory_UnexpectedVersion(t *testing.T) {
	agg := domain.NewAggregate("agg-test-id", "agg-test-name")
	agg.RecordEvent(domain.NewEvent("event-1", domain.CreateEventAggregateRef(agg)))
	events := makeEvents(t, agg.AggregateID(), agg.AggregateName(), 5)

	expectedError := domain.NewHistoryIntegrityError("history integrity error")
	err := domain.VerifyHistoryIntegrity(agg, events)
	require.ErrorAs(t, err, &expectedError)
}
