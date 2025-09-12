package retry_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xfrr/go-cqrsify/pkg/retry"
)

type fakeStrategy struct {
	delays     []time.Duration
	resetCount int
}

func (f *fakeStrategy) Reset() { f.resetCount++ }
func (f *fakeStrategy) NextDelay(attempt int, _ error) time.Duration {
	if attempt >= 0 && attempt < len(f.delays) {
		return f.delays[attempt]
	}
	if len(f.delays) > 0 {
		return f.delays[len(f.delays)-1]
	}
	return 0
}

type identJitter struct{}

func (identJitter) Apply(d time.Duration) time.Duration { return d }

type alwaysRetryClassifier struct{}

func (alwaysRetryClassifier) Retryable(error) bool { return true }

type neverRetryClassifier struct{}

func (neverRetryClassifier) Retryable(error) bool { return false }

// fakeSleeper advances a logical clock and records sleeps.
// If sleepErr is set, Sleep returns that error.
type fakeSleeper struct {
	now      time.Time
	slept    []time.Duration
	sleepErr error
}

func (s *fakeSleeper) Now() time.Time { return s.now }

func (s *fakeSleeper) Sleep(ctx context.Context, d time.Duration) error {
	s.slept = append(s.slept, d)
	// If the context is already canceled before/while sleeping, return that error.
	if err := ctx.Err(); err != nil {
		return err
	}
	// Advance time deterministically.
	s.now = s.now.Add(d)
	if s.sleepErr != nil {
		return s.sleepErr
	}
	return nil
}

// fakeStopper can stop at or after a given attempt index (0-based).
type fakeStopper struct {
	stopFromAttempt *int
	cause           error
	calls           []int
}

func (s *fakeStopper) ShouldStop(_ context.Context, attempt int, _ error, _ time.Duration) (bool, error) {
	s.calls = append(s.calls, attempt)
	if s.stopFromAttempt != nil && attempt >= *s.stopFromAttempt {
		return true, s.cause
	}
	return false, nil
}

func TestDo_SuccessFirstTry(t *testing.T) {
	fs := &fakeSleeper{now: time.Unix(0, 0)}
	strat := &fakeStrategy{delays: []time.Duration{50 * time.Millisecond}}
	var attempts []int
	var retries []int
	var giveUps []int

	r := retry.New(retry.Options{
		Strategy:   strat,
		Jitter:     identJitter{},
		Classifier: alwaysRetryClassifier{},
		Sleeper:    fs,
		Hooks: retry.Hooks{
			OnAttempt: func(n int) { attempts = append(attempts, n) },
			OnRetry:   func(n int, _ error, _ time.Duration) { retries = append(retries, n) },
			OnGiveUp:  func(n int, _ error, _ error) { giveUps = append(giveUps, n) },
		},
	})

	err := r.Do(context.Background(), func(_ context.Context) error {
		return nil
	})
	require.NoError(t, err)
	assert.Equal(t, []int{0}, attempts, "one attempt (0)")
	assert.Empty(t, retries, "no retries")
	assert.Empty(t, giveUps, "no give-up")
	assert.Equal(t, 1, strat.resetCount, "strategy reset once")
	assert.Empty(t, fs.slept, "no sleeping on success")
}

func TestDo_RetryableThenSuccess(t *testing.T) {
	fs := &fakeSleeper{now: time.Unix(0, 0)}
	strat := &fakeStrategy{delays: []time.Duration{10 * time.Millisecond, 20 * time.Millisecond}}
	var attempts, retries []int

	r := retry.New(retry.Options{
		Strategy:   strat,
		Jitter:     identJitter{},
		Classifier: alwaysRetryClassifier{},
		Sleeper:    fs,
		Hooks: retry.Hooks{
			OnAttempt: func(n int) { attempts = append(attempts, n) },
			OnRetry:   func(n int, _ error, _ time.Duration) { retries = append(retries, n) },
		},
	})

	call := 0
	err := r.Do(context.Background(), func(_ context.Context) error {
		defer func() { call++ }()
		if call < 1 {
			return errors.New("temporary")
		}
		return nil
	})
	require.NoError(t, err)

	assert.Equal(t, []int{0, 1}, attempts)
	assert.Equal(t, []int{0}, retries, "one retry (after attempt 0)")
	require.Len(t, fs.slept, 1)
	assert.Equal(t, 10*time.Millisecond, fs.slept[0])
}

