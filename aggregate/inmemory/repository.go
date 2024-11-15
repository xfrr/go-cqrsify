package inmemory

import (
	"context"
	"errors"
	"sync"

	"github.com/xfrr/go-cqrsify/aggregate"
)

var _ aggregate.VersionedRepository[string] = (*AggregateRepository)(nil)

var (
	ErrInvalidAggregateEventID = errors.New("invalid aggregate event id")
)

type aggregateDTO struct {
	id      string
	name    string
	version int
	events  []aggregate.Event
}

type AggregateRepository struct {
	mu sync.RWMutex

	indexedEvents map[string]*aggregateDTO
	events        []aggregate.Event
}

func NewAggregateRepository() *AggregateRepository {
	return &AggregateRepository{
		mu:            sync.RWMutex{},
		indexedEvents: make(map[string]*aggregateDTO),
		events:        make([]aggregate.Event, 0),
	}
}

func (repo *AggregateRepository) Exists(_ context.Context, agg aggregate.Aggregate[string]) (bool, error) {
	repo.mu.RLock()
	defer repo.mu.RUnlock()

	_, ok := repo.indexedEvents[agg.AggregateID()]
	return ok, nil
}

func (repo *AggregateRepository) ExistsVersion(_ context.Context, agg aggregate.Aggregate[string], version aggregate.Version) (bool, error) {
	repo.mu.RLock()
	defer repo.mu.RUnlock()

	dto, ok := repo.indexedEvents[agg.AggregateID()]
	if !ok {
		return false, nil
	}

	if dto.version < int(version) {
		return false, nil
	}

	return true, nil
}

func (repo *AggregateRepository) Load(_ context.Context, agg aggregate.Aggregate[string]) error {
	repo.mu.RLock()
	defer repo.mu.RUnlock()

	dto, ok := repo.indexedEvents[agg.AggregateID()]
	if !ok {
		return aggregate.ErrNotFound
	}

	return aggregate.RestoreStateFromHistory(agg, dto.events)
}

func (repo *AggregateRepository) LoadVersion(_ context.Context, agg aggregate.Aggregate[string], version aggregate.Version) error {
	repo.mu.RLock()
	defer repo.mu.RUnlock()

	dto, ok := repo.indexedEvents[agg.AggregateID()]
	if !ok {
		return aggregate.ErrNotFound
	}

	if dto.version < int(version) {
		return aggregate.ErrNotFound
	}

	return aggregate.RestoreStateFromHistory(agg, filterEventsFromVersion(version, dto.events))
}

func (repo *AggregateRepository) Save(_ context.Context, agg aggregate.Aggregate[string]) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	dto, ok := repo.indexedEvents[agg.AggregateID()]
	if !ok {
		dto = &aggregateDTO{
			id:      agg.AggregateID(),
			name:    agg.AggregateName(),
			version: 0,
			events:  make([]aggregate.Event, 0),
		}
	}

	// copy the aggregate events
	events := make([]aggregate.Event, len(agg.AggregateEvents()))
	copy(events, agg.AggregateEvents())
	if c, ok := agg.(aggregate.EventCommitter); ok {
		c.CommitEvents()
	}

	// update the aggregate version
	dto.version = int(agg.AggregateVersion())

	// append the events to the aggregate
	dto.events = append(dto.events, events...)
	repo.indexedEvents[agg.AggregateID()] = dto

	// append the events to the repository
	repo.events = append(repo.events, events...)

	return nil
}

func (repo *AggregateRepository) Search(_ context.Context, criteria *aggregate.SearchCriteriaOptions) ([]aggregate.Aggregate[string], error) {
	repo.mu.RLock()
	defer repo.mu.RUnlock()

	events := repo.events
	if len(events) == 0 {
		return nil, nil
	}

	if criteria == nil || criteria.IsEmpty() {
		return getAggregates(events)
	}

	events = filterEventsByAggregateIDs(criteria.AggregateIDs(), events)
	events = filterEventsByAggregateNames(criteria.AggregateNames(), events)
	events = filterEventsByAggregateVersions(criteria.AggregateVersions(), events)
	return getAggregates(events)
}

func (repo *AggregateRepository) Delete(_ context.Context, agg aggregate.Aggregate[string]) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	repo.deleteAggregate(agg.AggregateID())
	return nil
}

func (repo *AggregateRepository) deleteAggregate(aggID string) {
	for idx, event := range repo.events {
		aggregateRef := event.Aggregate()
		if aggregateRef == nil {
			continue
		}

		if aggregateRef.ID == aggID {
			repo.events = append(repo.events[:idx], repo.events[idx+1:]...)
		}
	}

	delete(repo.indexedEvents, aggID)
}

func getAggregates(events []aggregate.Event) ([]aggregate.Aggregate[string], error) {
	var (
		aggregates     = make([]aggregate.Aggregate[string], 0)
		aggregateIndex = make(map[string]aggregate.Aggregate[string])
	)

	for _, event := range events {
		aggregateRef := event.Aggregate()
		if aggregateRef == nil {
			continue
		}

		aggregateID, ok := aggregateRef.ID.(string)
		if !ok {
			return nil, ErrInvalidAggregateEventID
		}

		agg, ok := aggregateIndex[aggregateID]
		if !ok {
			agg = aggregate.New(aggregateID, aggregateRef.Name)
		}

		err := aggregate.RestoreStateFromHistory(agg, []aggregate.Event{event})
		if err != nil {
			return nil, err
		}

		aggregateIndex[aggregateID] = agg
	}

	for _, agg := range aggregateIndex {
		aggregates = append(aggregates, agg)
	}

	return aggregates, nil
}
