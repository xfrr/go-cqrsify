package command

import "fmt"

var (
	ErrNilHandler  = fmt.Errorf("handler is nil")
	ErrCastContext = fmt.Errorf("failed to cast context")
)

type ErrSubscribeFailed struct {
	err error
}

func (e ErrSubscribeFailed) Error() string {
	return fmt.Sprintf("failed to subscribe: %s", e.err)
}

func (e ErrSubscribeFailed) Wrap(err error) ErrSubscribeFailed {
	e.err = err
	return e
}

func (e ErrSubscribeFailed) Unwrap() error {
	return e.err
}
