package inmemory

import (
	"context"
	"sync"

	"github.com/xfrr/go-cqrsify/domain"
)

var _ domain.Repository[string] = (*BaseAggregateRepository)(nil)
var _ domain.VersionedRepository[string] = (*BaseAggregateRepository)(nil)

type BaseAggregateRepository struct {
	mu              sync.RWMutex
	aggregates      []domain.Aggregate[string]
	aggregatesIndex map[string]domain.Aggregate[string]
}

func NewBaseAggregateRepository() *BaseAggregateRepository {
	return &BaseAggregateRepository{
		aggregates:      make([]domain.Aggregate[string], 0),
		aggregatesIndex: make(map[string]domain.Aggregate[string]),
	}
}

func (repo *BaseAggregateRepository) Exists(_ context.Context, agg domain.Aggregate[string]) (bool, error) {
	repo.mu.RLock()
	defer repo.mu.RUnlock()

	_, exists := repo.aggregatesIndex[agg.AggregateID()]
	return exists, nil
}

func (repo *BaseAggregateRepository) Save(_ context.Context, agg domain.Aggregate[string]) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	repo.aggregates = append(repo.aggregates, agg)
	repo.aggregatesIndex[agg.AggregateID()] = agg
	return nil
}

// Delete removes an aggregate by its instance.
func (repo *BaseAggregateRepository) Delete(_ context.Context, agg domain.Aggregate[string]) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	id := agg.AggregateID()

	// Remove from index
	delete(repo.aggregatesIndex, id)

	// Remove from slice
	for i, a := range repo.aggregates {
		if a.AggregateID() == id {
			repo.aggregates = append(repo.aggregates[:i], repo.aggregates[i+1:]...)
			break
		}
	}
	return nil
}

// Load retrieves an aggregate by its ID.
func (repo *BaseAggregateRepository) Load(_ context.Context, agg domain.Aggregate[string]) error {
	repo.mu.RLock()
	defer repo.mu.RUnlock()

	loadedAgg, exists := repo.aggregatesIndex[agg.AggregateID()]
	if !exists {
		return domain.NewNotFoundError(agg.AggregateID())
	}

	agg = loadedAgg
	return nil
}

// Search retrieves aggregates based on the provided search criteria.
func (repo *BaseAggregateRepository) Search(_ context.Context, criteria *domain.SearchCriteriaOptions) ([]domain.Aggregate[string], error) {
	repo.mu.RLock()
	defer repo.mu.RUnlock()

	var results []domain.Aggregate[string]
	for _, agg := range repo.aggregates {
		if criteria.Matches(agg) {
			results = append(results, agg)
		}
	}
	return results, nil
}

func (repo *BaseAggregateRepository) ExistsVersion(_ context.Context, agg domain.VersionedAggregate[string], version domain.AggregateVersion) (bool, error) {
	repo.mu.RLock()
	defer repo.mu.RUnlock()

	loadedAgg, exists := repo.aggregatesIndex[agg.AggregateID()]
	if !exists {
		return false, nil
	}

	versionedAgg, ok := loadedAgg.(domain.VersionedAggregate[string])
	if !ok {
		return false, nil
	}

	if versionedAgg.AggregateVersion() < version {
		return false, nil
	}

	return true, nil
}

func (repo *BaseAggregateRepository) LoadVersion(_ context.Context, agg domain.VersionedAggregate[string], version domain.AggregateVersion) error {
	repo.mu.RLock()
	defer repo.mu.RUnlock()

	loadedAgg, exists := repo.aggregatesIndex[agg.AggregateID()]
	if !exists {
		return domain.NewNotFoundError(agg.AggregateID())
	}

	versionedAgg, ok := loadedAgg.(domain.VersionedAggregate[string])
	if !ok {
		return domain.NewNotFoundError(agg.AggregateID())
	}

	if versionedAgg.AggregateVersion() < version {
		return domain.NewNotFoundError(agg.AggregateID())
	}

	return nil
}
