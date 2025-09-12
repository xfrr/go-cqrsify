package retry

// Classifier decides if an error is retryable. Return false to stop immediately.
type Classifier interface {
	Retryable(err error) bool
}

// RetryAlways retries on any non-nil error (useful default).
//
//nolint:revive // allow name
type RetryAlways struct{}

func (RetryAlways) Retryable(err error) bool { return err != nil }

// RetryOn wraps a predicate for flexible classification.
//
//nolint:revive // allow name
type RetryOn struct{ Predicate func(error) bool }

func (r RetryOn) Retryable(err error) bool { return err != nil && r.Predicate(err) }

func defaultClassifier(c Classifier) Classifier {
	if c == nil {
		return RetryAlways{}
	}
	return c
}
