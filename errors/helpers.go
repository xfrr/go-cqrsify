package errors

import "errors"

// IsTemporary checks if an error is classified as temporary.
func IsTemporary(err error) bool {
	var c Classifier
	return errors.As(err, &c) && c.Temporary()
}

// IsRetryable checks if an error is retryable.
func IsRetryable(err error) bool {
	var c Classifier
	return errors.As(err, &c) && c.Retryable()
}

// IsPermanent checks if an error is permanent.
func IsPermanent(err error) bool {
	var c Classifier
	return errors.As(err, &c) && c.Permanent()
}
