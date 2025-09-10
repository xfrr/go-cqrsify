package retry

import (
	"errors"
	"math"
	"time"
)

// RetryAfterHint is implemented by errors that can surface a server-provided wait.
// (Duration, true) means "wait at least Duration before retrying".
type RetryAfterHint interface {
	RetryAfter() (time.Duration, bool)
}

// Strategy computes the next delay for a given attempt index (0-based).
// Implementations must be deterministic for given inputs and should be stateless.
// Returning a negative or zero duration is allowed and interpreted as "no wait".
type Strategy interface {
	NextDelay(attempt int, prevErr error) time.Duration
	// Optional Reset hook for stateful strategies; no-op for stateless.
	Reset()
}

// ConstantStrategy waits a fixed duration between attempts.
type ConstantStrategy struct{ Delay time.Duration }

func (c ConstantStrategy) NextDelay(_ int, _ error) time.Duration { return c.Delay }
func (c ConstantStrategy) Reset()                                 {}

// ExponentialStrategy grows delay as: min(cap, base * factor^attempt).
// base > 0, factor >= 1.0, cap >= base. Uses overflow-safe exponential growth.
type ExponentialStrategy struct {
	Base   time.Duration // e.g., 50 * time.Millisecond
	Factor float64       // e.g., 2.0
	Cap    time.Duration // e.g., 30 * time.Second
}

func (e ExponentialStrategy) NextDelay(attempt int, _ error) time.Duration {
	if attempt <= 0 {
		if e.Base <= 0 {
			return 0
		}
		return clampDur(e.Base, e.Cap)
	}
	if e.Factor < 1.0 {
		// Defensive: if misconfigured, don't shrink over time.
		return clampDur(e.Base, e.Cap)
	}
	// Use pow with guard; compute as float then clamp.
	base := float64(e.Base)
	delay := base * math.Pow(e.Factor, float64(attempt))
	if delay > float64(e.Cap) {
		return e.Cap
	}
	// Convert with overflow guard.
	if delay > float64(math.MaxInt64) {
		return e.Cap
	}
	return time.Duration(delay)
}
func (e ExponentialStrategy) Reset() {}

// StrategyWithHint decorates a base Strategy by honoring server "Retry-After" hints.
// delay := cap( combine(base, hint) ), where combine defaults to max(base, hint).
type StrategyWithHint struct {
	Base    Strategy
	Cap     time.Duration                                // 0 => no cap
	Extract func(err error) (time.Duration, bool)        // optional custom extractor
	Combine func(base, hint time.Duration) time.Duration // default: max
}

func (s StrategyWithHint) NextDelay(attempt int, prevErr error) time.Duration {
	if s.Base == nil {
		// Defensive default to prevent nil deref
		return 0
	}
	base := s.Base.NextDelay(attempt, prevErr)

	extract := s.Extract
	if extract == nil {
		extract = ExtractRetryAfter // default HTTP extractor + RetryAfterHint interface
	}
	hint, ok := extract(prevErr)
	if ok && hint > 0 {
		combine := s.Combine
		if combine == nil {
			// default: take the safer/larger wait
			combine = func(base, hint time.Duration) time.Duration {
				if hint > base {
					return hint
				}
				return base
			}
		}
		base = combine(base, hint)
	}

	if s.Cap > 0 && base > s.Cap {
		return s.Cap
	}
	return base
}

func (s StrategyWithHint) Reset() {
	if s.Base != nil {
		s.Base.Reset()
	}
}

// ExtractRetryAfter discovers retry-after hints in errors.
// Priority:
// 1) err implements RetryAfterHint
// 2) err wraps *HTTPError (below)
// 3) err wraps *http.Response via HTTPErrorFromResponse helper you may use at call-site
func ExtractRetryAfter(err error) (time.Duration, bool) {
	// 1) RetryAfterHint interface
	var rh RetryAfterHint
	if errors.As(err, &rh) {
		if d, ok := rh.RetryAfter(); ok && d > 0 {
			return d, true
		}
	}
	// 2) *HTTPError wrapper
	var he *HTTPError
	if errors.As(err, &he) {
		return parseRetryAfterHeader(he.Header.Get("Retry-After"), he.ReceivedAt)
	}
	// Best-effort: nothing found
	return 0, false
}

func clampDur(v, cap time.Duration) time.Duration {
	if cap <= 0 || v <= cap {
		return v
	}
	return cap
}

func defaultStrategy(s Strategy) Strategy {
	if s == nil {
		return ExponentialStrategy{
			Base:   50 * time.Millisecond,
			Factor: 2.0,
			Cap:    30 * time.Second,
		}
	}
	return s
}
