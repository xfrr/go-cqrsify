package messaging

import (
	"errors"
	"fmt"
)

var (
	ErrHandlerNotFound = errors.New("handler not found")
)

type InvalidMessageTypeError struct {
	Expected string
	Actual   string
}

func (e InvalidMessageTypeError) Error() string {
	return fmt.Sprintf("invalid message type: expected %q, got %q", e.Expected, e.Actual)
}

// NoSubscribersForMessageError is returned by Publish if no subscribers exist for an message.
type NoSubscribersForMessageError struct {
	MessageName string
}

func (e NoSubscribersForMessageError) Error() string {
	return "messagebus: no subscribers for message " + e.MessageName
}
