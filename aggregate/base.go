package aggregate

var (
	_ ChangeCommitter = (*Base)(nil)
)

// Base provides the basic functionality of an aggregate.
// It implements the Aggregate and ChangeCommitter interfaces.
// It can be embedded in a custom aggregate type to provide the core
// functionality of an aggregate.
type Base struct {
	id      ID
	name    string
	version Version

	changes  []Change
	handlers map[string][]func(Change) error
}

// AggregateID returns the aggregate's ID.
func (a *Base) AggregateID() ID {
	return a.id
}

// AggregateName returns the aggregate's name.
func (a *Base) AggregateName() string {
	return a.name
}

// AggregateChanges returns the aggregate changes (events).
func (a *Base) AggregateChanges() []Change {
	return a.changes
}

// AggregateVersion returns the current version of the aggregate.
func (a *Base) AggregateVersion() Version {
	return a.version
}

// RecordChange adds the given changes as uncommitted events to the aggregate.
// It implements the ChangeCommitter interface.
func (a *Base) RecordChange(changes ...Change) {
	a.changes = append(a.changes, changes...)
}

// CommitChanges commits the aggregate's changes incrementing the version to the last change's version
// and resetting the changes list.
// It implements the ChangeCommitter interface.
func (b *Base) CommitChanges() {
	if len(b.changes) == 0 {
		return
	}

	b.version = Version(b.changes[len(b.changes)-1].Aggregate().Version)
	b.changes = b.changes[:0]
}

// RollbackChanges resets the aggregate's changes list.
// It implements the ChangeCommitter interface.
func (b *Base) RollbackChanges() {
	b.changes = b.changes[:0]
}

// When registers a handler for the given change (event) reason.
// The handler is called when the event is applied to the aggregate.
func (a *Base) When(reason string, handler func(change Change) error) {
	if a.handlers == nil {
		a.handlers = make(map[string][]func(change Change) error)
	}

	if _, ok := a.handlers[reason]; !ok {
		a.handlers[reason] = []func(change Change) error{}
	}

	a.handlers[reason] = append(a.handlers[reason], handler)
}

// ApplyChange calls the handlers for the given change (event) reason.
func (a *Base) ApplyChange(ev Change) {
	if a.handlers == nil {
		a.handlers = make(map[string][]func(Change) error)
	}

	if handlers, ok := a.handlers[ev.Reason()]; ok {
		for _, handler := range handlers {
			handler(ev)
		}
	}
}

// New creates a new base aggregate with the given ID and name.
// The default version is 0.
func New(id, name string) *Base {
	return &Base{
		id:      ID(id),
		name:    name,
		version: 0,
	}
}
