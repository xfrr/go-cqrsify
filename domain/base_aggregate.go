package domain

import "github.com/xfrr/go-cqrsify/pkg/multierror"

var (
	_ EventCommitter = (*BaseAggregate[any])(nil)
)

// BaseAggregate implements the core functionality of an Aggregate.
// It must be embedded in a struct to implements the Aggregate interface.
type BaseAggregate[ID comparable] struct {
	id      ID
	name    string
	version AggregateVersion

	events   []Event
	handlers map[string][]func(Event) error
}

// AggregateID returns the aggregate's ID.
func (agb *BaseAggregate[ID]) AggregateID() ID {
	return agb.id
}

// AggregateName returns the aggregate's name.
func (agb *BaseAggregate[ID]) AggregateName() string {
	return agb.name
}

// AggregateEvents returns the aggregate uncommitted events.
func (agb *BaseAggregate[ID]) AggregateEvents() []Event {
	return agb.events
}

// AggregateVersion returns the current version of the aggregate.
func (agb *BaseAggregate[ID]) AggregateVersion() AggregateVersion {
	return agb.version
}

// RecordEvent adds the given events as uncommitted events to the aggregate.
// It implements the EventCommitter interface.
func (agb *BaseAggregate[ID]) RecordEvent(event Event) {
	agb.events = append(agb.events, event)
}

// CommitEvents commits the aggregate's events incrementing the version to the last event's version
// and resetting the events list.
// It implements the EventCommitter interface.
func (agb *BaseAggregate[ID]) CommitEvents() {
	if len(agb.events) == 0 {
		return
	}

	agb.version = AggregateVersion(UncommittedAggregateVersion(agb))
	agb.events = agb.events[:0]
}

// ClearEvents resets the aggregate's events list.
// It implements the EventCommitter interface.
func (agb *BaseAggregate[ID]) ClearEvents() {
	agb.events = agb.events[:0]
}

// HandleEvent registers a handler for the given event name.
// The handler is called when the event is applied to the aggregate.
func (agb *BaseAggregate[ID]) HandleEvent(name string, handler func(event Event) error) {
	if agb.handlers == nil {
		agb.handlers = make(map[string][]func(event Event) error)
	}

	if _, ok := agb.handlers[name]; !ok {
		agb.handlers[name] = []func(event Event) error{}
	}

	agb.handlers[name] = append(agb.handlers[name], handler)
}

// ApplyEvent calls the handlers for the given event (event) name.
func (agb *BaseAggregate[ID]) ApplyEvent(ev Event) error {
	if agb.handlers == nil {
		agb.handlers = make(map[string][]func(Event) error)
	}

	multiErr := multierror.New()
	if handlers, ok := agb.handlers[ev.Name()]; ok {
		for _, handler := range handlers {
			if err := handler(ev); err != nil {
				multiErr.Append(err)
			}
		}
	}

	return multiErr.ErrorOrNil()
}

// Any returns a copy of the aggregate with an arbitrary ID type.
func (agb *BaseAggregate[ID]) Any() *BaseAggregate[any] {
	return &BaseAggregate[any]{
		id:       agb.id,
		name:     agb.name,
		version:  agb.version,
		events:   agb.events,
		handlers: agb.handlers,
	}
}

// NewAggregate creates a new base aggregate with the given ID and name.
// The default version is 0.
func NewAggregate[ID comparable](id ID, name string) *BaseAggregate[ID] {
	return &BaseAggregate[ID]{
		id:       id,
		name:     name,
		version:  0,
		events:   []Event{},
		handlers: make(map[string][]func(Event) error),
	}
}
