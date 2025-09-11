package domain

import (
	"context"
	"time"
)

// RecoverMiddleware shields the bus from panics in handlers.
func RecoverMiddleware(cb func(r any)) EventHandlerMiddleware {
	return func(next EventHandler) EventHandler {
		return EventHandlerFunc(func(ctx context.Context, evt Event) (err error) {
			defer func() {
				if r := recover(); r != nil {
					cb(r)
				}
			}()
			return next.Handle(ctx, evt)
		})
	}
}

// TimeoutMiddleware enforces a per-event timeout; callers may pass ctx with deadline as well.
func TimeoutMiddleware(d time.Duration) EventHandlerMiddleware {
	return func(next EventHandler) EventHandler {
		return EventHandlerFunc(func(ctx context.Context, evt Event) error {
			if _, has := ctx.Deadline(); has || d <= 0 {
				return next.Handle(ctx, evt)
			}
			ctx, cancel := context.WithTimeout(ctx, d)
			defer cancel()
			return next.Handle(ctx, evt)
		})
	}
}

// RetryBackoffMiddleware retries the handler with exponential backoff on error.
func RetryBackoffMiddleware(attempts int, initialDelay time.Duration) EventHandlerMiddleware {
	return func(next EventHandler) EventHandler {
		return EventHandlerFunc(func(ctx context.Context, evt Event) error {
			var err error
			delay := initialDelay
			for i := 0; i < attempts; i++ {
				err = next.Handle(ctx, evt)
				if err == nil {
					return nil
				}
				time.Sleep(delay)
				delay *= 2
			}
			return err
		})
	}
}
