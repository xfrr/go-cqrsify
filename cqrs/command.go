package cqrs

// Command is an interface that represents a command request.
type Command interface {
	CommandName() string
}
