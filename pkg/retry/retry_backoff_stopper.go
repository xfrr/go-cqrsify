package retry

import (
	"context"
	"time"
)

// Stopper can short-circuit retries after a failed attempt (or external signal).
// Return (true, cause) to stop immediately. Cause is included in finalError.Cause.
type Stopper interface {
	ShouldStop(ctx context.Context, attempt int, lastErr error, elapsed time.Duration) (stop bool, cause error)
}

type StopperFunc func(ctx context.Context, attempt int, lastErr error, elapsed time.Duration) (stop bool, cause error)

func (f StopperFunc) ShouldStop(ctx context.Context, attempt int, lastErr error, elapsed time.Duration) (bool, error) {
	return f(ctx, attempt, lastErr, elapsed)
}

// CombineStoppers combines multiple stoppers into one. If any stopper signals to stop, the combined
// stopper will also signal to stop, returning the cause from the first stopper that signaled to stop.
func CombineStoppers(stoppers ...Stopper) Stopper {
	return StopperFunc(func(ctx context.Context, attempt int, lastErr error, elapsed time.Duration) (bool, error) {
		for _, stopper := range stoppers {
			if stop, cause := stopper.ShouldStop(ctx, attempt, lastErr, elapsed); stop {
				return true, cause
			}
		}
		return false, nil
	})
}

// TokenBucketStopper implements a token bucket algorithm to limit the rate of retries.
// It allows a certain number of retries (capacity) in a given time interval (refillInterval).
// If the bucket is empty, it signals to stop retries until tokens are refilled.
type TokenBucketStopper struct {
	capacity       int           // max tokens
	tokens         int           // current tokens
	refillInterval time.Duration // interval to add one token
	lastRefill     time.Time     // last refill time
}

func NewTokenBucketStopper(capacity int, refillInterval time.Duration) *TokenBucketStopper {
	if capacity <= 0 {
		panic("capacity must be greater than 0")
	}
	if refillInterval <= 0 {
		panic("refillInterval must be greater than 0")
	}
	return &TokenBucketStopper{
		capacity:       capacity,
		tokens:         capacity, // start full
		refillInterval: refillInterval,
		lastRefill:     time.Now(),
	}
}

func (t *TokenBucketStopper) ShouldStop(ctx context.Context, attempt int, lastErr error, elapsed time.Duration) (bool, error) {
	now := time.Now()
	// Refill tokens based on elapsed time since last refill
	elapsedSinceLastRefill := now.Sub(t.lastRefill)
	if elapsedSinceLastRefill >= t.refillInterval {
		tokensToAdd := int(elapsedSinceLastRefill / t.refillInterval)
		t.tokens += tokensToAdd
		if t.tokens > t.capacity {
			t.tokens = t.capacity
		}
		t.lastRefill = now
	}

	if t.tokens > 0 {
		t.tokens--
		return false, nil // allow retry
	}

	return true, nil // stop retries
}
