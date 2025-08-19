package inmemory_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xfrr/go-cqrsify/aggregate"

	inmemory "github.com/xfrr/go-cqrsify/aggregate/inmemory"
)

func TestNewRepository(t *testing.T) {
	sut := inmemory.NewEventSourcedAggregateRepository()
	require.NotNil(t, sut, "NewInMemory() should not return nil")
}

func TestInMemory_Delete(t *testing.T) {
	sut := inmemory.NewEventSourcedAggregateRepository()
	ctx := context.Background()

	agg := aggregate.New("1", "test")
	require.NotNil(t, agg, "expected aggregate to not be nil")

	// add and apply a event to the aggregate
	evt := aggregate.NewEvent("cname", aggregate.CreateEventAggregateRef(agg))
	aggregate.NextEvent(agg, evt)

	// save the aggregate
	require.NoError(t, sut.Save(ctx, agg), "Save() error = %v, want nil")

	// delete the aggregate
	require.NoError(t, sut.Delete(ctx, agg), "Delete() error = %v, want nil")
}

func TestInMemory_Exists(t *testing.T) {
	sut := inmemory.NewEventSourcedAggregateRepository()
	ctx := context.Background()

	agg := aggregate.New("1", "test")
	require.NotNil(t, agg, "expected aggregate to not be nil")

	// save the aggregate
	require.NoError(t, sut.Save(ctx, agg), "Save() error = %v, want nil")

	// check if the aggregate exists
	exists, err := sut.Exists(ctx, agg)
	require.NoError(t, err, "Exists() error = %v, want nil", err)
	require.True(t, exists, "Exists() = %v, want %v", exists, true)
}

func TestInMemory_ExistsVersion(t *testing.T) {
	sut := inmemory.NewEventSourcedAggregateRepository()
	ctx := context.Background()

	agg := aggregate.New("1", "test")
	require.NotNil(t, agg, "expected aggregate to not be nil")

	// save the aggregate
	err := sut.Save(ctx, agg)
	require.NoError(t, err, "Save() error = %v, want nil", err)

	// check if the aggregate exists
	exists, err := sut.ExistsVersion(ctx, agg, 0)
	require.NoError(t, err, "ExistsVersion() error = %v, want nil", err)
	require.True(t, exists, "ExistsVersion() = %v, want %v", exists, true)
}

func TestInMemory_Load(t *testing.T) {
	sut := inmemory.NewEventSourcedAggregateRepository()
	ctx := context.Background()

	agg := aggregate.New("1", "test")
	require.NotNil(t, agg, "expected aggregate to not be nil")

	// add and apply a event to the aggregate
	evt := aggregate.NewEvent("cname", aggregate.CreateEventAggregateRef(agg))
	aggregate.NextEvent(agg, evt)

	// save the aggregate
	err := sut.Save(ctx, agg)
	require.NoError(t, err, "Save() error = %v, want nil", err)

	// set new aggregate to avoid version conflicts
	expectedAgg := aggregate.New("1", "test")

	// load the aggregate
	err = sut.Load(ctx, expectedAgg)
	require.NoError(t, err, "Load() error = %v, want nil", err)
	assert.Equal(t, agg.AggregateID(), expectedAgg.AggregateID(), "expected aggregate id to be %s, got %s", agg.AggregateID(), expectedAgg.AggregateID())
	assert.Equal(t, aggregate.Version(1), expectedAgg.AggregateVersion(), "expected aggregate version to be 1, got %d", expectedAgg.AggregateVersion())
	assert.Empty(t, expectedAgg.AggregateEvents(), "expected aggregate events to be empty, got %d", len(expectedAgg.AggregateEvents()))
}

func TestInMemory_LoadVersion(t *testing.T) {
	sut := inmemory.NewEventSourcedAggregateRepository()
	ctx := context.Background()

	agg := aggregate.New("1", "test")
	require.NotNil(t, agg, "expected aggregate to not be nil")

	// add and apply a event to the aggregate
	evt := aggregate.NewEvent("cname", aggregate.CreateEventAggregateRef(agg))
	aggregate.NextEvent(agg, evt)

	// save the aggregate
	err := sut.Save(ctx, agg)
	require.NoError(t, err, "Save() error = %v, want nil", err)

	// load the aggregate
	agg2 := aggregate.New("1", "test")
	err = sut.LoadVersion(ctx, agg2, 1)
	require.NoError(t, err, "LoadVersion() error = %v, want nil", err)
	assert.Equal(t, agg.AggregateID(), agg2.AggregateID(), "expected aggregate id to be %s, got %s", agg.AggregateID(), agg2.AggregateID())
	assert.Equal(t, aggregate.Version(1), agg2.AggregateVersion(), "expected aggregate version to be 1, got %d", agg2.AggregateVersion())
	assert.Empty(t, agg2.AggregateEvents(), "expected aggregate events to be empty, got %d", len(agg2.AggregateEvents()))
	assert.Equal(t, agg.AggregateName(), agg2.AggregateName(), "expected aggregate name to be %s, got %s", agg.AggregateName(), agg2.AggregateName())
}