func TestDo_NonRetryableStopsImmediately(t *testing.T) {
	fs := &fakeSleeper{now: time.Unix(0, 0)}
	r := retry.New(retry.Options{
		Strategy:   &fakeStrategy{delays: []time.Duration{100 * time.Millisecond}},
		Jitter:     identJitter{},
		Classifier: neverRetryClassifier{},
		Sleeper:    fs,
	})

	opErr := errors.New("boom")
	err := r.Do(context.Background(), func(_ context.Context) error { return opErr })
	require.Error(t, err)

	require.ErrorIs(t, retry.Cause(err), retry.ErrNonRetryable)
	require.ErrorIs(t, retry.LastError(err), opErr)
	assert.Equal(t, 1, retry.Attempts(err))
	assert.Empty(t, fs.slept, "no sleeps when non-retryable")
}

func TestDo_ContextCanceledBetweenAttempts(t *testing.T) {
	fs := &fakeSleeper{now: time.Unix(0, 0)}
	r := retry.New(retry.Options{
		Strategy:   &fakeStrategy{delays: []time.Duration{5 * time.Second}},
		Jitter:     identJitter{},
		Classifier: alwaysRetryClassifier{},
		Sleeper:    fs,
	})

	ctx, cancel := context.WithCancel(context.Background())

	call := 0
	err := r.Do(ctx, func(_ context.Context) error {
		call++
		// Cancel after first failing return to hit handleContextCancel next loop
		if call == 1 {
			cancel()
			return errors.New("retryable")
		}
		return nil
	})
	require.Error(t, err)
	require.ErrorIs(t, retry.Cause(err), context.Canceled)
	assert.Equal(t, 1, retry.Attempts(err))
	assert.Empty(t, fs.slept, "no sleep because canceled before sleeping")
}

func TestDo_StopperShortCircuits(t *testing.T) {
	fs := &fakeSleeper{now: time.Unix(0, 0)}
	stopCause := errors.New("circuit-open")
	zero := 0
	stopper := &fakeStopper{stopFromAttempt: &zero, cause: stopCause}

	var giveUpCalled bool
	r := retry.New(retry.Options{
		Strategy:   &fakeStrategy{delays: []time.Duration{10 * time.Millisecond}},
		Jitter:     identJitter{},
		Classifier: alwaysRetryClassifier{},
		Sleeper:    fs,
		Stopper:    stopper,
		Hooks: retry.Hooks{
			OnGiveUp: func(_ int, _ error, cause error) {
				giveUpCalled = true
				assert.ErrorIs(t, cause, stopCause)
			},
		},
	})

	opErr := errors.New("x")
	err := r.Do(context.Background(), func(_ context.Context) error { return opErr })
	require.Error(t, err)
	assert.True(t, giveUpCalled)
	require.ErrorIs(t, retry.Cause(err), stopCause)
	require.ErrorIs(t, retry.LastError(err), opErr)
	assert.Equal(t, 1, retry.Attempts(err))
	assert.Empty(t, fs.slept, "stopped before any sleep")
}

func TestDo_MaxAttempts(t *testing.T) {
	fs := &fakeSleeper{now: time.Unix(0, 0)}
	r := retry.New(retry.Options{
		Strategy:    &fakeStrategy{delays: []time.Duration{1, 1, 1}},
		Jitter:      identJitter{},
		Classifier:  alwaysRetryClassifier{},
		Sleeper:     fs,
		MaxAttempts: 3, // allow attempts 0,1,2 then stop
	})

	opErr := errors.New("keep failing")
	err := r.Do(context.Background(), func(_ context.Context) error { return opErr })
	require.Error(t, err)

	require.ErrorIs(t, retry.Cause(err), retry.ErrGiveUp)
	require.ErrorIs(t, retry.LastError(err), opErr)
	assert.Equal(t, 3, retry.Attempts(err))
	// It should have slept after attempt 0 and 1 (but not after last attempt 2)
	require.Len(t, fs.slept, 2)
}

