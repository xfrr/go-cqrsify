package valueobject

import (
	"reflect"
)

// BaseValueObject provides common functionality for all value objects
type BaseValueObject struct{}

// Equals implements structural equality comparison using reflection
func (bvo BaseValueObject) Equals(other ValueObject) bool {
	if other == nil {
		return false
	}

	thisValue := reflect.ValueOf(bvo)
	otherValue := reflect.ValueOf(other)

	// Check if types are the same
	if thisValue.Type() != otherValue.Type() {
		return false
	}

	return reflect.DeepEqual(bvo, other)
}
