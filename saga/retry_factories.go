package saga

import (
	"time"

	"github.com/xfrr/go-cqrsify/pkg/retry"
)

func stepActionRetryFactory() RetryFactory {
	return func(step Step) *retry.Retrier {
		opts := retry.DefaultOptions()
		opts.Strategy = retry.ExponentialStrategy{
			Base:   time.Second,
			Factor: 2.0,
			Cap:    30 * time.Second,
		}
		opts.Jitter = retry.FullJitter{}

		so := step.RetryOptions
		if so.Strategy != nil {
			opts.Strategy = so.Strategy
		}
		if so.Jitter != nil {
			opts.Jitter = so.Jitter
		}
		if so.Classifier != nil {
			opts.Classifier = so.Classifier
		}
		if so.Stopper != nil {
			opts.Stopper = so.Stopper
		}
		if so.MaxAttempts != 0 {
			opts.MaxAttempts = so.MaxAttempts
		}
		if so.MaxElapsed > 0 {
			opts.MaxElapsed = so.MaxElapsed
		}
		if so.Sleeper != nil {
			opts.Sleeper = so.Sleeper
		}
		if so.Hooks.OnAttempt != nil {
			opts.Hooks.OnAttempt = so.Hooks.OnAttempt
		}
		if so.Hooks.OnRetry != nil {
			opts.Hooks.OnRetry = so.Hooks.OnRetry
		}
		if so.Hooks.OnGiveUp != nil {
			opts.Hooks.OnGiveUp = so.Hooks.OnGiveUp
		}
		return retry.New(opts)
	}
}

func stepCompensationRetryFactory(step Step) *retry.Retrier {
	opts := retry.DefaultOptions()
	opts.Strategy = retry.ExponentialStrategy{
		Base:   500 * time.Millisecond,
		Factor: 2.0,
		Cap:    10 * time.Second,
	}
	opts.Jitter = retry.FullJitter{}
	opts.MaxAttempts = 5

	co := step.CompensationRetryOptions
	if co.Strategy != nil {
		opts.Strategy = co.Strategy
	}
	if co.Jitter != nil {
		opts.Jitter = co.Jitter
	}
	if co.Classifier != nil {
		opts.Classifier = co.Classifier
	}
	if co.Stopper != nil {
		opts.Stopper = co.Stopper
	}
	if co.MaxAttempts != 0 {
		opts.MaxAttempts = co.MaxAttempts
	}
	if co.MaxElapsed > 0 {
		opts.MaxElapsed = co.MaxElapsed
	}
	if co.Sleeper != nil {
		opts.Sleeper = co.Sleeper
	}
	if co.Hooks.OnAttempt != nil {
		opts.Hooks.OnAttempt = co.Hooks.OnAttempt
	}
	if co.Hooks.OnRetry != nil {
		opts.Hooks.OnRetry = co.Hooks.OnRetry
	}
	if co.Hooks.OnGiveUp != nil {
		opts.Hooks.OnGiveUp = co.Hooks.OnGiveUp
	}
	return retry.New(opts)
}
