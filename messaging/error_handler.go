package messaging

import "log"

// ErrorHandler defines a function type for handling errors that occur during message processing.
type ErrorHandler func(msg Message, err error)

// DefaultErrorHandler is a basic implementation of ErrorHandler that logs the error.
func DefaultErrorHandler(msg Message, err error) {
	log.Printf("error processing message %v: %v", msg, err)
}