func TestInMemory_Search(t *testing.T) {
	ctx := context.Background()

	sut := newEventSourcedAggregateRepositoryWithAggregates(
		t,
		aggregate.New("id-1", "test-1"),
		aggregate.New("id-2", "test-2"),
		aggregate.New("id-3", "test-3"),
	)

	t.Run("should return all aggregates when no criteria is provided", func(t *testing.T) {
		aggs, err := sut.Search(ctx, aggregate.SearchCriteria())
		require.NoError(t, err, "Search() error = %v, want nil", err)

		expectedN := 3
		require.Equal(t, expectedN, len(aggs), "expected %d aggregates, got %d", expectedN, len(aggs))
	})

	t.Run("should return aggregates with provided ids", func(t *testing.T) {
		ids := []string{"id-1", "id-2"}
		aggs, err := sut.Search(ctx, aggregate.SearchCriteria().WithSearchAggregateIDs(ids...))
		require.NoError(t, err, "Search() error = %v, want nil", err)

		expectedN := 2
		require.Equal(t, expectedN, len(aggs), "expected %d aggregates, got %d", expectedN, len(aggs))

		expected := ids
		for i, agg := range aggs {
			assert.Contains(t, expected, agg.AggregateID(), "expected aggregate id to be %s, got %s", expected[i], agg.AggregateID())
		}
	})

	t.Run("should return aggregates with provided names", func(t *testing.T) {
		names := []string{"test-1", "test-2"}
		aggs, err := sut.Search(ctx, aggregate.SearchCriteria().WithSearchAggregateNames(names...))
		require.NoError(t, err, "Search() error = %v, want nil", err)

		expectedN := 2
		require.Equal(t, expectedN, len(aggs), "expected %d aggregates, got %d", expectedN, len(aggs))

		expected := names
		for i, agg := range aggs {
			assert.Contains(t, expected, agg.AggregateName(), "expected aggregate name to be %s, got %s", expected[i], agg.AggregateName())
		}
	})

	t.Run("should return aggregates with provided versions", func(t *testing.T) {
		version := 1
		aggs, err := sut.Search(ctx, aggregate.SearchCriteria().WithSearchAggregateVersions(version))
		require.NoError(t, err, "Search() error = %v, want nil", err)

		expectedN := 3
		require.Equal(t, expectedN, len(aggs), "expected %d aggregates, got %d", expectedN, len(aggs))

		expectedV := version
		for _, agg := range aggs {
			versionedAgg, ok := agg.(aggregate.VersionedAggregate[string])
			require.True(t, ok, "expected aggregate to be versioned, got %T", agg)
			require.Equal(t, aggregate.Version(expectedV), versionedAgg.AggregateVersion(), "expected aggregate version to be %d, got %d", expectedV, versionedAgg.AggregateVersion())
		}
	})

}

func TestInMemory_Save(t *testing.T) {
	sut := inmemory.NewEventSourcedAggregateRepository()
	ctx := context.Background()

	agg := aggregate.New("1", "test")
	require.NotNil(t, agg, "expected aggregate to not be nil")

	// save the aggregate
	err := sut.Save(ctx, agg)
	require.NoError(t, err, "Save() error = %v, want nil", err)
}

func newEventSourcedAggregateRepositoryWithAggregates(t *testing.T, aggregates ...aggregate.EventSourcedAggregate[string]) *inmemory.EventSourcedAggregateRepository {
	repo := inmemory.NewEventSourcedAggregateRepository()
	for _, agg := range aggregates {
		evt := aggregate.NewEvent(
			"test.name",
			aggregate.CreateEventAggregateRef(agg),
		)

		err := aggregate.NextEvent(agg, evt)
		require.NoError(t, err, "failed to raise event: %v", err)

		err = repo.Save(context.Background(), agg)
		require.NoError(t, err, "failed to save aggregate: %v", err)
	}
	return repo
}

type testEventPayload struct {
	id string
}
