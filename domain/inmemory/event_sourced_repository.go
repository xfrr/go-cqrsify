package inmemory

import (
	"context"
	"errors"
	"sync"

	"github.com/xfrr/go-cqrsify/domain"
)

var _ domain.EventSourcedRepository[domain.EventSourcedAggregate[string], string] = (*EventSourcedAggregateRepository)(nil)

var (
	ErrInvalidAggregateEventID = errors.New("invalid aggregate event id")
)

type eventSourcedAggregateDTO struct {
	id      string
	name    string
	version int
	events  []domain.Event
}

type EventSourcedAggregateRepository struct {
	mu sync.RWMutex

	dtosIndex map[string]*eventSourcedAggregateDTO
	events    []domain.Event
}

func NewEventSourcedAggregateRepository() *EventSourcedAggregateRepository {
	return &EventSourcedAggregateRepository{
		mu:        sync.RWMutex{},
		dtosIndex: make(map[string]*eventSourcedAggregateDTO),
		events:    make([]domain.Event, 0),
	}
}

func (repo *EventSourcedAggregateRepository) Exists(_ context.Context, agg domain.EventSourcedAggregate[string]) (bool, error) {
	repo.mu.RLock()
	defer repo.mu.RUnlock()

	_, ok := repo.dtosIndex[agg.AggregateID()]
	return ok, nil
}

func (repo *EventSourcedAggregateRepository) ExistsVersion(_ context.Context, agg domain.EventSourcedAggregate[string], version domain.AggregateVersion) (bool, error) {
	repo.mu.RLock()
	defer repo.mu.RUnlock()

	dto, ok := repo.dtosIndex[agg.AggregateID()]
	if !ok {
		return false, nil
	}

	if dto.version < int(version) {
		return false, nil
	}

	return true, nil
}

func (repo *EventSourcedAggregateRepository) Load(_ context.Context, agg domain.EventSourcedAggregate[string]) error {
	repo.mu.RLock()
	defer repo.mu.RUnlock()

	dto, ok := repo.dtosIndex[agg.AggregateID()]
	if !ok {
		return domain.NewNotFoundError(agg.AggregateID())
	}

	return domain.RestoreAggregateFromHistory(agg, dto.events)
}

func (repo *EventSourcedAggregateRepository) LoadVersion(_ context.Context, agg domain.EventSourcedAggregate[string], version domain.AggregateVersion) error {
	repo.mu.RLock()
	defer repo.mu.RUnlock()

	dto, ok := repo.dtosIndex[agg.AggregateID()]
	if !ok {
		return domain.NewNotFoundError(agg.AggregateID())
	}

	if dto.version < int(version) {
		return domain.NewNotFoundError(agg.AggregateID())
	}

	return domain.RestoreAggregateFromHistory(agg, filterEventsFromVersion(version, dto.events))
}

func (repo *EventSourcedAggregateRepository) Save(ctx context.Context, agg domain.EventSourcedAggregate[string]) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	return repo.saveEvents(agg)
}

func (repo *EventSourcedAggregateRepository) saveEvents(agg domain.EventSourcedAggregate[string]) error {
	dto, ok := repo.dtosIndex[agg.AggregateID()]
	if !ok {
		dto = &eventSourcedAggregateDTO{
			id:      agg.AggregateID(),
			name:    agg.AggregateName(),
			version: 0,
			events:  make([]domain.Event, 0),
		}
	}

	// copy the aggregate events
	events := make([]domain.Event, len(agg.AggregateEvents()))
	copy(events, agg.AggregateEvents())
	if c, ok := agg.(domain.EventCommitter); ok {
		c.CommitEvents()
	}

	// update the aggregate version
	dto.version = int(agg.AggregateVersion())

	// append the events to the aggregate
	dto.events = append(dto.events, events...)
	repo.dtosIndex[agg.AggregateID()] = dto

	// append the events to the repository
	repo.events = append(repo.events, events...)

	return nil
}

func (repo *EventSourcedAggregateRepository) Search(_ context.Context, criteria *domain.SearchCriteriaOptions) ([]domain.EventSourcedAggregate[string], error) {
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

func (repo *EventSourcedAggregateRepository) Delete(_ context.Context, agg domain.EventSourcedAggregate[string]) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	repo.deleteAggregate(agg.AggregateID())
	return nil
}

func (repo *EventSourcedAggregateRepository) deleteAggregate(aggID string) {
	for idx, event := range repo.events {
		aggregateRef := event.AggregateRef()
		if aggregateRef.ID() == aggID {
			repo.events = append(repo.events[:idx], repo.events[idx+1:]...)
		}
	}

	delete(repo.dtosIndex, aggID)
}

func getAggregates(events []domain.Event) ([]domain.EventSourcedAggregate[string], error) {
	var (
		aggregates     = make([]domain.EventSourcedAggregate[string], 0)
		aggregateIndex = make(map[string]domain.EventSourcedAggregate[string])
	)

	for _, event := range events {
		aggregateRef := event.AggregateRef()
		aggregateID, ok := aggregateRef.ID().(string)
		if !ok {
			return nil, ErrInvalidAggregateEventID
		}

		agg, ok := aggregateIndex[aggregateID]
		if !ok {
			agg = domain.NewAggregate(aggregateID, aggregateRef.Name())
		}

		err := domain.RestoreAggregateFromHistory(agg, []domain.Event{event})
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
