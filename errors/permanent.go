package errors

// Permanent error is typically used for
// errors that are not recoverable and will
// not change after retries.
// e.g. invalid user input, missing resources, etc.
type PermanentError struct {
	*baseError
}

func (e *PermanentError) Temporary() bool { return false }

func (e *PermanentError) Retryable() bool { return false }

func (e *PermanentError) Permanent() bool { return true }

func NewPermanentError(err error) *PermanentError {
	return &PermanentError{
		baseError: &baseError{
			error: err,
			kind:  "permanent",
		},
	}
}
