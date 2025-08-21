package aggregate_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/xfrr/go-cqrsify/domain/aggregate"
)

func TestNextEvent(t *testing.T) {
	const (
		aggID     = "agg-id"
		aggName   = "TestAggregate"
		eventID   = "evt-1"
		eventName = "test"
	)

	type mockAggregate struct {
		*aggregate.Base[string]
	}

	agg := &mockAggregate{
		Base: aggregate.New(aggID, aggName),
	}

	handlerCalls := 0
	agg.HandleEvent(eventName, func(_ aggregate.Event) error {
		handlerCalls++
		return nil
	})

	t.Run("should return error if aggregate is nil", func(t *testing.T) {
		evt := aggregate.NewEvent(eventName, aggregate.CreateEventAggregateRef(agg))
		err := aggregate.NextEvent[string](
			nil,
			evt,
		)
		require.Error(t, err)
		require.True(t, errors.Is(err, aggregate.ErrNilAggregate))
	})

	t.Run("should create and record a new event successfully", func(t *testing.T) {
		evt := aggregate.NewEvent(eventName, aggregate.CreateEventAggregateRef(agg))
		err := aggregate.NextEvent(agg, evt)
		require.NoError(t, err)
		require.Len(t, agg.AggregateEvents(), 1)
		require.Equal(t, 1, handlerCalls)

		// apply second event
		evt2 := aggregate.NewEvent(eventName, aggregate.CreateEventAggregateRef(agg))
		err = aggregate.NextEvent(agg, evt2)
		require.NoError(t, err)
		require.Len(t, agg.AggregateEvents(), 2)
		require.Equal(t, 2, handlerCalls)
	})
}
