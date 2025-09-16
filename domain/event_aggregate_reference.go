package domain

type EventAggregateReference struct {
	aggregateID      any
	aggregateName    string
	aggregateVersion AggregateVersion
}

func (r EventAggregateReference) ID() any {
	return r.aggregateID
}

func (r EventAggregateReference) Name() string {
	return r.aggregateName
}

func (r EventAggregateReference) Version() AggregateVersion {
	return r.aggregateVersion
}

// CreateEventAggregateRef creates a new EventAggregateReference for the given EventSourcedAggregate.
// It increments the aggregate version.
func CreateEventAggregateRef[T comparable](agg EventSourcedAggregate[T]) *EventAggregateReference {
	return newEventAggregateReference(
		agg.AggregateID(),
		agg.AggregateName(),
		UncommittedAggregateVersion(agg)+1,
	)
}

func newEventAggregateReference(
	aggregateID any,
	aggregateName string,
	version AggregateVersion,
) *EventAggregateReference {
	return &EventAggregateReference{
		aggregateID:      aggregateID,
		aggregateName:    aggregateName,
		aggregateVersion: version,
	}
}
