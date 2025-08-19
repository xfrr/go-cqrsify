package aggregate_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xfrr/go-cqrsify/aggregate"
)

func TestBase(t *testing.T) {
	t.Run("it should create a new base aggregate", func(t *testing.T) {
		base := aggregate.New("test-id", "test-name")
		require.NotNil(t, base)

		assert.Equal(t, "test-id", base.AggregateID())
		assert.Equal(t, "test-name", base.AggregateName())
		assert.Empty(t, base.AggregateEvents())
		assert.Equal(t, aggregate.Version(0), base.AggregateVersion())
		assert.NotNil(t, base.Any())
	})

	t.Run("it should record a event", func(t *testing.T) {
		baseAggregate := aggregate.New("test-id", "test-name")
		evt := aggregate.NewEvent("test.name", aggregate.CreateEventAggregateRef(baseAggregate))

		baseAggregate.RecordEvent(evt)
		assert.Len(t, baseAggregate.AggregateEvents(), 1)
	})

	t.Run("it should commit events", func(t *testing.T) {
		baseAggregate := aggregate.New("test-id", "test-name")
		evt := aggregate.NewEvent("test.name", aggregate.CreateEventAggregateRef(baseAggregate))

		baseAggregate.RecordEvent(evt)
		require.Len(t, baseAggregate.AggregateEvents(), 1)

		baseAggregate.CommitEvents()
		require.Len(t, baseAggregate.AggregateEvents(), 0)
		require.Equal(t, aggregate.Version(1), baseAggregate.AggregateVersion())
	})

	t.Run("it should apply events", func(t *testing.T) {
		baseAggregate := aggregate.New("test-id", "test-name")
		evt := aggregate.NewEvent("test.name", aggregate.CreateEventAggregateRef(baseAggregate))

		handlerCalls := 0
		baseAggregate.HandleEvent("test.name", func(_ aggregate.Event) error {
			handlerCalls++
			return nil
		})

		baseAggregate.ApplyEvent(evt)
		assert.Equal(t, 1, handlerCalls)
	})

	t.Run("it should return the aggregate's id", func(t *testing.T) {
		base := aggregate.New("test-id", "test-name")
		assert.Equal(t, "test-id", base.AggregateID())
	})

	t.Run("it should return the aggregate's name", func(t *testing.T) {
		base := aggregate.New("test-id", "test-name")
		assert.Equal(t, "test-name", base.AggregateName())
	})

	t.Run("it should return the aggregate's name", func(t *testing.T) {
		base := aggregate.New("test-id", "test-name")
		assert.Equal(t, "test-name", base.AggregateName())
	})

	t.Run("it should return the aggregate's version", func(t *testing.T) {
		base := aggregate.New("test-id", "test-name")
		assert.Equal(t, aggregate.Version(0), base.AggregateVersion())
	})
}
