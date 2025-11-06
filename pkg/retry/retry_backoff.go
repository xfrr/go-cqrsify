// Package retry provides composable, production-grade retry with backoff and jitter.
//
// Design overview (SOLID, Clean Architecture):
// - Strategy (interface) encapsulates the backoff policy (SRP). Built-in: Constant, Exponential.
// - Jitter (interface) decorates delay distribution with randomness; built-in: None, Full, Equal, Decorrelated.
// - Classifier (interface) decides whether an error is retryable (Open/Closed principle).
// - Sleeper (interface) abstracts time.Sleep for testability (Dependency Inversion).
// - Hooks allow observability without coupling to logging/metrics providers (DIP).
// - Retrier orchestrates attempts with context, deadlines, max-attempts/max-elapsed guards.
//
// Thread-safety: Retrier is immutable after creation and safe for concurrent use.
// Strategy instances are NOT required to be goroutine-safe; Retrier clones or uses stateless strategies per run.
package retry

import (
	"context"
	"errors"
	"fmt"
	"time"
)

const (
	defaultMaxAttempts = 10
	defaultMaxElapsed  = 0 // no limit
)

var (
	// ErrGiveUp is returned when retries exhausted (max attempts or elapsed).
	ErrGiveUp = errors.New("retry: give up")
	// ErrNonRetryable is returned when classifier decides the error must not be retried.
	ErrNonRetryable = errors.New("retry: non-retryable")
)

// Options configures a Retrier.
type Options struct {
	// Strategy defines base backoff per attempt.
	Strategy Strategy
	// Jitter decorates Strategy delays. Defaults to NoJitter.
	Jitter Jitter
	// Classifier determines retryability. Defaults to RetryAll.
	Classifier Classifier
	// Stopper can short-circuit retries (circuit breaker / kill-switch / budget guard).
	// If both Stopper and Classifier are set, Stopper is consulted first.
	Stopper Stopper
	// MaxAttempts caps the number of tries. 0/negative means unlimited (bounded by MaxElapsed/ctx).
	MaxAttempts int
	// MaxElapsed is a wall-clock limit across all attempts including sleeps. 0 means no limit.
	MaxElapsed time.Duration
	// Sleeper allows injecting fake/time-travel clocks for tests. Defaults to RealSleeper.
	Sleeper Sleeper
	// Hooks (optional) for observability.
	Hooks Hooks
}

// Retrier orchestrates retrying a function with the configured backoff and policies.
type Retrier struct {
	strategy    Strategy
	jitter      Jitter
	classifier  Classifier
	maxAttempts int
	maxElapsed  time.Duration
	sleeper     Sleeper
	stopper     Stopper
	hooks       Hooks
}

func DefaultOptions() Options {
	return Options{
		Strategy:    defaultStrategy(nil),
		Jitter:      defaultJitter(nil),
		Classifier:  defaultClassifier(nil),
		Sleeper:     defaultSleeper(nil),
		Hooks:       Hooks{},
		Stopper:     defaultStopper(nil),
		MaxAttempts: defaultMaxAttempts,
		MaxElapsed:  defaultMaxElapsed,
	}
}

// New creates a new Retrier with validated defaults.
func New(opts Options) *Retrier {
	r := &Retrier{
		strategy:    defaultStrategy(opts.Strategy),
		jitter:      defaultJitter(opts.Jitter),
		classifier:  defaultClassifier(opts.Classifier),
		maxAttempts: opts.MaxAttempts,
		maxElapsed:  opts.MaxElapsed,
		sleeper:     opts.Sleeper,
		hooks:       opts.Hooks,
		stopper:     opts.Stopper,
	}
	if r.sleeper == nil {
		r.sleeper = RealSleeper{}
	}
	return r
}

// Do executes fn, retrying according to configuration.
// It returns nil on success, or a wrapped error on failure indicating cause (give up vs non-retryable vs ctx).
//
// Contract of fn:
// - Must be idempotent or otherwise safe to run multiple times.
// - Must be quick to return; long sleeps should be avoided (backoff handles pacing).
func (r *Retrier) Do(ctx context.Context, fn func(ctx context.Context) error) error {
	start := r.sleeper.Now()
	r.strategy.Reset()

	var attempt int
	for {
		if r.hooks.OnAttempt != nil {
			r.hooks.OnAttempt(attempt)
		}

		err := fn(ctx)
		if err == nil {
			return nil
		}

		if finalErr := r.handleNonRetryable(attempt, err); finalErr != nil {
			return finalErr
		}

		if finalErr := r.handleContextCancel(ctx, attempt, err); finalErr != nil {
			return finalErr
		}

		if finalErr := r.handleStopper(ctx, attempt, err, start); finalErr != nil {
			return finalErr
		}

		nextAttempt := attempt + 1
		if finalErr := r.handleMaxAttempts(nextAttempt, attempt, err); finalErr != nil {
			return finalErr
		}

		if finalErr := r.handleMaxElapsedBeforeSleep(attempt, err, start); finalErr != nil {
			return finalErr
		}

		delay := r.computeDelay(attempt, err)
		delay = r.adjustDelayForMaxElapsed(delay, attempt, err, start)
		if delay == -1 {
			return wrapFinalError(err, ErrGiveUp, attempt)
		}

		if r.hooks.OnRetry != nil {
			r.hooks.OnRetry(attempt, err, delay)
		}

		if finalErr := r.handleSleep(ctx, attempt, err, delay); finalErr != nil {
			return finalErr
		}

		attempt = nextAttempt
	}
}

