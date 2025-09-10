package domain

// History represents a sequence of events
// that have occurred in an Aggregate.
type History []Event

// RestoreAggregateFromHistory restores the state of the Aggregate from the given History of events.
//
// It records and commits the events if the Aggregate implements the EventCommitter interface.
//
// It returns an error if the events cannot be applied.
func RestoreAggregateFromHistory[ID comparable](agg EventSourcedAggregate[ID], events History) error {
	if err := VerifyHistoryIntegrity(agg, events); err != nil {
		return err
	}

	for _, event := range events {
		agg.ApplyEvent(event)
	}

	if r, ok := agg.(EventRecorder); ok {
		for _, event := range events {
			r.RecordEvent(event)
		}
	}

	if c, ok := agg.(EventCommitter); ok {
		c.CommitEvents()
	}

	return nil
}
