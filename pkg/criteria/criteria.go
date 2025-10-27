package criteria

import (
	"cmp"
	"strings"
)

// Criteria defines the interface for all criteria implementations
type Criteria[T any] interface {
	// MeetsCriteria checks if the given entities meet this criteria
	MeetsCriteria(entities []T) []T
}

// AndCriteria represents a logical AND operation between multiple criteria
type AndCriteria[T any] struct {
	criteria      Criteria[T]
	otherCriteria Criteria[T]
}

// NewAndCriteria creates a new AND criteria combining two criteria
func NewAndCriteria[T any](criteria, otherCriteria Criteria[T]) *AndCriteria[T] {
	return &AndCriteria[T]{
		criteria:      criteria,
		otherCriteria: otherCriteria,
	}
}

// MeetsCriteria returns entities that meet both criteria
func (ac *AndCriteria[T]) MeetsCriteria(entities []T) []T {
	firstFilter := ac.criteria.MeetsCriteria(entities)
	return ac.otherCriteria.MeetsCriteria(firstFilter)
}

// OrCriteria represents a logical OR operation between multiple criteria
type OrCriteria[T comparable] struct {
	criteria      Criteria[T]
	otherCriteria Criteria[T]
}

// NewOrCriteria creates a new OR criteria combining two criteria
func NewOrCriteria[T comparable](criteria, otherCriteria Criteria[T]) *OrCriteria[T] {
	return &OrCriteria[T]{
		criteria:      criteria,
		otherCriteria: otherCriteria,
	}
}

// MeetsCriteria returns entities that meet either criteria (without duplicates)
func (oc *OrCriteria[T]) MeetsCriteria(entities []T) []T {
	firstSet := oc.criteria.MeetsCriteria(entities)
	secondSet := oc.otherCriteria.MeetsCriteria(entities)

	// For primitive types and simple structs, we can use a map for deduplication
	// For complex types, this approach works with comparable types
	return removeDuplicates[T](append(firstSet, secondSet...))
}

// NotCriteria represents a logical NOT operation for a criteria
type NotCriteria[T any] struct {
	criteria Criteria[T]
}

// NewNotCriteria creates a new NOT criteria that inverts the given criteria
func NewNotCriteria[T any](criteria Criteria[T]) *NotCriteria[T] {
	return &NotCriteria[T]{
		criteria: criteria,
	}
}

// MeetsCriteria returns entities that do NOT meet the wrapped criteria
func (nc *NotCriteria[T]) MeetsCriteria(entities []T) []T {
	meetsCriteria := nc.criteria.MeetsCriteria(entities)
	meetsSet := make(map[any]bool)

	for _, entity := range meetsCriteria {
		meetsSet[entity] = true
	}

	result := make([]T, 0)
	for _, entity := range entities {
		if !meetsSet[entity] {
			result = append(result, entity)
		}
	}

	return result
}

// PredicateCriteria allows using a custom function as criteria
type PredicateCriteria[T any] struct {
	predicate func(T) bool
}

// NewPredicateCriteria creates a new criteria based on a predicate function
func NewPredicateCriteria[T any](predicate func(T) bool) *PredicateCriteria[T] {
	return &PredicateCriteria[T]{
		predicate: predicate,
	}
}

// MeetsCriteria returns entities that satisfy the predicate function
func (pc *PredicateCriteria[T]) MeetsCriteria(entities []T) []T {
	result := make([]T, 0)
	for _, entity := range entities {
		if pc.predicate == nil {
			continue
		}

		if pc.predicate(entity) {
			result = append(result, entity)
		}
	}
	return result
}

// FieldCriteria provides a generic way to filter based on field values using a field accessor function
type FieldCriteria[T any, F any] struct {
	fieldAccessor func(T) F
	value         F
	compareFn     func(F, F) bool
}

// NewFieldCriteria creates a new criteria that filters based on a field value
func NewFieldCriteria[T any, F any](fieldAccessor func(T) F, value F, compareFn func(F, F) bool) *FieldCriteria[T, F] {
	return &FieldCriteria[T, F]{
		fieldAccessor: fieldAccessor,
		value:         value,
		compareFn:     compareFn,
	}
}

// MeetsCriteria returns entities where the specified field meets the comparison criteria
func (fc *FieldCriteria[T, F]) MeetsCriteria(entities []T) []T {
	result := make([]T, 0)

	for _, entity := range entities {
		fieldValue := fc.fieldAccessor(entity)
		if fc.compareFn(fieldValue, fc.value) {
			result = append(result, entity)
		}
	}

	return result
}

// ComparisonFunctions provides type-safe comparison functions
type ComparisonFunctions[T any] struct{}

// Equal returns a comparison function for equality
func (ComparisonFunctions[T]) Equal() func(T, T) bool {
	return func(a, b T) bool {
		return any(a) == any(b)
	}
}

