package aggregate

type HistoryIntegrityError struct {
	desc string
}

func (e *HistoryIntegrityError) Error() string {
	return e.desc
}

func NewHistoryIntegrityError(desc string) *HistoryIntegrityError {
	return &HistoryIntegrityError{desc: desc}
}

// VerifyHistoryIntegrity verifies the integrity of the given history against the given Aggregate.
func VerifyHistoryIntegrity[ID comparable](a Aggregate[ID], events []Event) error {
	for i, c := range events {
		if c.Aggregate() == nil {
			return NewHistoryIntegrityError("event has no aggregate")
		}

		if c.Aggregate().ID != a.AggregateID() {
			return NewHistoryIntegrityError("event has different aggregate ID")
		}

		if c.Aggregate().Name != a.AggregateName() {
			return NewHistoryIntegrityError("event has different aggregate name")
		}

		expectedVersion := UncommittedVersion(a) + i + 1
		if c.Aggregate().Version != expectedVersion {
			return NewHistoryIntegrityError("event has unexpected version")
		}
	}

	return nil
}
