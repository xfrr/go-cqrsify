package memory

import (
	"context"
	"errors"
	"sync"

	"github.com/xfrr/cqrsify/aggregate"
)

var _ aggregate.VersionedRepository[string] = (*InMemory)(nil)

var (
	ErrInvalidAggregateChangeID = errors.New("invalid aggregate change id")
)

type aggregateDTO struct {
	id      string
	name    string
	version int
	changes []aggregate.Change
}

type InMemory struct {
	mu sync.RWMutex

	indexedChanges map[string]*aggregateDTO
	changes        []aggregate.Change
}

func NewRepository() *InMemory {
	return &InMemory{
		indexedChanges: make(map[string]*aggregateDTO),
		changes:        make([]aggregate.Change, 0),
	}
}

func (repo *InMemory) Exists(_ context.Context, agg aggregate.Aggregate[string]) (bool, error) {
	repo.mu.RLock()
	defer repo.mu.RUnlock()

	_, ok := repo.indexedChanges[agg.AggregateID()]
	return ok, nil
}

func (repo *InMemory) ExistsVersion(_ context.Context, agg aggregate.Aggregate[string], version aggregate.Version) (bool, error) {
	repo.mu.RLock()
	defer repo.mu.RUnlock()

	dto, ok := repo.indexedChanges[agg.AggregateID()]
	if !ok {
		return false, nil
	}

	if dto.version < int(version) {
		return false, nil
	}

	return true, nil
}

func (repo *InMemory) Load(_ context.Context, agg aggregate.Aggregate[string]) error {
	repo.mu.RLock()
	defer repo.mu.RUnlock()

	dto, ok := repo.indexedChanges[agg.AggregateID()]
	if !ok {
		return aggregate.ErrAggregateNotFound
	}

	return aggregate.Hydrate(agg, dto.changes)
}

func (repo *InMemory) LoadVersion(_ context.Context, agg aggregate.Aggregate[string], version aggregate.Version) error {
	repo.mu.RLock()
	defer repo.mu.RUnlock()

	dto, ok := repo.indexedChanges[agg.AggregateID()]
	if !ok {
		return aggregate.ErrAggregateNotFound
	}

	if dto.version < int(version) {
		return aggregate.ErrAggregateNotFound
	}

	return aggregate.Hydrate(agg, filterChangesFromVersion(version, dto.changes))
}

func (repo *InMemory) Save(_ context.Context, agg aggregate.Aggregate[string]) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	dto, ok := repo.indexedChanges[agg.AggregateID()]
	if !ok {
		dto = &aggregateDTO{
			id:      agg.AggregateID(),
			name:    agg.AggregateName(),
			version: 0,
			changes: make([]aggregate.Change, 0),
		}
	}

	// copy the aggregate changes
	changes := make([]aggregate.Change, len(agg.AggregateChanges()))
	copy(changes, agg.AggregateChanges())
	if c, ok := agg.(aggregate.ChangeCommitter); ok {
		c.CommitChanges()
	}

	// update the aggregate version
	dto.version = int(agg.AggregateVersion())

	// append the changes to the aggregate
	dto.changes = append(dto.changes, changes...)
	repo.indexedChanges[agg.AggregateID()] = dto

	// append the changes to the repository
	repo.changes = append(repo.changes, changes...)

	return nil
}

func (repo *InMemory) Search(_ context.Context, opts ...aggregate.SearchOption) ([]aggregate.Aggregate[string], error) {
	repo.mu.RLock()
	defer repo.mu.RUnlock()

	changes := repo.changes
	if len(changes) == 0 {
		return nil, nil
	}

	criteria := aggregate.NewSearchCriteria(opts...)
	if criteria == nil {
		return getAggregates(changes)
	}

	changes = filterChangesByAggregateIDs(criteria.AggregateIDs(), changes)
	changes = filterChangesByAggregateNames(criteria.AggregateNames(), changes)
	changes = filterChangesByAggregateVersions(criteria.AggregateVersions(), changes)
	return getAggregates(changes)
}

func (repo *InMemory) Delete(_ context.Context, agg aggregate.Aggregate[string]) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	repo.deleteAggregate(agg.AggregateID())
	return nil
}

func (repo *InMemory) deleteAggregate(aggID string) {
	for i, ev := range repo.changes {
		aggregateRef := ev.Aggregate()
		if aggregateRef == nil {
			continue
		}

		if aggregateRef.ID == aggID {
			repo.changes = append(repo.changes[:i], repo.changes[i+1:]...)
		}
	}

	delete(repo.indexedChanges, aggID)
}

func getAggregates(changes []aggregate.Change) ([]aggregate.Aggregate[string], error) {
	var (
		aggregates = make([]aggregate.Aggregate[string], 0)
		aggMap     = make(map[string]aggregate.Aggregate[string])
	)

	for _, ev := range changes {
		aggregateRef := ev.Aggregate()
		if aggregateRef == nil {
			continue
		}

		aggID, ok := aggregateRef.ID.(string)
		if !ok {
			return nil, ErrInvalidAggregateChangeID
		}

		agg, ok := aggMap[aggID]
		if !ok {
			agg = aggregate.New(aggID, aggregateRef.Name)
		}

		err := aggregate.Hydrate(agg, []aggregate.Change{ev})
		if err != nil {
			return nil, err
		}

		aggMap[aggID] = agg
	}

	for _, agg := range aggMap {
		aggregates = append(aggregates, agg)
	}

	return aggregates, nil
}