func (r *Retrier) handleNonRetryable(attempt int, err error) error {
	if !r.classifier.Retryable(err) {
		if r.hooks.OnGiveUp != nil {
			r.hooks.OnGiveUp(attempt, err, ErrNonRetryable)
		}
		return wrapFinalError(err, ErrNonRetryable, attempt)
	}
	return nil
}

func (r *Retrier) handleContextCancel(ctx context.Context, attempt int, err error) error {
	if ctx.Err() != nil {
		if r.hooks.OnGiveUp != nil {
			r.hooks.OnGiveUp(attempt, ctx.Err(), ctx.Err())
		}
		return wrapFinalError(err, ctx.Err(), attempt)
	}
	return nil
}

func (r *Retrier) handleStopper(ctx context.Context, attempt int, err error, start time.Time) error {
	if r.stopper != nil {
		elapsed := r.sleeper.Now().Sub(start)
		if stop, cause := r.stopper.ShouldStop(ctx, attempt, err, elapsed); stop {
			if r.hooks.OnGiveUp != nil {
				r.hooks.OnGiveUp(attempt, err, cause)
			}
			return wrapFinalError(err, cause, attempt)
		}
	}
	return nil
}

func (r *Retrier) handleMaxAttempts(nextAttempt, attempt int, err error) error {
	if r.maxAttempts > 0 && nextAttempt >= r.maxAttempts {
		if r.hooks.OnGiveUp != nil {
			r.hooks.OnGiveUp(attempt, err, ErrGiveUp)
		}
		return wrapFinalError(err, ErrGiveUp, attempt)
	}
	return nil
}

func (r *Retrier) handleMaxElapsedBeforeSleep(attempt int, err error, start time.Time) error {
	if r.maxElapsed > 0 {
		now := r.sleeper.Now()
		elapsed := now.Sub(start)
		if elapsed >= r.maxElapsed {
			if r.hooks.OnGiveUp != nil {
				r.hooks.OnGiveUp(attempt, err, ErrGiveUp)
			}
			return wrapFinalError(err, ErrGiveUp, attempt)
		}
	}
	return nil
}

func (r *Retrier) computeDelay(attempt int, err error) time.Duration {
	base := r.strategy.NextDelay(attempt, err)
	return r.jitter.Apply(base)
}

func (r *Retrier) adjustDelayForMaxElapsed(delay time.Duration, attempt int, err error, start time.Time) time.Duration {
	if r.maxElapsed > 0 {
		now := r.sleeper.Now()
		remaining := r.maxElapsed - now.Sub(start)
		if remaining <= 0 {
			if r.hooks.OnGiveUp != nil {
				r.hooks.OnGiveUp(attempt, err, ErrGiveUp)
			}
			return -1
		}
		if delay > remaining {
			delay = remaining
		}
	}
	return delay
}

func (r *Retrier) handleSleep(ctx context.Context, attempt int, err error, delay time.Duration) error {
	if delay > 0 {
		if sleepErr := r.sleeper.Sleep(ctx, delay); sleepErr != nil {
			if r.hooks.OnGiveUp != nil {
				r.hooks.OnGiveUp(attempt, sleepErr, sleepErr)
			}
			return wrapFinalError(err, sleepErr, attempt)
		}
	}
	return nil
}

// DoResult behaves like Do but returns a value on success. On failure, it returns
// the zero value of T and a wrapped error (finalError).
func DoResult[T any](ctx context.Context, r *Retrier, fn func(ctx context.Context) (T, error)) (T, error) {
	var zero T
	var lastErr error

	err := r.Do(ctx, func(ctx context.Context) error {
		var err error
		var res T
		res, err = fn(ctx)
		if err == nil {
			// Success: store result in closure var and return nil error to stop retries.
			zero = res
			return nil
		}
		// Failure: store last error and return it to continue retries.
		lastErr = err
		return err
	})
	if err != nil {
		// Return the last error from fn, wrapped in finalError.
		return zero, wrapFinalError(lastErr, Cause(err), Attempts(err)-1)
	}
	return zero, nil
}

type finalError struct {
	Err         error // original last error
	Cause       error // why we stopped: ErrNonRetryable, ErrGiveUp, context error
	LastAttempt int   // 0-based index of last attempt executed
}

func (f finalError) Error() string {
	return fmt.Sprintf("retry stopped: cause=%v attempts=%d lastErr=%v", f.Cause, f.LastAttempt+1, f.Err)
}
func (f finalError) Unwrap() error { return f.Err }

func wrapFinalError(last error, cause error, attempt int) error {
	return finalError{Err: last, Cause: cause, LastAttempt: attempt}
}

// Cause extracts the reason why retries stopped from err, or nil if not a finalError.
func Cause(err error) error {
	var fe finalError
	if errors.As(err, &fe) {
		return fe.Cause
	}
	return nil
}

// LastError extracts the last error returned by fn from err, or nil if not a finalError.
func LastError(err error) error {
	var fe finalError
	if errors.As(err, &fe) {
		return fe.Err
	}
	return nil
}

// Attempts extracts the number of attempts executed from err, or 0 if not a finalError.
func Attempts(err error) int {
	var fe finalError
	if errors.As(err, &fe) {
		return fe.LastAttempt + 1
	}
	return 0
}