func TestDo_MaxElapsed_BeforeSleepTriggersGiveUp(t *testing.T) {
	fs := &fakeSleeper{now: time.Unix(0, 0)}
	r := retry.New(retry.Options{
		Strategy:   &fakeStrategy{delays: []time.Duration{10 * time.Second}},
		Jitter:     identJitter{},
		Classifier: alwaysRetryClassifier{},
		Sleeper:    fs,
		MaxElapsed: 1 * time.Nanosecond,
	})

	// Make Now() report time >= start + MaxElapsed right after first failure.
	fs.now = time.Unix(0, 1) // 1ns

	opErr := errors.New("fail")
	err := r.Do(context.Background(), func(_ context.Context) error { return opErr })
	require.Error(t, err)

	require.ErrorIs(t, retry.Cause(err), retry.ErrGiveUp)
	assert.Equal(t, 2, retry.Attempts(err))
}

func TestDo_MaxElapsed_AdjustsDelayToRemainingAndSleepsCapped(t *testing.T) {
	start := time.Unix(1000, 0)
	fs := &fakeSleeper{now: start}
	// Big base delay that should get capped by remaining time
	strat := &fakeStrategy{delays: []time.Duration{10 * time.Second}}
	r := retry.New(retry.Options{
		Strategy:   strat,
		Jitter:     identJitter{},
		Classifier: alwaysRetryClassifier{},
		Sleeper:    fs,
		MaxElapsed: 1500 * time.Millisecond, // total budget
	})

	// First failure at t=1000s, remaining = 1.5s, so sleep should cap to 1.5s
	failOnce := true
	err := r.Do(context.Background(), func(_ context.Context) error {
		if failOnce {
			failOnce = false
			return errors.New("retryable")
		}
		return nil
	})
	require.NoError(t, err)

	require.Len(t, fs.slept, 1)
	assert.Equal(t, 1500*time.Millisecond, fs.slept[0], "delay capped to remaining")
}

func TestDo_MaxElapsed_RemainingZeroReturnsGiveUp(t *testing.T) {
	start := time.Unix(0, 0)
	fs := &fakeSleeper{now: start}
	r := retry.New(retry.Options{
		Strategy:   &fakeStrategy{delays: []time.Duration{1 * time.Second}},
		Jitter:     identJitter{},
		Classifier: alwaysRetryClassifier{},
		Sleeper:    fs,
		MaxElapsed: 1 * time.Second,
	})

	// First attempt fails, then before computing delay we move time to exactly start+MaxElapsed
	opErr := errors.New("x")
	move := true
	err := r.Do(context.Background(), func(_ context.Context) error {
		if move {
			move = false
			// Simulate work consuming entire budget instantly
			fs.now = start.Add(1 * time.Second)
			return opErr
		}
		return nil
	})
	require.Error(t, err)
	require.ErrorIs(t, retry.Cause(err), retry.ErrGiveUp)
	assert.Equal(t, 1, retry.Attempts(err))
	assert.Empty(t, fs.slept)
}

func TestDo_SleepErrorIsFinalCause(t *testing.T) {
	fs := &fakeSleeper{now: time.Unix(0, 0), sleepErr: errors.New("sleep aborted")}
	r := retry.New(retry.Options{
		Strategy:   &fakeStrategy{delays: []time.Duration{123 * time.Millisecond}},
		Jitter:     identJitter{},
		Classifier: alwaysRetryClassifier{},
		Sleeper:    fs,
	})

	opErr := errors.New("retryable")
	err := r.Do(context.Background(), func(_ context.Context) error { return opErr })
	require.Error(t, err)

	require.ErrorIs(t, retry.Cause(err), fs.sleepErr)
	require.ErrorIs(t, retry.LastError(err), opErr)
	assert.Equal(t, 1, retry.Attempts(err))
	require.Len(t, fs.slept, 1)
	assert.Equal(t, 123*time.Millisecond, fs.slept[0])
}

