package retry

import (
	"context"
	"errors"
	"sync"
	"time"
)

// ErrBackpressure is returned by TokenBucketStopper to explain why retries stopped.
var ErrBackpressure = errors.New("retry: backpressure (token bucket empty)")

// TokenBucketStopper halts retries while a shared token bucket is empty.
// Typical use: wire this into Retrier.Options.Stopper and consume/add tokens
// from an external SLO monitor (e.g., p99 latency breach burns tokens; recovery adds them).
//
// Thread-safety: all public methods are safe for concurrent use.
type TokenBucketStopper struct {
	mu sync.Mutex
	// configuration
	capacity   float64       // maximum tokens (burst)
	refillRate float64       // tokens per second
	minDelay   time.Duration // advisory only; align Strategy caps if desired

	// state
	tokens float64   // current tokens
	last   time.Time // last refill instant

	// behavior
	consumePerFailure bool             // optional: decrement token on each failed attempt
	clock             func() time.Time // injectable clock for tests
}

// NewTokenBucketStopper creates a stopper with a classic token bucket.
//
//   - capacity: maximum tokens (burst). If <= 0, defaults to 1.
//   - refillRate: tokens added per second. If < 0, treated as 0.
//   - startFull: if true, bucket starts at capacity; otherwise 0.
//   - minDelay: optional advisory minimum delay (you can read it to align Strategy caps).
func NewTokenBucketStopper(capacity, refillRate float64, startFull bool, minDelay time.Duration) *TokenBucketStopper {
	if capacity <= 0 {
		capacity = 1
	}
	if refillRate < 0 {
		refillRate = 0
	}
	now := time.Now()
	tb := &TokenBucketStopper{
		mu:                sync.Mutex{},
		capacity:          capacity,
		refillRate:        refillRate,
		minDelay:          minDelay,
		last:              now,
		clock:             time.Now,
		tokens:            0,
		consumePerFailure: false,
	}
	if startFull {
		tb.tokens = capacity
	}
	return tb
}

// WithClock overrides the internal clock (useful for unit tests).
func (t *TokenBucketStopper) WithClock(f func() time.Time) *TokenBucketStopper {
	t.mu.Lock()
	defer t.mu.Unlock()
	if f == nil {
		t.clock = time.Now
	} else {
		// keep state consistent on clock swap
		t.refillUnlocked()
		t.clock = f
		t.last = t.clock()
	}
	return t
}

// SetConsumePerFailure enables/disables consuming one token per failed attempt inspected by Stopper.
// When enabled, the bucket drains proportional to failure density even without external signals.
func (t *TokenBucketStopper) SetConsumePerFailure(on bool) {
	t.mu.Lock()
	t.consumePerFailure = on
	t.mu.Unlock()
}

// Reconfigure updates capacity and refillRate at runtime.
// It first applies a refill at the current rate, then updates the config.
func (t *TokenBucketStopper) Reconfigure(capacity, refillRate float64) {
	if capacity <= 0 {
		capacity = 1
	}
	if refillRate < 0 {
		refillRate = 0
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	t.refillUnlocked()
	t.capacity = capacity
	if t.tokens > t.capacity {
		t.tokens = t.capacity
	}
	t.refillRate = refillRate
}

// Consume burns n tokens immediately (n > 0). Safe to call from external SLO controllers.
func (t *TokenBucketStopper) Consume(n float64) {
	if n <= 0 {
		return
	}
	t.mu.Lock()
	t.refillUnlocked()
	t.tokens -= n
	if t.tokens < 0 {
		t.tokens = 0
	}
	t.mu.Unlock()
}

// AddTokens adds up to n tokens immediately (bounded by capacity). Useful when signals clear.
func (t *TokenBucketStopper) AddTokens(n float64) {
	if n <= 0 {
		return
	}
	t.mu.Lock()
	t.refillUnlocked()
	t.tokens += n
	if t.tokens > t.capacity {
		t.tokens = t.capacity
	}
	t.mu.Unlock()
}

// MinDelay returns the advisory minimum delay (you can use it to align Strategy caps).
func (t *TokenBucketStopper) MinDelay() time.Duration {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.minDelay
}

// ShouldStop implements Stopper. It runs after a failed attempt and before the backoff sleep.
// If there are no tokens available, it stops retries immediately with ErrBackpressure.
// If consumePerFailure is enabled and tokens are available, it decrements one token per failure.
func (t *TokenBucketStopper) ShouldStop(_ context.Context, _ int, _ error, _ time.Duration) (bool, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.refillUnlocked()

	if t.consumePerFailure {
		if t.tokens >= 1.0 {
			t.tokens -= 1.0
		}
	}

	// If empty, instruct Retrier to stop right now.
	if t.tokens < 1.0 {
		return true, ErrBackpressure
	}
	// Otherwise, allow retry to proceed.
	return false, nil
}

// refillUnlocked performs time-based refill at the configured rate.
// Assumes caller holds t.mu.
func (t *TokenBucketStopper) refillUnlocked() {
	now := t.clock()
	dt := now.Sub(t.last).Seconds()
	if dt > 0 && t.refillRate > 0 {
		t.tokens += dt * t.refillRate
		if t.tokens > t.capacity {
			t.tokens = t.capacity
		}
	}
	t.last = now
}
