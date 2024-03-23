package aggregate

// Version represents the version of an aggregate.
type Version int

// UncommittedVersion returns the aggregate version
// based on the uncommitted events.
// If there are no uncommitted events, it returns the current version.
func UncommittedVersion[ID comparable](agg Aggregate[ID]) int {
	if len(agg.AggregateChanges()) == 0 {
		return int(agg.AggregateVersion())
	}

	changes := agg.AggregateChanges()
	if len(changes) == 0 {
		return int(agg.AggregateVersion())
	}

	change := changes[len(changes)-1]
	if change.Aggregate() == nil {
		return int(agg.AggregateVersion())
	}

	return change.Aggregate().Version
}

// nextVersion returns the next version of the aggregate.
func nextVersion[ID comparable](agg Aggregate[ID]) int {
	return UncommittedVersion(agg) + 1
}
