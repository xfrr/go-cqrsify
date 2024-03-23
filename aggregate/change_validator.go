package aggregate

import "errors"

var (
	// ErrInvalidAggregateID is returned when the change's aggregate ID does not match the aggregate's ID.
	ErrInvalidAggregateID = errors.New("invalid aggregate ID")

	// ErrInvalidAggregateName is returned when the change's aggregate name does not match the aggregate's name.
	ErrInvalidAggregateName = errors.New("invalid aggregate name")

	// ErrInvalidVersion is returned when the change's version does not match the aggregate's version.
	ErrInvalidVersion = errors.New("invalid aggregate version")

	// ErrInvalidChangePayload is returned when the change's payload is nil.
	ErrInvalidChangePayload = errors.New("invalid change payload")

	// ErrInvalidChangeTime is returned when the change's time is zero.
	ErrInvalidChangeTime = errors.New("invalid change time")

	// ErrInvalidEventAggregateReference is returned when the change's aggregate is nil.
	ErrInvalidEventAggregateReference = errors.New("invalid event aggregate, event must carry an aggregate reference")
)

// ValidateChanges validates the given change (event) against the Aggregate's state.
// It returns an error if the change is not valid.
func ValidateChanges[ID comparable](a Aggregate[ID], changes []Change) error {
	for i, c := range changes {
		if c.Aggregate() == nil {
			return ErrInvalidEventAggregateReference
		}

		if c.Aggregate().ID != a.AggregateID() {
			return ErrInvalidAggregateID
		}

		if c.Aggregate().Name != a.AggregateName() {
			return ErrInvalidAggregateName
		}

		if c.Payload() == nil {
			return ErrInvalidChangePayload
		}

		expectedVersion := UncommittedVersion(a) + i + 1
		if c.Aggregate().Version != expectedVersion {
			return ErrInvalidVersion
		}

		if c.Time().IsZero() {
			return ErrInvalidChangeTime
		}
	}

	return nil
}
