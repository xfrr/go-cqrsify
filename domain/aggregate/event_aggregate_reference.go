package aggregate

type EventAggregateReference struct {
	aggregateID      any
	aggregateType    string
	aggregateVersion Version
}

func (r EventAggregateReference) ID() any {
	return r.aggregateID
}

func (r EventAggregateReference) Type() string {
	return r.aggregateType
}

func (r EventAggregateReference) Version() Version {
	return r.aggregateVersion
}

// CreateEventAggregateRef creates a new EventAggregateReference for the given EventSourcedAggregate.
// It increments the aggregate version.
func CreateEventAggregateRef[T comparable](agg EventSourcedAggregate[T]) *EventAggregateReference {
	return newEventAggregateReference(
		agg.AggregateID(),
		agg.AggregateName(),
		UncommittedVersion(agg)+1,
	)
}

func newEventAggregateReference(id any, aggregateType string, version Version) *EventAggregateReference {
	return &EventAggregateReference{
		aggregateID:      id,
		aggregateType:    aggregateType,
		aggregateVersion: version,
	}
}
