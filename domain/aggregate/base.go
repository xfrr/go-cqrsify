package aggregate

import "github.com/xfrr/go-cqrsify/pkg/multierror"

var (
	_ EventCommitter = (*Base[any])(nil)
)

// Base implements the core functionality of an Aggregate.
// It must be embedded in a struct to implements the Aggregate interface.
type Base[ID comparable] struct {
	id      ID
	name    string
	version Version

	events   []Event
	handlers map[string][]func(Event) error
}

// AggregateID returns the aggregate's ID.
func (agb *Base[ID]) AggregateID() ID {
	return agb.id
}

// AggregateName returns the aggregate's name.
func (agb *Base[ID]) AggregateName() string {
	return agb.name
}

// AggregateEvents returns the aggregate uncommitted events.
func (agb *Base[ID]) AggregateEvents() []Event {
	return agb.events
}

// AggregateVersion returns the current version of the aggregate.
func (agb *Base[ID]) AggregateVersion() Version {
	return agb.version
}

// RecordEvent adds the given events as uncommitted events to the aggregate.
// It implements the EventCommitter interface.
func (agb *Base[ID]) RecordEvent(event Event) {
	agb.events = append(agb.events, event)
}

// CommitEvents commits the aggregate's events incrementing the version to the last event's version
// and resetting the events list.
// It implements the EventCommitter interface.
func (agb *Base[ID]) CommitEvents() {
	if len(agb.events) == 0 {
		return
	}

	agb.version = Version(UncommittedVersion(agb))
	agb.events = agb.events[:0]
}

// ClearEvents resets the aggregate's events list.
// It implements the EventCommitter interface.
func (agb *Base[ID]) ClearEvents() {
	agb.events = agb.events[:0]
}

// HandleEvent registers a handler for the given event name.
// The handler is called when the event is applied to the aggregate.
func (agb *Base[ID]) HandleEvent(name string, handler func(event Event) error) {
	if agb.handlers == nil {
		agb.handlers = make(map[string][]func(event Event) error)
	}

	if _, ok := agb.handlers[name]; !ok {
		agb.handlers[name] = []func(event Event) error{}
	}

	agb.handlers[name] = append(agb.handlers[name], handler)
}

// ApplyEvent calls the handlers for the given event (event) name.
func (agb *Base[ID]) ApplyEvent(ev Event) error {
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
func (agb *Base[ID]) Any() *Base[any] {
	return &Base[any]{
		id:       agb.id,
		name:     agb.name,
		version:  agb.version,
		events:   agb.events,
		handlers: agb.handlers,
	}
}

// New creates a new base aggregate with the given ID and name.
// The default version is 0.
func New[ID comparable](id ID, name string) *Base[ID] {
	return &Base[ID]{
		id:       id,
		name:     name,
		version:  0,
		events:   []Event{},
		handlers: make(map[string][]func(Event) error),
	}
}
