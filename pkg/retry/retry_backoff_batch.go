package retry

import (
	"context"
	"runtime"
	"sync"

	"github.com/xfrr/go-cqrsify/pkg/multierror"
)

// DoNOptions tunes batch execution.
type DoNOptions struct {
	// Maximum parallel workers. If <= 0, defaults to runtime.NumCPU().
	Concurrency int
	// If true, cancel the whole batch on the first error.
	FailFast bool
	// If true, return per-item errors; otherwise collapse to first non-nil.
	CollectErrors bool
}

// DoN retries a batch of items using the same Retrier configuration and a
// *shared* time budget (MaxElapsed). Each item gets the **same deadline** so
// the whole batch respects one wall-clock budget.
// Returns either a slice of per-item errors (if CollectErrors) or the first error encountered.
func DoN[T any](
	ctx context.Context,
	r *Retrier,
	items []T,
	fn func(context.Context, T) error,
	opt DoNOptions,
) error {
	n := len(items)
	if n == 0 {
		return nil
	}
	if opt.Concurrency <= 0 {
		opt.Concurrency = runtime.NumCPU()
	}
	sem := make(chan struct{}, opt.Concurrency)
	var (
		wg     sync.WaitGroup
		errsMu sync.Mutex
		errs   = make([]error, n)
	)
	// Shared deadline (if configured)
	sharedCtx := ctx
	cancel := func() {}
	if r.maxElapsed > 0 {
		deadline := r.sleeper.Now().Add(r.maxElapsed)
		sharedCtx, cancel = context.WithDeadline(ctx, deadline)
	} else {
		sharedCtx, cancel = context.WithCancel(ctx)
	}
	defer cancel()

	var onceCancel sync.Once
	recordErr := func(idx int, err error) {
		if err == nil {
			return
		}
		errsMu.Lock()
		errs[idx] = err
		errsMu.Unlock()
		if opt.FailFast {
			onceCancel.Do(func() { cancel() })
		}
	}

	wg.Add(n)
	for i := range items {
		i := i
		sem <- struct{}{}
		go func() {
			defer func() { <-sem; wg.Done() }()
			// Each item uses the same Retrier and the sharedCtx (shared budget).
			e := r.Do(sharedCtx, func(ctx context.Context) error {
				return fn(ctx, items[i])
			})
			recordErr(i, e)
		}()
	}
	wg.Wait()

	// Collapse errors if requested
	if opt.CollectErrors {
		// Return a multi-error for convenience
		var first error
		errsMu.Lock()
		defer errsMu.Unlock()
		for _, e := range errs {
			if e != nil && first == nil {
				first = e
			}
		}
		if first == nil {
			return nil
		}
		multierr := multierror.New()
		multierr.Append(errs...)
		return multierr
	}
	// First-error mode
	errsMu.Lock()
	defer errsMu.Unlock()
	for _, e := range errs {
		if e != nil {
			return e
		}
	}
	return nil
}
