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
	"math/rand"
	"time"
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
	// RandomSource allows deterministic tests. Defaults to crypto-weak math/rand seeded with time.Now.
	RandomSource rand.Source
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
	randSrc     rand.Source
	sleeper     Sleeper
	stopper     Stopper
	hooks       Hooks
}

// New creates a new Retrier with validated defaults.
func New(opts Options) *Retrier {
	r := &Retrier{
		strategy:    defaultStrategy(opts.Strategy),
		jitter:      defaultJitter(opts.Jitter),
		classifier:  defaultClassifier(opts.Classifier),
		maxAttempts: opts.MaxAttempts,
		maxElapsed:  opts.MaxElapsed,
		randSrc:     opts.RandomSource,
		sleeper:     opts.Sleeper,
		hooks:       opts.Hooks,
		stopper:     opts.Stopper,
	}
	if r.randSrc == nil {
		r.randSrc = rand.NewSource(time.Now().UnixNano())
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
	// We snapshot start wall-clock to enforce MaxElapsed.
	start := r.sleeper.Now()
	// Reset strategy in case it's stateful.
	r.strategy.Reset()

	// Per-run RNG; do not share math.Rand across goroutines without a mutex.
	prng := rand.New(r.randSrc)

	var attempt int
	for {
		if r.hooks.OnAttempt != nil {
			r.hooks.OnAttempt(attempt)
		}

		// Execute user function.
		err := fn(ctx)
		if err == nil {
			return nil
		}

		// Classify error.
		if !r.classifier.Retryable(err) {
			if r.hooks.OnGiveUp != nil {
				r.hooks.OnGiveUp(attempt, err, ErrNonRetryable)
			}
			return wrapFinalError(err, ErrNonRetryable, attempt)
		}

		// Respect context cancelation.
		if ctx.Err() != nil {
			if r.hooks.OnGiveUp != nil {
				r.hooks.OnGiveUp(attempt, ctx.Err(), ctx.Err())
			}
			return wrapFinalError(err, ctx.Err(), attempt)
		}

		// Stopper check (after failure, before we proceed)
		if r.stopper != nil {
			elapsed := r.sleeper.Now().Sub(start)
			if stop, cause := r.stopper.ShouldStop(ctx, attempt, err, elapsed); stop {
				if r.hooks.OnGiveUp != nil {
					r.hooks.OnGiveUp(attempt, err, cause)
				}
				return wrapFinalError(err, cause, attempt)
			}
		}

		// Check attempts guard BEFORE sleeping (i.e., attempts counts executions).
		nextAttempt := attempt + 1
		if r.maxAttempts > 0 && nextAttempt >= r.maxAttempts {
			if r.hooks.OnGiveUp != nil {
				r.hooks.OnGiveUp(attempt, err, ErrGiveUp)
			}
			return wrapFinalError(err, ErrGiveUp, attempt)
		}

		// Check elapsed limit BEFORE sleeping.
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

		// Compute delay = jitter(strategy(attempt, err)).
		base := r.strategy.NextDelay(attempt, err)
		delay := r.jitter.Apply(base, prng)

		// Adjust delay if it would exceed MaxElapsed.
		if r.maxElapsed > 0 {
			now := r.sleeper.Now()
			remaining := r.maxElapsed - now.Sub(start)
			if remaining <= 0 {
				if r.hooks.OnGiveUp != nil {
					r.hooks.OnGiveUp(attempt, err, ErrGiveUp)
				}
				return wrapFinalError(err, ErrGiveUp, attempt)
			}
			if delay > remaining {
				delay = remaining
			}
		}

		if r.hooks.OnRetry != nil {
			r.hooks.OnRetry(attempt, err, delay)
		}

		// Sleep (context-aware).
		if delay > 0 {
			if sleepErr := r.sleeper.Sleep(ctx, delay); sleepErr != nil {
				// Context canceled during sleep.
				if r.hooks.OnGiveUp != nil {
					r.hooks.OnGiveUp(attempt, sleepErr, sleepErr)
				}
				return wrapFinalError(err, sleepErr, attempt)
			}
		}

		// Next attempt.
		attempt = nextAttempt
	}
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
