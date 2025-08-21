package aggregate

import (
	"errors"
	"fmt"
)

type HistoryIntegrityError struct {
	desc          string
	EventIndex    int
	ExpectedValue interface{}
	ActualValue   interface{}
	ErrorType     string // "ID_MISMATCH", "TYPE_MISMATCH", "VERSION_MISMATCH"
}

func (e *HistoryIntegrityError) WithDetails(eventIndex int, expected, actual interface{}, errorType string) *HistoryIntegrityError {
	e.EventIndex = eventIndex
	e.ExpectedValue = expected
	e.ActualValue = actual
	e.ErrorType = errorType
	return e
}

func (e *HistoryIntegrityError) Error() string {
	return "history integrity error. " + e.desc + fmt.Sprintf(" (event index: %d, expected: %v, actual: %v, type: %s)",
		e.EventIndex, e.ExpectedValue, e.ActualValue, e.ErrorType)
}

func NewHistoryIntegrityError(desc string) *HistoryIntegrityError {
	return &HistoryIntegrityError{desc: desc}
}

var (
	ErrEmptyEventHistory     = errors.New("event history cannot be empty")
	ErrNilAggregate          = errors.New("aggregate cannot be nil")
	ErrAggregateIDMismatch   = NewHistoryIntegrityError("event has different aggregate ID")
	ErrAggregateTypeMismatch = NewHistoryIntegrityError("event has different aggregate type")
)

// VerifyHistoryIntegrity verifies the integrity of the given event history against the given Aggregate.
// It performs the following checks:
//   - All events belong to the same aggregate (ID and type)
//   - Event versions are sequential and start from the aggregate's uncommitted version + 1
//   - No events are nil
//
// Returns a HistoryIntegrityError if any integrity check fails, nil otherwise.
func VerifyHistoryIntegrity[ID comparable](agg EventSourcedAggregate[ID], events []Event) error {
	// Input validation
	if agg == nil {
		return ErrNilAggregate
	}

	if len(events) == 0 {
		return ErrEmptyEventHistory
	}

	aggregateID := agg.AggregateID()
	aggregateName := agg.AggregateName()
	baseVersion := UncommittedVersion(agg)

	for i, evt := range events {
		// Check for nil events
		if evt == nil {
			return NewHistoryIntegrityError(fmt.Sprintf("event at index %d is nil", i)).
				WithDetails(i, nil, evt, "NIL_EVENT")
		}

		eventRef := evt.AggregateRef()
		if eventRef == nil {
			return NewHistoryIntegrityError(fmt.Sprintf("event at index %d has nil aggregate reference", i)).
				WithDetails(i, nil, eventRef, "NIL_EVENT_REF")
		}

		// Verify aggregate ID matches
		if eventRef.ID() != aggregateID {
			return NewHistoryIntegrityError(fmt.Sprintf("event at index %d has different aggregate ID: got %v, want %v",
				i, eventRef.ID(), aggregateID)).
				WithDetails(i, aggregateID, eventRef.ID(), "ID_MISMATCH")
		}

		// Verify aggregate type matches
		if eventRef.Type() != aggregateName {
			return NewHistoryIntegrityError(fmt.Sprintf("event at index %d has different aggregate type: got %q, want %q",
				i, eventRef.Type(), aggregateName)).
				WithDetails(i, aggregateName, eventRef.Type(), "TYPE_MISMATCH")
		}

		// Verify version is sequential
		expectedVersion := baseVersion + Version(i) + 1
		actualVersion := eventRef.Version()
		if actualVersion != expectedVersion {
			return NewHistoryIntegrityError(fmt.Sprintf("event at index %d has unexpected version: got %d, want %d",
				i, actualVersion, expectedVersion)).
				WithDetails(i, expectedVersion, actualVersion, "VERSION_MISMATCH")
		}
	}

	return nil
}
