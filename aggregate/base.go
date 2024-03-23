package aggregate

var (
	_ ChangeCommitter = (*Base[any])(nil)
)

// Base provides the basic functionality of an aggregate.
// It implements the Aggregate and ChangeCommitter interfaces.
// It can be embedded in a custom aggregate type to provide the core
// functionality of an aggregate.
type Base[ID comparable] struct {
	id      ID
	name    string
	version Version

	changes  []Change
	handlers map[string][]func(Change)
}

// AggregateID returns the aggregate's ID.
func (a *Base[ID]) AggregateID() ID {
	return a.id
}

// AggregateName returns the aggregate's name.
func (a *Base[ID]) AggregateName() string {
	return a.name
}

// AggregateChanges returns the aggregate changes (events).
func (a *Base[ID]) AggregateChanges() []Change {
	return a.changes
}

// AggregateVersion returns the current version of the aggregate.
func (a *Base[ID]) AggregateVersion() Version {
	return a.version
}

// RecordChange adds the given changes as uncommitted events to the aggregate.
// It implements the ChangeCommitter interface.
func (a *Base[ID]) RecordChange(changes ...Change) {
	a.changes = append(a.changes, changes...)
}

// CommitChanges commits the aggregate's changes incrementing the version to the last change's version
// and resetting the changes list.
// It implements the ChangeCommitter interface.
func (b *Base[ID]) CommitChanges() {
	if len(b.changes) == 0 {
		return
	}

	b.version = Version(UncommittedVersion(b.Any()))
	b.changes = b.changes[:0]
}

// ClearChanges resets the aggregate's changes list.
// It implements the ChangeCommitter interface.
func (b *Base[ID]) ClearChanges() {
	b.changes = b.changes[:0]
}

// When registers a handler for the given change (event) reason.
// The handler is called when the event is applied to the aggregate.
func (a *Base[ID]) When(reason string, handler func(change Change)) {
	if a.handlers == nil {
		a.handlers = make(map[string][]func(change Change))
	}

	if _, ok := a.handlers[reason]; !ok {
		a.handlers[reason] = []func(change Change){}
	}

	a.handlers[reason] = append(a.handlers[reason], handler)
}

// ApplyChange calls the handlers for the given change (event) reason.
func (a *Base[ID]) ApplyChange(ev Change) {
	if a.handlers == nil {
		a.handlers = make(map[string][]func(Change))
	}

	if handlers, ok := a.handlers[ev.Reason()]; ok {
		for _, handler := range handlers {
			handler(ev)
		}
	}
}

// Any returns a copy of the aggregate with an arbitrary ID type.
func (a *Base[ID]) Any() *Base[any] {
	return &Base[any]{
		id:       a.id,
		name:     a.name,
		version:  a.version,
		changes:  a.changes,
		handlers: a.handlers,
	}
}

// New creates a new base aggregate with the given ID and name.
// The default version is 0.
func New[ID comparable](id ID, name string) *Base[ID] {
	return &Base[ID]{
		id:      id,
		name:    name,
		version: 0,
	}
}
