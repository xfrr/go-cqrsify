package retry

import "time"

// Hooks expose observability without coupling to specific logging/metrics libs.
type Hooks struct {
	// OnAttempt is called just before executing fn for attempt i (0-based).
	OnAttempt func(i int)
	// OnRetry is called after a failed attempt that is retryable and will wait/sleep.
	// delay is the planned sleep time (post-jitter).
	OnRetry func(i int, err error, delay time.Duration)
	// OnGiveUp is called when we give up (non-retryable, max attempts, max elapsed, or context cancel).
	OnGiveUp func(i int, finalErr error, cause error)
}
