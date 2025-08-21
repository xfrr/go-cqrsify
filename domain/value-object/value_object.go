package valueobject

// ValueObject defines the interface that all value objects must implement
type ValueObject interface {
	// Equals compares two value objects for equality
	Equals(other ValueObject) bool
	// String returns a string representation of the value object
	String() string
	// Validate ensures the value object is in a valid state
	Validate() error
}
