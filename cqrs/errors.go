package cqrs

import (
	"errors"
)

var (
	ErrHandlerNotFound          = errors.New("handler not found")
	ErrHandlerAlreadyRegistered = errors.New("handler already registered")
	ErrNilHandler               = errors.New("the provided handler is nil")
	ErrNilBus                   = errors.New("the provided bus is nil")
	ErrBadRequest               = errors.New("bad request")
)
