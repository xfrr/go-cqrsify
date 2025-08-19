package aggregate

// Version represents the version of an aggregate.
type Version int

// UncommittedVersion returns the aggregate version
// based on the uncommitted events.
// If there are no uncommitted events, it returns the current version.
func UncommittedVersion[ID comparable](agg EventSourcedAggregate[ID]) Version {
	events := agg.AggregateEvents()
	if len(events) == 0 {
		return agg.AggregateVersion()
	}

	return agg.AggregateVersion() + Version(len(events))
}
