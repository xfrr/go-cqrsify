package multierror

import (
	"errors"
	"fmt"
	"strings"
)

// MultiError represents a collection of errors that implements the error interface.
// It's designed to be safe for concurrent use when only reading, but writes should be synchronized.
type MultiError struct {
	errors []error
}

// New creates a new MultiError from the provided errors.
// Nil errors are filtered out automatically.
func New(errs ...error) *MultiError {
	me := &MultiError{}
	me.Append(errs...)
	return me
}

// Error implements the error interface.
// Returns a formatted string containing all error messages.
func (me *MultiError) Error() string {
	if len(me.errors) == 0 {
		return ""
	}

	if len(me.errors) == 1 {
		return me.errors[0].Error()
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d errors occurred:", len(me.errors)))

	for i, err := range me.errors {
		sb.WriteString(fmt.Sprintf("\n\t* [%d] %s", i+1, err.Error()))
	}

	return sb.String()
}

// Append adds one or more errors to the MultiError.
// Nil errors are automatically filtered out.
// If any of the provided errors is itself a MultiError, its errors are flattened.
func (me *MultiError) Append(errs ...error) {
	for _, err := range errs {
		if err == nil {
			continue
		}

		// Flatten nested MultiErrors
		if multiErr, ok := err.(*MultiError); ok {
			me.errors = append(me.errors, multiErr.errors...)
		} else {
			me.errors = append(me.errors, err)
		}
	}
}

// Errors returns a copy of the underlying error slice.
// This prevents external modification of the internal state.
func (me *MultiError) Errors() []error {
	if me == nil || len(me.errors) == 0 {
		return nil
	}

	result := make([]error, len(me.errors))
	copy(result, me.errors)
	return result
}

// Len returns the number of errors contained in the MultiError.
func (me *MultiError) Len() int {
	if me == nil {
		return 0
	}
	return len(me.errors)
}

// Is implements error matching for Go 1.13+ errors.Is functionality.
// It returns true if any contained error matches the target.
func (me *MultiError) Is(target error) bool {
	if me == nil {
		return false
	}

	for _, err := range me.errors {
		if errors.Is(err, target) {
			return true
		}
	}
	return false
}

// As implements error unwrapping for Go 1.13+ errors.As functionality.
// It finds the first error in the chain that matches the target type.
func (me *MultiError) As(target any) bool {
	if me == nil {
		return false
	}

	for _, err := range me.errors {
		if errors.As(err, target) {
			return true
		}
	}
	return false
}

// Unwrap returns the underlying errors for Go 1.20+ multiple error unwrapping.
func (me *MultiError) Unwrap() []error {
	return me.Errors()
}

// ErrorOrNil returns the MultiError if it contains any errors, otherwise returns nil.
// This is useful for idiomatic Go error handling where you want to return nil
// instead of an empty error collection.
func (me *MultiError) ErrorOrNil() error {
	if me == nil || len(me.errors) == 0 {
		return nil
	}
	return me
}

// HasErrors returns true if the MultiError contains any errors.
func (me *MultiError) HasErrors() bool {
	return me != nil && len(me.errors) > 0
}
