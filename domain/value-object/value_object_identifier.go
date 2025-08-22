package valueobject

import (
	"fmt"
	"reflect"
	"strconv"
)

var _ ValueObject = (Identifier[any]{})

// Identifier is the value object that represents a unique identifier.
type Identifier[T comparable] struct {
	value T
}

// NewIdentifier creates a new Identifier value object.
func NewIdentifier[T comparable](value T) Identifier[T] {
	return Identifier[T]{value: value}
}

// String returns the string representation of the Identifier.
func (id Identifier[T]) String() string {
	switch kind := reflect.TypeOf(id.value).Kind(); kind {
	case reflect.String:
		return fmt.Sprintf("%v", id.value)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(reflect.ValueOf(id.value).Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.FormatUint(reflect.ValueOf(id.value).Uint(), 10)
	case reflect.Float32, reflect.Float64:
		return strconv.FormatFloat(reflect.ValueOf(id.value).Float(), 'f', -1, 64)
	default:
		return fmt.Sprintf("%v", id.value)
	}
}

// Equals checks if two Identifier values are equal.
func (id Identifier[T]) Equals(other ValueObject) bool {
	if otherId, ok := other.(Identifier[T]); ok {
		return id.value == otherId.value
	}
	return false
}

// Value returns the underlying value of the Identifier.
func (id Identifier[T]) Value() T {
	return id.value
}

// Validate checks if the Identifier value is valid.
func (id Identifier[T]) Validate() error {
	val := reflect.ValueOf(id.value)
	kind := val.Kind()

	switch kind {
	case reflect.String:
		if val.String() == "" {
			return fmt.Errorf("invalid identifier: %q", val.String())
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if val.Int() == 0 {
			return fmt.Errorf("invalid identifier: %d", val.Int())
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if val.Uint() == 0 {
			return fmt.Errorf("invalid identifier: %d", val.Uint())
		}
	case reflect.Float32, reflect.Float64:
		if val.Float() == 0 {
			return fmt.Errorf("invalid identifier: %f", val.Float())
		}
	default:
		return fmt.Errorf("invalid identifier: %v", id.value)
	}

	return nil
}
