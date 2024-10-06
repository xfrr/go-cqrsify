package cqrs

import "fmt"

// Command is an interface that represents a command request.
type Command interface {
	CommandName() string
}

// getIdentifier returns the unique identifier of the given interface based on its type.
func getIdentifier(t interface{}) string {
	switch t := t.(type) {
	case nil:
		return ""
	case Command:
		return t.CommandName()
	case fmt.Stringer:
		return t.String()
	case fmt.GoStringer:
		return t.GoString()
	default:
		return fmt.Sprintf("%T", t)
	}
}
