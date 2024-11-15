package aggregate

// History represents a sequence of events
// that have occurred in an Aggregate.
type History []Event

// RestoreStateFromHistory restores the state of the Aggregate from the given History of events.
//
// It records and commits the events if the Aggregate implements the EventCommitter interface.
//
// It returns an error if the events cannot be applied.
func RestoreStateFromHistory[ID comparable](a Aggregate[ID], events History) error {
	if err := VerifyHistoryIntegrity(a, events); err != nil {
		return err
	}

	for _, event := range events {
		a.ApplyEvent(event)
	}

	if r, ok := a.(EventRecorder); ok {
		for _, event := range events {
			r.RecordEvent(event)
		}
	}

	if c, ok := a.(EventCommitter); ok {
		c.CommitEvents()
	}

	return nil
}