// NotEqual returns a comparison function for inequality
func (ComparisonFunctions[T]) NotEqual() func(T, T) bool {
	return func(a, b T) bool {
		return any(a) != any(b)
	}
}

// OrderedComparisonFunctions provides comparison functions for ordered types
type OrderedComparisonFunctions[T cmp.Ordered] struct{}

// Equal returns a comparison function for equality
func (OrderedComparisonFunctions[T]) Equal() func(T, T) bool {
	return func(a, b T) bool {
		return a == b
	}
}

// NotEqual returns a comparison function for inequality
func (OrderedComparisonFunctions[T]) NotEqual() func(T, T) bool {
	return func(a, b T) bool {
		return a != b
	}
}

// GreaterThan returns a comparison function for greater than
func (OrderedComparisonFunctions[T]) GreaterThan() func(T, T) bool {
	return func(fieldValue, expectedValue T) bool {
		return fieldValue > expectedValue
	}
}

// LessThan returns a comparison function for less than
func (OrderedComparisonFunctions[T]) LessThan() func(T, T) bool {
	return func(fieldValue, expectedValue T) bool {
		return fieldValue < expectedValue
	}
}

// GreaterThanOrEqual returns a comparison function for greater than or equal
func (OrderedComparisonFunctions[T]) GreaterThanOrEqual() func(T, T) bool {
	return func(fieldValue, expectedValue T) bool {
		return fieldValue >= expectedValue
	}
}

// LessThanOrEqual returns a comparison function for less than or equal
func (OrderedComparisonFunctions[T]) LessThanOrEqual() func(T, T) bool {
	return func(fieldValue, expectedValue T) bool {
		return fieldValue <= expectedValue
	}
}

// StringComparisonFunctions provides string-specific comparison functions
type StringComparisonFunctions struct{}

// Contains returns a comparison function that checks if the field contains the expected substring
func (StringComparisonFunctions) Contains() func(string, string) bool {
	return func(fieldValue, expectedValue string) bool {
		return strings.Contains(fieldValue, expectedValue)
	}
}

// HasPrefix returns a comparison function that checks if the field starts with the expected prefix
func (StringComparisonFunctions) HasPrefix() func(string, string) bool {
	return func(fieldValue, expectedValue string) bool {
		return strings.HasPrefix(fieldValue, expectedValue)
	}
}

// HasSuffix returns a comparison function that checks if the field ends with the expected suffix
func (StringComparisonFunctions) HasSuffix() func(string, string) bool {
	return func(fieldValue, expectedValue string) bool {
		return strings.HasSuffix(fieldValue, expectedValue)
	}
}

// EqualFold returns a comparison function for case-insensitive string equality
func (StringComparisonFunctions) EqualFold() func(string, string) bool {
	return func(fieldValue, expectedValue string) bool {
		return strings.EqualFold(fieldValue, expectedValue)
	}
}

// CriteriaBuilder provides a fluent interface for building complex criteria
type CriteriaBuilder[T comparable] struct {
	criteria Criteria[T]
}

// NewCriteriaBuilder creates a new criteria builder
func NewCriteriaBuilder[T comparable]() *CriteriaBuilder[T] {
	return &CriteriaBuilder[T]{}
}

// WithCriteria sets the base criteria
func (cb *CriteriaBuilder[T]) WithCriteria(criteria Criteria[T]) *CriteriaBuilder[T] {
	cb.criteria = criteria
	return cb
}

// And adds an AND criteria
func (cb *CriteriaBuilder[T]) And(criteria Criteria[T]) *CriteriaBuilder[T] {
	if cb.criteria == nil {
		cb.criteria = criteria
	} else {
		cb.criteria = NewAndCriteria(cb.criteria, criteria)
	}
	return cb
}

// Or adds an OR criteria
func (cb *CriteriaBuilder[T]) Or(criteria Criteria[T]) *CriteriaBuilder[T] {
	if cb.criteria == nil {
		cb.criteria = criteria
	} else {
		cb.criteria = NewOrCriteria(cb.criteria, criteria)
	}
	return cb
}

// Not wraps the current criteria in a NOT criteria
func (cb *CriteriaBuilder[T]) Not() *CriteriaBuilder[T] {
	if cb.criteria != nil {
		cb.criteria = NewNotCriteria(cb.criteria)
	}
	return cb
}

// Build returns the final criteria
func (cb *CriteriaBuilder[T]) Build() Criteria[T] {
	return cb.criteria
}

// Helper function to remove duplicates from a slice
func removeDuplicates[T comparable](slice []T) []T {
	keys := make(map[T]bool)
	result := make([]T, 0)

	for _, item := range slice {
		if !keys[item] {
			keys[item] = true
			result = append(result, item)
		}
	}

	return result
}
