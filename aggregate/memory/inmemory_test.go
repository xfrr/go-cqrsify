package memory_test

import (
	"context"
	"fmt"
	"slices"
	"testing"

	"github.com/xfrr/cqrsify/aggregate"
	repository "github.com/xfrr/cqrsify/aggregate/memory"
	"github.com/xfrr/cqrsify/event"
)

func TestNewInMemory(t *testing.T) {
	sut := repository.NewRepository()
	if sut == nil {
		t.Fatal("NewInMemory() should not return nil")
	}

	if sut == nil {
		t.Fatal("NewInMemory() should not return nil")
	}
}

func TestInMemory_Delete(t *testing.T) {
	sut := repository.NewRepository()
	ctx := context.Background()

	agg := aggregate.New("1", "test")
	if agg == nil {
		t.Fatal("expected aggregate to not be nil")
	}

	// add and apply a change to the aggregate
	aggregate.NextChange(agg, "cid", "cname", "payload")

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
	sut := repository.NewRepository()
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
	sut := repository.NewRepository()
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
	sut := repository.NewRepository()
	ctx := context.Background()

	agg := aggregate.New("1", "test")
	if agg == nil {
		t.Fatal("expected aggregate to not be nil")
	}

	// add and apply a change to the aggregate
	aggregate.NextChange(agg, "cid", "cname", "payload")

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

	if len(foundAgg.AggregateChanges()) != 0 {
		t.Errorf("expected aggregate changes to be empty, got %d", len(foundAgg.AggregateChanges()))
	}
}

func TestInMemory_LoadVersion(t *testing.T) {
	sut := repository.NewRepository()
	ctx := context.Background()

	agg := aggregate.New("1", "test")
	if agg == nil {
		t.Fatal("expected aggregate to not be nil")
	}

	// add and apply a change to the aggregate
	aggregate.NextChange(agg, "cid", "cname", "payload")

	// save the aggregate
	if err := sut.Save(ctx, agg); err != nil {
		t.Errorf("Save() error = %v, want nil", err)
	}

	// load the aggregate
	if err := sut.LoadVersion(ctx, agg, 0); err != nil {
		t.Errorf("LoadVersion() error = %v, want nil", err)
	}

	// commit the changes to increase the version
	agg.CommitChanges()
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
		aggs, err := sut.Search(ctx)
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
		aggs, err := sut.Search(ctx, aggregate.WithSearchAggregateIDs(ids...))
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
		aggs, err := sut.Search(ctx, aggregate.WithSearchAggregateNames(names...))
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
		aggs, err := sut.Search(ctx, aggregate.WithSearchAggregateVersions(version))
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
	sut := repository.NewRepository()
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

func newRepositoryWithAggregates(t *testing.T, aggregates ...aggregate.Aggregate[string]) *repository.InMemory {
	repo := repository.NewRepository()
	for i, agg := range aggregates {
		evt := event.New(
			agg.AggregateID(),
			"test",
			testEventPayload{id: fmt.Sprintf("event-id-%d", i)},
			event.WithAggregate(agg.AggregateID(), agg.AggregateName(), 1),
		)

		commiter, ok := agg.(aggregate.ChangeCommitter)
		if !ok {
			t.Fatalf("aggregate does not implement ChangeCommitter")
		}
		commiter.RecordChange(evt.Any())

		err := repo.Save(context.Background(), agg)
		if err != nil {
			t.Fatalf("failed to save aggregate: %v", err)
		}
	}
	return repo
}

type testEventPayload struct {
	id string
}
