package domain

// EventCommitter provides methods for tracking aggregate events
// as a series of events and committing them to update the aggregate's state.
type EventCommitter interface {
	// CommitEvents increments the aggregate version and clears the uncommitted events.
	CommitEvents()

	// ClearEvents clears the uncommitted events.
	ClearEvents()
}
