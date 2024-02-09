package aggregate

// Hydrate applies the given changes (events) to the given Aggregate, ensuring
// consistency and updating the Aggregate's state accordingly.
// It records and commits the changes if the Aggregate implements the ChangeCommiter interface.
// It returns an error if the events cannot be applied.
func Hydrate[Changes ~[]Change](a Aggregate, changes Changes) error {
	for _, change := range changes {
		if err := ValidateChange(a, change); err != nil {
			if c, ok := a.(ChangeCommiter); len(a.AggregateChanges()) > 0 && ok {
				c.RollbackChanges()
			}
			return err
		}

		a.ApplyChange(change)
		if r, ok := a.(ChangeCommiter); ok {
			r.RecordChange(change)
		}
	}

	if c, ok := a.(ChangeCommiter); len(a.AggregateChanges()) > 0 && ok {
		c.CommitChanges()
	}

	return nil
}
