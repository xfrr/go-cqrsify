package inmemory

import (
	"slices"

	"github.com/xfrr/go-cqrsify/domain/aggregate"
)

func filterEventsFromVersion(version aggregate.Version, events []aggregate.Event) []aggregate.Event {
	filtered := make([]aggregate.Event, 0)
	for _, ch := range events {
		aggregateRef := ch.AggregateRef()
		if aggregateRef.Version() <= version {
			filtered = append(filtered, ch)
		}
	}

	return filtered
}

func filterEventsByAggregateIDs(aggIDs []string, events []aggregate.Event) []aggregate.Event {
	if len(aggIDs) == 0 {
		return events
	}

	filtered := make([]aggregate.Event, 0)
	for _, ev := range events {
		aggregateRef := ev.AggregateRef()

		aggID, _ := aggregateRef.ID().(string)
		if slices.Contains(aggIDs, aggID) {
			filtered = append(filtered, ev)
		}
	}

	return filtered
}

func filterEventsByAggregateNames(aggNames []string, events []aggregate.Event) []aggregate.Event {
	if len(aggNames) == 0 {
		return events
	}

	filtered := make([]aggregate.Event, 0)
	for _, ev := range events {
		aggregateRef := ev.AggregateRef()
		if slices.Contains(aggNames, aggregateRef.Type()) {
			filtered = append(filtered, ev)
		}
	}

	return filtered
}

func filterEventsByAggregateVersions(aggVersions []int, events []aggregate.Event) []aggregate.Event {
	if len(aggVersions) == 0 {
		return events
	}

	filtered := make([]aggregate.Event, 0)
	for _, ev := range events {
		aggregateRef := ev.AggregateRef()
		if slices.Contains(aggVersions, int(aggregateRef.Version())) {
			filtered = append(filtered, ev)
		}
	}

	return filtered
}
