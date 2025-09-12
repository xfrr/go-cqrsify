package retry

import (
	"context"
	"fmt"
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
				return true, fmt.Errorf("stopper %T signaled to stop: %w", stopper, cause)
			}
		}
		return false, nil
	})
}
