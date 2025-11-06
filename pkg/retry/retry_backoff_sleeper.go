package retry

import (
	"context"
	"fmt"
	"time"
)

// Sleeper abstracts sleeping (time). Useful for tests or virtual clocks.
type Sleeper interface {
	Sleep(ctx context.Context, d time.Duration) error
	Now() time.Time
}

type RealSleeper struct{}

func (RealSleeper) Sleep(ctx context.Context, d time.Duration) error {
	t := time.NewTimer(d)
	defer t.Stop()
	select {
	case <-ctx.Done():
		return fmt.Errorf("sleep interrupted: %w", ctx.Err())
	case <-t.C:
		return nil
	}
}
func (RealSleeper) Now() time.Time { return time.Now() }

func defaultSleeper(sleeper Sleeper) Sleeper {
	if sleeper != nil {
		return sleeper
	}
	return RealSleeper{}
}
