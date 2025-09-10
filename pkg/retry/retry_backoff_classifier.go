package retry

// Classifier decides if an error is retryable. Return false to stop immediately.
type Classifier interface {
	Retryable(err error) bool
}

// RetryAll retries on any non-nil error (useful default).
type RetryAll struct{}

func (RetryAll) Retryable(err error) bool { return err != nil }

// RetryOn wraps a predicate for flexible classification.
type RetryOn struct{ Predicate func(error) bool }

func (r RetryOn) Retryable(err error) bool { return err != nil && r.Predicate(err) }

func defaultClassifier(c Classifier) Classifier {
	if c == nil {
		return RetryAll{}
	}
	return c
}
