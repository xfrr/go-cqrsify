package inmemory

import (
	"slices"

	"github.com/xfrr/go-cqrsify/domain"
)

func filterEventsFromVersion(version domain.AggregateVersion, events []domain.Event) []domain.Event {
	filtered := make([]domain.Event, 0)
	for _, ch := range events {
		aggregateRef := ch.AggregateRef()
		if aggregateRef.Version() <= version {
			filtered = append(filtered, ch)
		}
	}

	return filtered
}

func filterEventsByAggregateIDs(aggIDs []string, events []domain.Event) []domain.Event {
	if len(aggIDs) == 0 {
		return events
	}

	filtered := make([]domain.Event, 0)
	for _, ev := range events {
		aggregateRef := ev.AggregateRef()

		aggID, _ := aggregateRef.ID().(string)
		if slices.Contains(aggIDs, aggID) {
			filtered = append(filtered, ev)
		}
	}

	return filtered
}

func filterEventsByAggregateNames(aggNames []string, events []domain.Event) []domain.Event {
	if len(aggNames) == 0 {
		return events
	}

	filtered := make([]domain.Event, 0)
	for _, ev := range events {
		aggregateRef := ev.AggregateRef()
		if slices.Contains(aggNames, aggregateRef.Name()) {
			filtered = append(filtered, ev)
		}
	}

	return filtered
}

func filterEventsByAggregateVersions(aggVersions []int, events []domain.Event) []domain.Event {
	if len(aggVersions) == 0 {
		return events
	}

	filtered := make([]domain.Event, 0)
	for _, ev := range events {
		aggregateRef := ev.AggregateRef()
		if slices.Contains(aggVersions, int(aggregateRef.Version())) {
			filtered = append(filtered, ev)
		}
	}

	return filtered
}
