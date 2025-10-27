package errors

// Classifier is an interface for classifying errors.
type Classifier interface {
	error
	Temporary() bool
	Retryable() bool
	Permanent() bool
	Unwrap() error
}
