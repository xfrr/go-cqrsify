package aggregate

import "errors"

var (
	// ErrInvalidBaseAggregate is returned when the given Aggregate is not a BaseAggregate.
	ErrInvalidBaseAggregate = errors.New("invalid base aggregate")
)

// Hydrate applies the given changes (events) to the given Aggregate, ensuring
// consistency and updating the Aggregate's state accordingly.
// It records and commits the changes if the Aggregate implements the ChangeCommitter interface.
// It returns an error if the events cannot be applied.
func Hydrate[ID comparable](a Aggregate[ID], changes []Change) error {
	base, ok := Cast[ID](a)
	if !ok {
		return ErrInvalidBaseAggregate
	}

	if err := ValidateChanges(base, changes); err != nil {
		return err
	}

	for _, change := range changes {
		a.ApplyChange(change)
	}

	if c, ok := a.(ChangeCommitter); ok {
		c.RecordChange(changes...)
		c.CommitChanges()
	}

	return nil
}
