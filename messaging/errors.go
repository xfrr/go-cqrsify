package messaging

import (
	"errors"
	"fmt"
)

var (
	ErrHandlerNotFound     = errors.New("handler not found")
	ErrMessageIsNotEvent   = errors.New("message is not an event")
	ErrMessageIsNotCommand = errors.New("message is not a command")
	ErrMessageIsNotQuery   = errors.New("message is not a query")
	ErrPublishOnClosedBus  = errors.New("cannot publish on closed bus")
)

type InvalidMessageTypeError struct {
	Expected string
	Actual   string
}

func (e InvalidMessageTypeError) Error() string {
	return fmt.Sprintf("invalid message type: expected %q, got %q", e.Expected, e.Actual)
}

// NoHandlersForMessageError is returned when there are no handlers for a given message type.
type NoHandlersForMessageError struct {
	MessageType string
}

func (e NoHandlersForMessageError) Error() string {
	return "no handlers found for message type: " + e.MessageType
}
