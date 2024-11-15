package aggregate

// Version represents the version of an aggregate.
type Version int

// UncommittedVersion returns the aggregate version
// based on the uncommitted events.
// If there are no uncommitted events, it returns the current version.
func UncommittedVersion[ID comparable](agg Aggregate[ID]) int {
	events := agg.AggregateEvents()
	if len(events) == 0 {
		return int(agg.AggregateVersion())
	}

	latestEvent := events[len(events)-1]
	if latestEvent.Aggregate() == nil {
		return int(agg.AggregateVersion())
	}

	return latestEvent.Aggregate().Version
}

// nextVersion returns the next version of the aggregate.
func nextVersion[ID comparable](agg Aggregate[ID]) int {
	return UncommittedVersion(agg) + 1
}
