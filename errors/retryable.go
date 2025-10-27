package errors

// Retryable error is typically used for
// errors that can be retried after a short delay.
// Unlike temporary errors, retryable errors may not
// resolve themselves and require explicit retries.
// e.g. rate limiting, service unavailability, etc.
type RetryableError struct {
	*baseError
}

func (e *RetryableError) Temporary() bool { return false }

func (e *RetryableError) Retryable() bool { return true }

func (e *RetryableError) Permanent() bool { return false }

func NewRetryableError(err error) *RetryableError {
	return &RetryableError{
		baseError: &baseError{
			error: err,
			kind:  "retryable",
		},
	}
}