func TestDoResult_SuccessReturnsValue(t *testing.T) {
	fs := &fakeSleeper{now: time.Unix(0, 0)}
	r := retry.New(retry.Options{
		Strategy:   &fakeStrategy{delays: []time.Duration{1, 1}},
		Jitter:     identJitter{},
		Classifier: alwaysRetryClassifier{},
		Sleeper:    fs,
	})

	call := 0
	val, err := retry.DoResult(context.Background(), r, func(_ context.Context) (int, error) {
		call++
		if call < 3 {
			return 0, errors.New("again")
		}
		return 42, nil
	})
	require.NoError(t, err)
	assert.Equal(t, 42, val)
	require.Len(t, fs.slept, 2, "slept between the two failures")
}

func TestDoResult_FailurePropagatesFinalErrorAndZeroValue(t *testing.T) {
	fs := &fakeSleeper{now: time.Unix(0, 0)}
	r := retry.New(retry.Options{
		Strategy:    &fakeStrategy{delays: []time.Duration{1, 1}},
		Jitter:      identJitter{},
		Classifier:  alwaysRetryClassifier{},
		Sleeper:     fs,
		MaxAttempts: 2, // attempts 0 and 1, then stop
	})

	last := errors.New("last-err")
	call := 0
	v, err := retry.DoResult(context.Background(), r, func(_ context.Context) (string, error) {
		call++
		return "", last
	})
	require.Error(t, err)
	assert.Empty(t, v, "zero value on failure")
	require.ErrorIs(t, retry.LastError(err), last)
	require.ErrorIs(t, retry.Cause(err), retry.ErrGiveUp)
	assert.Equal(t, 2, retry.Attempts(err))
}

func TestHooks_AreCalledWithExpectedSemantics(t *testing.T) {
	fs := &fakeSleeper{now: time.Unix(0, 0)}
	var attempts []int
	var retries []int
	var giveUps []int

	r := retry.New(retry.Options{
		Strategy:    &fakeStrategy{delays: []time.Duration{5, 5, 5}},
		Jitter:      identJitter{},
		Classifier:  alwaysRetryClassifier{},
		Sleeper:     fs,
		MaxAttempts: 2, // attempts 0 and 1, then give up
		Hooks: retry.Hooks{
			OnAttempt: func(n int) { attempts = append(attempts, n) },
			OnRetry:   func(n int, _ error, _ time.Duration) { retries = append(retries, n) },
			OnGiveUp:  func(n int, _ error, _ error) { giveUps = append(giveUps, n) },
		},
	})

	opErr := errors.New("x")
	err := r.Do(context.Background(), func(_ context.Context) error { return opErr })
	require.Error(t, err)

	assert.Equal(t, []int{0, 1}, attempts, "called once per attempt")
	assert.Equal(t, []int{0}, retries, "one retry after attempt 0")
	assert.Equal(t, []int{1}, giveUps, "give-up reported at last attempt index")
}

func TestStopper_IsConsultedBeforeClassifierAndOthers(t *testing.T) {
	// We ensure that even if the classifier would retry, the stopper halts first.
	fs := &fakeSleeper{now: time.Unix(0, 0)}
	zero := 0
	stopCause := errors.New("budget-exhausted")
	stopper := &fakeStopper{stopFromAttempt: &zero, cause: stopCause}

	r := retry.New(retry.Options{
		Strategy:   &fakeStrategy{delays: []time.Duration{1 * time.Second}},
		Jitter:     identJitter{},
		Classifier: alwaysRetryClassifier{},
		Sleeper:    fs,
		Stopper:    stopper,
	})

	errBoom := errors.New("boom")
	err := r.Do(context.Background(), func(_ context.Context) error { return errBoom })
	require.Error(t, err)
	require.ErrorIs(t, retry.Cause(err), stopCause)
	require.ErrorIs(t, retry.LastError(err), errBoom)
	assert.Equal(t, 1, retry.Attempts(err))
	assert.Equal(t, []int{0}, stopper.calls)
}

/* -------------------------
   Helpers coverage tests
------------------------- */

func TestFinalErrorHelpers_WorkOnNonFinalErrors(t *testing.T) {
	err := errors.New("plain")
	assert.NoError(t, retry.Cause(err))
	assert.NoError(t, retry.LastError(err))
	assert.Equal(t, 0, retry.Attempts(err))
}
