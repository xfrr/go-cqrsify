package command

// Option is an option for creating a command.
type Option func(*command[any])

// Aggregate returns an Option that links a command to an aggregate.
func WithAggregate(id string, name string) Option {
	return func(b *command[any]) {
		b.props.AggregateID = id
		b.props.AggregateName = name
	}
}

// ID represents a command identifier.
// Must be unique across all commands.
type ID string

// Command represents a command with an arbitrary message.
type Command[M any] interface {
	// ID returns the command id.
	ID() ID

	// Message returns the command message.
	Message() M

	// AggregateID returns the id of the aggregate that the command acts on.
	AggregateID() string

	// AggregateName returns the name of the aggregate that the command acts on.
	AggregateName() string
}

// command is the internal implementation of Command.
type command[Message any] struct {
	props properties[Message]
}

// properties contains the fields of a Cmd.
type properties[Message any] struct {
	ID            ID
	Message       Message
	AggregateName string
	AggregateID   string
}

// New returns a new command with the given name and message.
func New[M any](id string, msg M, opts ...Option) command[M] {
	c := command[any]{
		props: properties[any]{
			ID:      ID(id),
			Message: msg,
		},
	}
	for _, opt := range opts {
		opt(&c)
	}
	return command[M]{
		props: properties[M]{
			ID:            c.props.ID,
			Message:       c.props.Message.(M),
			AggregateName: c.props.AggregateName,
			AggregateID:   c.props.AggregateID,
		},
	}
}

// ID returns the command id.
func (c command[M]) ID() ID {
	return c.props.ID
}

// Message returns the command message.
func (c command[M]) Message() M {
	return c.props.Message
}

// AggregateID returns the id of the aggregate that the command acts on.
func (c command[M]) AggregateID() string {
	return c.props.AggregateID
}

// AggregateName returns the name of the aggregate that the command acts on.
func (c command[M]) AggregateName() string {
	return c.props.AggregateName
}

// Any returns a new command with the same id and aggregate as the given command
// but with an arbitrary message.
func (c command[M]) Any() Command[any] {
	return New[any](string(c.ID()), c.Message(), WithAggregate(c.AggregateID(), c.AggregateName()))
}

// Cast tries to cast the message of the given command to the given `To`
// type. If the message is not of type `To`, false is returned.
func Cast[Dest, Source any](cmd Command[Source]) (Command[Dest], bool) {
	out, ok := any(cmd.Message()).(Dest)
	if !ok {
		return nil, false
	}
	return New(string(cmd.ID()), out, WithAggregate(cmd.AggregateID(), cmd.AggregateName())), true
}
