package errors

// Temporary error is typically used for
// transient issues that may resolve themselves
// after a retry or a short period of time.
// e.g. network timeouts, temporary file system issues, etc.
type TemporaryError struct {
	*baseError
}

func (e *TemporaryError) Temporary() bool { return true }

func (e *TemporaryError) Retryable() bool { return true }

func (e *TemporaryError) Permanent() bool { return false }

func NewTemporaryError(err error) *TemporaryError {
	return &TemporaryError{
		baseError: &baseError{
			error: err,
			kind:  "temporary",
		},
	}
}
