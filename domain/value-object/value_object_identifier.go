package valueobject

import (
	"fmt"
	"reflect"
	"strconv"
)

var _ ValueObject = (Identifier[any]{})

// Identifier is the value object that represents a unique identifier.
type Identifier[T comparable] struct {
	value      T
	validateFn func(id Identifier[any]) error
}

// NewIdentifier creates a new Identifier value object.
func NewIdentifier[T comparable](value T, opts ...IdentifierOption) Identifier[T] {
	options := &IdentifierOptions{}
	options.Apply(opts...)
	return Identifier[T]{value: value, validateFn: options.customValidationFn}
}

// String returns the string representation of the Identifier.
func (id Identifier[T]) String() string {
	switch v := any(id.value).(type) {
	case string:
		return v
	case int:
		return strconv.FormatInt(int64(v), 10)
	case int8:
		return strconv.FormatInt(int64(v), 10)
	case int16:
		return strconv.FormatInt(int64(v), 10)
	case int32:
		return strconv.FormatInt(int64(v), 10)
	case int64:
		return strconv.FormatInt(v, 10)
	case uint:
		return strconv.FormatUint(uint64(v), 10)
	case uint8:
		return strconv.FormatUint(uint64(v), 10)
	case uint16:
		return strconv.FormatUint(uint64(v), 10)
	case uint32:
		return strconv.FormatUint(uint64(v), 10)
	case uint64:
		return strconv.FormatUint(v, 10)
	case float32:
		return strconv.FormatFloat(float64(v), 'f', -1, 32)
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	default:
		return fmt.Sprintf("%v", id.value)
	}
}

// Equals checks if two Identifier values are equal.
func (id Identifier[T]) Equals(other ValueObject) bool {
	if otherID, ok := other.(Identifier[T]); ok {
		return id.value == otherID.value
	}
	return false
}

// Value returns the underlying value of the Identifier.
func (id Identifier[T]) Value() T {
	return id.value
}

// Validate checks if the Identifier value is valid.
func (id Identifier[T]) Validate() error {
	if id.validateFn != nil {
		return id.validateFn(id.any())
	}

	switch val := reflect.ValueOf(id.value); val.Kind() {
	case reflect.Invalid:
		return fmt.Errorf("invalid identifier: %v", id.value)
	case reflect.String:
		if val.String() == "" {
			return fmt.Errorf("invalid identifier: %q", val.String())
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if val.Int() == 0 {
			return fmt.Errorf("invalid identifier: %d", val.Int())
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		if val.Uint() == 0 {
			return fmt.Errorf("invalid identifier: %d", val.Uint())
		}
	case reflect.Float32, reflect.Float64:
		if val.Float() == 0 {
			return fmt.Errorf("invalid identifier: %f", val.Float())
		}
	case reflect.Bool, reflect.Complex64, reflect.Complex128, reflect.Array, reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Ptr, reflect.Slice, reflect.Struct, reflect.UnsafePointer:
		return fmt.Errorf("invalid identifier: %v", id.value)
	default:
		return fmt.Errorf("invalid identifier: %v", id.value)
	}

	return nil
}

func (id Identifier[T]) any() Identifier[any] {
	return Identifier[any]{value: id.value, validateFn: id.validateFn}
}
