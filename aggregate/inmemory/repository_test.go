package inmemory_test

import (
	"context"
	"fmt"
	"slices"
	"testing"

	"github.com/xfrr/go-cqrsify/aggregate"
	"github.com/xfrr/go-cqrsify/aggregate/event"

	inmemory "github.com/xfrr/go-cqrsify/aggregate/inmemory"
)

func TestNewRepository(t *testing.T) {
	sut := inmemory.NewAggregateRepository()
	if sut == nil {
		t.Fatal("NewInMemory() should not return nil")
	}

	if sut == nil {
		t.Fatal("NewInMemory() should not return nil")
	}
}

func TestInMemory_Delete(t *testing.T) {
	sut := inmemory.NewAggregateRepository()
	ctx := context.Background()

	agg := aggregate.New("1", "test")
	if agg == nil {
		t.Fatal("expected aggregate to not be nil")
	}

	// add and apply a event to the aggregate
	aggregate.RaiseEvent(agg, "cid", "cname", "payload")

	// save the aggregate
	if err := sut.Save(ctx, agg); err != nil {
		t.Errorf("Save() error = %v, want nil", err)
	}

	// delete the aggregate
	if err := sut.Delete(ctx, agg); err != nil {
		t.Errorf("Delete() error = %v, want nil", err)
	}
}

func TestInMemory_Exists(t *testing.T) {
	sut := inmemory.NewAggregateRepository()
	ctx := context.Background()

	agg := aggregate.New("1", "test")
	if agg == nil {
		t.Fatal("expected aggregate to not be nil")
	}

	// save the aggregate
	if err := sut.Save(ctx, agg); err != nil {
		t.Errorf("Save() error = %v, want nil", err)
	}

	// check if the aggregate exists
	if exists, err := sut.Exists(ctx, agg); err != nil {
		t.Errorf("Exists() error = %v, want nil", err)
	} else if !exists {
		t.Errorf("Exists() = %v, want %v", exists, true)
	}
}

func TestInMemory_ExistsVersion(t *testing.T) {
	sut := inmemory.NewAggregateRepository()
	ctx := context.Background()

	agg := aggregate.New("1", "test")
	if agg == nil {
		t.Fatal("expected aggregate to not be nil")
	}

	// save the aggregate
	if err := sut.Save(ctx, agg); err != nil {
		t.Errorf("Save() error = %v, want nil", err)
	}

	// check if the aggregate exists
	if exists, err := sut.ExistsVersion(ctx, agg, 0); err != nil {
		t.Errorf("ExistsVersion() error = %v, want nil", err)
	} else if !exists {
		t.Errorf("ExistsVersion() = %v, want %v", exists, true)
	}
}

func TestInMemory_Load(t *testing.T) {
	sut := inmemory.NewAggregateRepository()
	ctx := context.Background()

	agg := aggregate.New("1", "test")
	if agg == nil {
		t.Fatal("expected aggregate to not be nil")
	}

	// add and apply a event to the aggregate
	aggregate.RaiseEvent(agg, "cid", "cname", "payload")

	// save the aggregate
	if err := sut.Save(ctx, agg); err != nil {
		t.Errorf("Save() error = %v, want nil", err)
	}

	// set new aggregate to avoid version conflicts
	foundAgg := aggregate.New("1", "test")

	// load the aggregate
	if err := sut.Load(ctx, foundAgg); err != nil {
		t.Errorf("Load() error = %v, want nil", err)
	}

	if foundAgg.AggregateID() != agg.AggregateID() {
		t.Errorf("expected aggregate id to be %s, got %s", agg.AggregateID(), foundAgg.AggregateID())
	}

	if foundAgg.AggregateVersion() != 1 {
		t.Errorf("expected aggregate version to be 1, got %d", foundAgg.AggregateVersion())
	}

	if len(foundAgg.AggregateEvents()) != 0 {
		t.Errorf("expected aggregate events to be empty, got %d", len(foundAgg.AggregateEvents()))
	}
}

func TestInMemory_LoadVersion(t *testing.T) {
	sut := inmemory.NewAggregateRepository()
	ctx := context.Background()

	agg := aggregate.New("1", "test")
	if agg == nil {
		t.Fatal("expected aggregate to not be nil")
	}

	// add and apply a event to the aggregate
	aggregate.RaiseEvent(agg, "cid", "cname", "payload")

	// save the aggregate
	if err := sut.Save(ctx, agg); err != nil {
		t.Errorf("Save() error = %v, want nil", err)
	}

	// load the aggregate
	if err := sut.LoadVersion(ctx, agg, 0); err != nil {
		t.Errorf("LoadVersion() error = %v, want nil", err)
	}

	// commit the events to increase the version
	agg.CommitEvents()
	if agg.AggregateVersion() != 1 {
		t.Errorf("expected aggregate version to be 1, got %d", agg.AggregateVersion())
	}
}

