package inmemory

import (
	"context"
	"errors"
	"sync"

	"github.com/xfrr/go-cqrsify/domain/aggregate"
)

var _ aggregate.EventSourcedRepository[string] = (*EventSourcedAggregateRepository)(nil)

var (
	ErrInvalidAggregateEventID = errors.New("invalid aggregate event id")
)

type eventSourcedAggregateDTO struct {
	id      string
	name    string
	version int
	events  []aggregate.Event
}

type EventSourcedAggregateRepository struct {
	mu sync.RWMutex

	dtosIndex map[string]*eventSourcedAggregateDTO
	events    []aggregate.Event
}

func NewEventSourcedAggregateRepository() *EventSourcedAggregateRepository {
	return &EventSourcedAggregateRepository{
		mu:        sync.RWMutex{},
		dtosIndex: make(map[string]*eventSourcedAggregateDTO),
		events:    make([]aggregate.Event, 0),
	}
}

func (repo *EventSourcedAggregateRepository) Exists(_ context.Context, agg aggregate.EventSourcedAggregate[string]) (bool, error) {
	repo.mu.RLock()
	defer repo.mu.RUnlock()

	_, ok := repo.dtosIndex[agg.AggregateID()]
	return ok, nil
}

func (repo *EventSourcedAggregateRepository) ExistsVersion(_ context.Context, agg aggregate.EventSourcedAggregate[string], version aggregate.Version) (bool, error) {
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

func (repo *EventSourcedAggregateRepository) Load(_ context.Context, agg aggregate.EventSourcedAggregate[string]) error {
	repo.mu.RLock()
	defer repo.mu.RUnlock()

	dto, ok := repo.dtosIndex[agg.AggregateID()]
	if !ok {
		return aggregate.NewNotFoundError(agg.AggregateID())
	}

	return aggregate.RestoreFromHistory(agg, dto.events)
}

func (repo *EventSourcedAggregateRepository) LoadVersion(_ context.Context, agg aggregate.EventSourcedAggregate[string], version aggregate.Version) error {
	repo.mu.RLock()
	defer repo.mu.RUnlock()

	dto, ok := repo.dtosIndex[agg.AggregateID()]
	if !ok {
		return aggregate.NewNotFoundError(agg.AggregateID())
	}

	if dto.version < int(version) {
		return aggregate.NewNotFoundError(agg.AggregateID())
	}

	return aggregate.RestoreFromHistory(agg, filterEventsFromVersion(version, dto.events))
}

func (repo *EventSourcedAggregateRepository) Save(ctx context.Context, agg aggregate.EventSourcedAggregate[string]) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	return repo.saveEvents(agg)
}

func (repo *EventSourcedAggregateRepository) saveEvents(agg aggregate.EventSourcedAggregate[string]) error {
	dto, ok := repo.dtosIndex[agg.AggregateID()]
	if !ok {
		dto = &eventSourcedAggregateDTO{
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
	repo.dtosIndex[agg.AggregateID()] = dto

	// append the events to the repository
	repo.events = append(repo.events, events...)

	return nil
}

func (repo *EventSourcedAggregateRepository) Search(_ context.Context, criteria *aggregate.SearchCriteriaOptions) ([]aggregate.EventSourcedAggregate[string], error) {
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

func (repo *EventSourcedAggregateRepository) Delete(_ context.Context, agg aggregate.EventSourcedAggregate[string]) error {
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

func getAggregates(events []aggregate.Event) ([]aggregate.EventSourcedAggregate[string], error) {
	var (
		aggregates     = make([]aggregate.EventSourcedAggregate[string], 0)
		aggregateIndex = make(map[string]aggregate.EventSourcedAggregate[string])
	)

	for _, event := range events {
		aggregateRef := event.AggregateRef()
		aggregateID, ok := aggregateRef.ID().(string)
		if !ok {
			return nil, ErrInvalidAggregateEventID
		}

		agg, ok := aggregateIndex[aggregateID]
		if !ok {
			agg = aggregate.New(aggregateID, aggregateRef.Type())
		}

		err := aggregate.RestoreFromHistory(agg, []aggregate.Event{event})
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
