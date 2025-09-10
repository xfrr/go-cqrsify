package domain

// AggregateVersion represents the version of an aggregate.
type AggregateVersion int

// UncommittedAggregateVersion returns the aggregate version
// based on the uncommitted events.
// If there are no uncommitted events, it returns the current version.
func UncommittedAggregateVersion[ID comparable](agg EventSourcedAggregate[ID]) AggregateVersion {
	events := agg.AggregateEvents()
	if len(events) == 0 {
		return agg.AggregateVersion()
	}

	return agg.AggregateVersion() + AggregateVersion(len(events))
}