func TestInMemory_Search(t *testing.T) {
	ctx := context.Background()

	sut := newRepositoryWithAggregates(
		t,
		aggregate.New("id-1", "test-1"),
		aggregate.New("id-2", "test-2"),
		aggregate.New("id-3", "test-3"),
	)

	t.Run("should return all aggregates when no criteria is provided", func(t *testing.T) {
		aggs, err := sut.Search(ctx, aggregate.SearchCriteria())
		if err != nil {
			t.Fatalf("Search() error = %v, want nil", err)
		}

		expectedN := 3
		if len(aggs) != expectedN {
			t.Fatalf("expected %d aggregates, got %d", expectedN, len(aggs))
		}
	})

	t.Run("should return aggregates with provided ids", func(t *testing.T) {
		ids := []string{"id-1", "id-2"}
		aggs, err := sut.Search(ctx, aggregate.SearchCriteria().WithSearchAggregateIDs(ids...))
		if err != nil {
			t.Fatalf("Search() error = %v, want nil", err)
		}

		expectedN := 2
		if len(aggs) != expectedN {
			t.Fatalf("expected %d aggregates, got %d", expectedN, len(aggs))
		}

		expected := ids
		for i, agg := range aggs {
			if !slices.Contains(expected, agg.AggregateID()) {
				t.Errorf("expected aggregate id to be %s, got %s", expected[i], agg.AggregateID())
			}
		}
	})

	t.Run("should return aggregates with provided names", func(t *testing.T) {
		names := []string{"test-1", "test-2"}
		aggs, err := sut.Search(ctx, aggregate.SearchCriteria().WithSearchAggregateNames(names...))
		if err != nil {
			t.Fatalf("Search() error = %v, want nil", err)
		}

		expectedN := 2
		if len(aggs) != expectedN {
			t.Fatalf("expected %d aggregates, got %d", expectedN, len(aggs))
		}

		expected := names
		for i, agg := range aggs {
			if !slices.Contains(expected, agg.AggregateName()) {
				t.Errorf("expected aggregate name to be %s, got %s", expected[i], agg.AggregateName())
			}
		}
	})

	t.Run("should return aggregates with provided versions", func(t *testing.T) {
		version := 1
		aggs, err := sut.Search(ctx, aggregate.SearchCriteria().WithSearchAggregateVersions(version))
		if err != nil {
			t.Fatalf("Search() error = %v, want nil", err)
		}

		expectedN := 3
		if len(aggs) != expectedN {
			t.Fatalf("expected %d aggregates, got %d", expectedN, len(aggs))
		}

		expectedV := version
		for _, agg := range aggs {
			if agg.AggregateVersion() != aggregate.Version(expectedV) {
				t.Errorf("expected aggregate version to be %d, got %d", expectedV, agg.AggregateVersion())
			}
		}
	})

}

func TestInMemory_Save(t *testing.T) {
	sut := inmemory.NewAggregateRepository()
	ctx := context.Background()

	agg := aggregate.New("1", "test")
	if agg == nil {
		t.Fatal("expected aggregate to not be nil")
	}

	// save the aggregate
	if err := sut.Save(ctx, agg); err != nil {
		t.Errorf("Save() error = %v, want nil", err)
	}
}

func newRepositoryWithAggregates(t *testing.T, aggregates ...aggregate.Aggregate[string]) *inmemory.AggregateRepository {
	repo := inmemory.NewAggregateRepository()
	for i, agg := range aggregates {
		evt, err := event.New(
			agg.AggregateID(),
			"test",
			testEventPayload{id: fmt.Sprintf("event-id-%d", i)},
			event.WithAggregate(agg.AggregateID(), agg.AggregateName(), 1),
		)
		if err != nil {
			t.Fatalf("failed to create event: %v", err)
		}

		eventRecorder, ok := agg.(aggregate.EventRecorder)
		if !ok {
			t.Fatalf("aggregate does not implement EventRecorder")
		}

		eventRecorder.RecordEvent(evt.Any())

		err = repo.Save(context.Background(), agg)
		if err != nil {
			t.Fatalf("failed to save aggregate: %v", err)
		}
	}
	return repo
}

type testEventPayload struct {
	id string
}
