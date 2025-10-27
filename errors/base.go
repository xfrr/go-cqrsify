package errors

type baseError struct {
	error
	kind string
}

func (e *baseError) Error() string {
	if e.error != nil {
		return e.kind + " error: " + e.error.Error()
	}
	return e.kind + " error"
}

func (e *baseError) Unwrap() error { return e.error }
