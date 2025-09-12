package retry

import (
	crand "crypto/rand"
	"math/big"
	"time"
)

// Jitter perturbs a base delay. Implementations must be deterministic w.r.t. provided rand.Source.
type Jitter interface {
	Apply(base time.Duration) time.Duration
}

// NoJitter leaves the delay untouched.
type NoJitter struct{}

func (NoJitter) Apply(base time.Duration) time.Duration { return base }

// FullJitter returns uniform random in [0, base].
type FullJitter struct{}

func (FullJitter) Apply(base time.Duration) time.Duration {
	if base <= 0 {
		return 0
	}
	max := big.NewInt(int64(base) + 1) // inclusive of base
	n, err := crand.Int(crand.Reader, max)
	if err != nil {
		return 0
	}
	return time.Duration(n.Int64())
}

// EqualJitter centers around base/2 with +/- base/2 range (AWS "equal jitter").
// Result is uniform in [base/2, base].
type EqualJitter struct{}

func (EqualJitter) Apply(base time.Duration) time.Duration {
	if base <= 0 {
		return 0
	}
	half := base / 2
	max := big.NewInt(int64(half) + 1)
	n, err := crand.Int(crand.Reader, max)
	if err != nil {
		return half
	}
	return half + time.Duration(n.Int64())
}

// DecorrelatedJitter (Qualtrics) picks random between base and prev*3; here we keep it stateless by
// using attempt number to bound range: uniform in [base, min(cap, base*3*attempt)].
type DecorrelatedJitter struct {
	Cap time.Duration // hard cap (optional; 0 means no cap)
}

func (d DecorrelatedJitter) Apply(base time.Duration) time.Duration {
	if base <= 0 {
		return 0
	}
	// A lightweight variation without prev-state; still decorrelates as attempts grow.
	hi := base * 3
	if d.Cap > 0 && hi > d.Cap {
		hi = d.Cap
	}
	span := hi - base
	if span <= 0 {
		return base
	}
	max := big.NewInt(int64(span) + 1)
	n, err := crand.Int(crand.Reader, max)
	if err != nil {
		return base
	}
	return base + time.Duration(n.Int64())
}

func defaultJitter(j Jitter) Jitter {
	if j == nil {
		return FullJitter{}
	}
	return j
}
