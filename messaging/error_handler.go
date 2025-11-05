package messaging

import "log"

// ErrorHandler defines a function type for handling errors that occur during message processing.
type ErrorHandler interface {
	Handle(msg Message, err error)
}

// ErrorHandlerFunc is a function type that implements the ErrorHandler interface.
type ErrorHandlerFunc func(msg Message, err error)

// Handle calls the ErrorHandlerFunc with the given message and error.
func (f ErrorHandlerFunc) Handle(msg Message, err error) {
	f(msg, err)
}

// DefaultErrorHandler is a basic implementation of ErrorHandler that logs the error.
var DefaultErrorHandler ErrorHandler = ErrorHandlerFunc(func(msg Message, err error) {
	if msg == nil {
		log.Printf("error processing message: %v", err)
		return
	}

	log.Printf("error processing message %T: %v", msg, err)
})
