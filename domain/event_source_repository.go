package domain

import (
	"context"
	"fmt"
)

var _ EventSourcedRepository[any] = (*EventSourceRepository[any])(nil)

// EventSourceRepository represents a repository that provides access to an EventStore.
type EventSourceRepository[ID comparable] struct {
	eventStore EventStore[ID]
}

// NewEventSourceRepository creates a new EventSourceRepository with the given EventStore.
func NewEventSourceRepository[ID comparable](eventStore EventStore[ID]) *EventSourceRepository[ID] {
	return &EventSourceRepository[ID]{
		eventStore: eventStore,
	}
}

func (e *EventSourceRepository[ID]) Exists(ctx context.Context, agg EventSourcedAggregate[ID]) (bool, error) {
	events, err := e.eventStore.RetrieveMany(
		ctx,
		agg.AggregateID(),
		RetrieveEventsFromVersion(1),
		RetrieveEventsToVersion(1),
	)
	if err != nil {
		return false, fmt.Errorf("could not retrieve events: %w", err)
	}
	return len(events) > 0, nil
}

func (e *EventSourceRepository[ID]) ExistsVersion(
	ctx context.Context,
	agg EventSourcedAggregate[ID],
	version AggregateVersion,
) (bool, error) {
	events, err := e.eventStore.RetrieveMany(
		ctx,
		agg.AggregateID(),
		RetrieveEventsFromVersion(int(version)),
		RetrieveEventsToVersion(int(version)),
	)
	if err != nil {
		return false, fmt.Errorf("could not retrieve events: %w", err)
	}
	return len(events) > 0, nil
}

func (e *EventSourceRepository[ID]) Load(ctx context.Context, agg EventSourcedAggregate[ID]) error {
	events, err := e.eventStore.RetrieveMany(ctx, agg.AggregateID())
	if err != nil {
		return fmt.Errorf("could not retrieve events: %w", err)
	}

	if len(events) == 0 {
		return NewNotFoundError(agg.AggregateID())
	}

	return RestoreAggregateFromHistory(agg, events)
}

func (e *EventSourceRepository[ID]) LoadVersion(
	ctx context.Context,
	agg EventSourcedAggregate[ID],
	version AggregateVersion,
) error {
	events, err := e.eventStore.RetrieveMany(
		ctx,
		agg.AggregateID(),
		RetrieveEventsToVersion(int(version)),
	)
	if err != nil {
		return fmt.Errorf("could not retrieve events: %w", err)
	}

	if len(events) == 0 {
		return NewNotFoundError(agg.AggregateID())
	}

	return RestoreAggregateFromHistory(agg, events)
}

func (e *EventSourceRepository[ID]) Save(ctx context.Context, agg EventSourcedAggregate[ID]) error {
	err := e.eventStore.Save(ctx, agg.AggregateEvents())
	if err != nil {
		return fmt.Errorf("could not save aggregate events: %w", err)
	}

	if c, ok := agg.(EventCommitter); ok {
		c.CommitEvents()
	}

	return nil
}

func (e *EventSourceRepository[ID]) Search(
	ctx context.Context,
	opts *SearchCriteriaOptions,
) ([]EventSourcedAggregate[ID], error) {
	events, err := e.eventStore.Search(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("could not search events: %w", err)
	}

	if len(events) == 0 {
		return nil, nil
	}

	aggEvents := make(map[ID][]Event)
	for _, event := range events {
		id, ok := event.AggregateRef().ID().(ID)
		if !ok {
			return nil, fmt.Errorf("invalid aggregate ID type: %T", event.AggregateRef().ID())
		}
		aggEvents[id] = append(aggEvents[id], event)
	}

	result := make([]EventSourcedAggregate[ID], 0, len(aggEvents))
	for id, events := range aggEvents {
		firstEvent := events[0]
		agg := NewAggregate(id, firstEvent.AggregateRef().Name())
		if err != nil {
			return nil, fmt.Errorf("could not create aggregate: %w", err)
		}

		err = RestoreAggregateFromHistory(agg, events)
		if err != nil {
			return nil, fmt.Errorf("could not restore aggregate from history: %w", err)
		}

		result = append(result, agg)
	}

	return result, nil
}
