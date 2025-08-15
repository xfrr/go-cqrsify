package message

import (
	"errors"
	"fmt"
)

var (
	ErrHandlerNotFound          = errors.New("handler not found")
	ErrHandlerAlreadyRegistered = errors.New("handler already registered")
	ErrNilHandler               = errors.New("the provided handler is nil")
	ErrNilRegisterer            = errors.New("the provided registerer is nil")
	ErrBadRequest               = errors.New("bad request")
)

type InvalidMessageTypeError struct {
	Expected string
	Actual   string
}

func (e InvalidMessageTypeError) Error() string {
	return fmt.Sprintf("invalid message type: expected %q, got %q", e.Expected, e.Actual)
}
