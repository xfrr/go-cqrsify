package aggregate

// ChangeCommitter provides methods for tracking aggregate changes
// as a series of events and committing them to update the aggregate's state.
type ChangeCommitter interface {
	// RecordChange adds the given changes as events to the Aggregate history.
	// It should not update the Aggregate's state.
	RecordChange(...Change)

	// CommitChanges increments the aggregate version and clears the uncommitted changes.
	CommitChanges()

	// ClearChanges clears the uncommitted changes.
	ClearChanges()
}
