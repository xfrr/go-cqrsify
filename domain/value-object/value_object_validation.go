package valueobject

import (
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
