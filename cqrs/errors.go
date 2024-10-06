package cqrs

import (
	"errors"
	"fmt"
)

var (
	ErrHandlerNotFound          = errors.New("handler not found")
	ErrHandlerAlreadyRegistered = errors.New("handler already registered")
	ErrNilHandler               = fmt.Errorf("the provided handler is nil")
	ErrNilBus                   = fmt.Errorf("the provided bus is nil")
	ErrHandleFailed             = fmt.Errorf("failed to handle request")
	ErrInvalidRequest           = fmt.Errorf("invalid request")
)
