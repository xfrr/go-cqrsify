package memory

import (
	"slices"

	"github.com/xfrr/cqrsify/aggregate"
)

func filterChangesFromVersion(version aggregate.Version, changes []aggregate.Change) []aggregate.Change {
	filtered := make([]aggregate.Change, 0)
	for _, ch := range changes {
		aggregateRef := ch.Aggregate()
		if aggregateRef == nil {
			continue
		}

		if int(aggregateRef.Version) <= int(version) {
			filtered = append(filtered, ch)
		}
	}

	return filtered
}

func filterChangesByAggregateIDs(aggIDs []string, changes []aggregate.Change) []aggregate.Change {
	if len(aggIDs) == 0 {
		return changes
	}

	filtered := make([]aggregate.Change, 0)
	for _, ev := range changes {
		aggregateRef := ev.Aggregate()
		if aggregateRef == nil {
			continue
		}

		aggID, _ := aggregateRef.ID.(string)
		if slices.Contains(aggIDs, aggID) {
			filtered = append(filtered, ev)
		}
	}

	return filtered
}

func filterChangesByAggregateNames(aggNames []string, changes []aggregate.Change) []aggregate.Change {
	if len(aggNames) == 0 {
		return changes
	}

	filtered := make([]aggregate.Change, 0)
	for _, ev := range changes {
		aggregateRef := ev.Aggregate()
		if aggregateRef == nil {
			continue
		}

		if slices.Contains(aggNames, aggregateRef.Name) {
			filtered = append(filtered, ev)
		}
	}

	return filtered
}

func filterChangesByAggregateVersions(aggVersions []int, changes []aggregate.Change) []aggregate.Change {
	if len(aggVersions) == 0 {
		return changes
	}

	filtered := make([]aggregate.Change, 0)
	for _, ev := range changes {
		aggregateRef := ev.Aggregate()
		if aggregateRef == nil {
			continue
		}

		if slices.Contains(aggVersions, int(aggregateRef.Version)) {
			filtered = append(filtered, ev)
		}
	}

	return filtered
}
