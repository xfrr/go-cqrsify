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

// Command represents a command with an arbitrary payload.
type Command[Payload any] interface {
	// ID returns the command id.
	ID() ID

	// Payload returns the command payload.
	Payload() Payload

	// AggregateID returns the id of the aggregate that the command acts on.
	AggregateID() string

	// AggregateName returns the name of the aggregate that the command acts on.
	AggregateName() string
}

// command is the internal implementation of Command.
type command[Payload any] struct {
	props properties[Payload]
}

// properties contains the fields of a Cmd.
type properties[Payload any] struct {
	ID            ID
	Payload       Payload
	AggregateName string
	AggregateID   string
}

// New returns a new command with the given id, name and payload.
func New[Payload any](id string, payload Payload, opts ...Option) command[Payload] {
	c := command[any]{
		props: properties[any]{
			ID:      ID(id),
			Payload: payload,
		},
	}
	for _, opt := range opts {
		opt(&c)
	}
	return command[Payload]{
		props: properties[Payload]{
			ID:            c.props.ID,
			Payload:       c.props.Payload.(Payload),
			AggregateName: c.props.AggregateName,
			AggregateID:   c.props.AggregateID,
		},
	}
}

// ID returns the command id.
func (c command[Payload]) ID() ID {
	return c.props.ID
}

// Payload returns the command payload.
func (c command[Payload]) Payload() Payload {
	return c.props.Payload
}

// AggregateID returns the id of the aggregate that the command acts on.
func (c command[Payload]) AggregateID() string {
	return c.props.AggregateID
}

// AggregateName returns the name of the aggregate that the command acts on.
func (c command[Payload]) AggregateName() string {
	return c.props.AggregateName
}

// Any returns a new command with the same id and aggregate as the given command
// but with an arbitrary payload.
func (c command[Payload]) Any() Command[any] {
	return New[any](string(c.ID()), c.Payload(), WithAggregate(c.AggregateID(), c.AggregateName()))
}

// Cast tries to cast the payload of the given command to the given `To`
// type. If the payload is not of type `To`, false is returned.
func Cast[Dest, Source any](cmd Command[Source]) (Command[Dest], bool) {
	out, ok := any(cmd.Payload()).(Dest)
	if !ok {
		return nil, false
	}
	return New(string(cmd.ID()), out, WithAggregate(cmd.AggregateID(), cmd.AggregateName())), true
}
