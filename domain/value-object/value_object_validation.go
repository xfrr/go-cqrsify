package valueobject

import (
	"errors"
	"fmt"
	"strings"
)

// ValidationError represents validation errors for value objects
type ValidationError struct {
	Field   string
	Message string
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("validation error on field '%s': %s", e.Field, e.Message)
}

// MultiValidationError holds multiple validation errors
type MultiValidationError struct {
	Errors []ValidationError
}

// UnWrap returns the first error if there's only one, nil if none, or itself if multiple
// for compatibility with errors.Is and errors.As.
func (e MultiValidationError) UnWrap() error {
	if len(e.Errors) == 0 {
		return nil
	}
	if len(e.Errors) == 1 {
		return e.Errors[0]
	}
	return e
}

func (e MultiValidationError) Error() string {
	messages := make([]string, 0, len(e.Errors))
	for _, err := range e.Errors {
		messages = append(messages, err.Error())
	}
	return strings.Join(messages, "; ")
}

func ValidationErrors(errs []ValidationError) error {
	if len(errs) == 0 {
		return nil
	}
	return MultiValidationError{Errors: errs}
}

func IsMultiValidationError(err error) bool {
	var multiErr MultiValidationError
	return errors.As(err, &multiErr)
}

func (e MultiValidationError) Contains(target ValidationError) bool {
	for _, err := range e.Errors {
		if errors.Is(err, target) {
			return true
		}
	}
	return false
}
