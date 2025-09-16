package messaging

import (
	"context"
	"time"
)

// MessageHandlerMiddleware is a middleware for message handlers.
type MessageHandlerMiddleware func(next MessageHandler[Message]) MessageHandler[Message]

// RecoverMiddleware shields the bus from panics in handlers.
func RecoverMiddleware(cb func(r any)) MessageHandlerMiddleware {
	return func(next MessageHandler[Message]) MessageHandler[Message] {
		return MessageHandlerFn[Message](func(ctx context.Context, msg Message) error {
			defer func() {
				if r := recover(); r != nil {
					cb(r)
				}
			}()
			return next.Handle(ctx, msg)
		})
	}
}

// TimeoutMiddleware enforces a per-message timeout; callers may pass ctx with deadline as well.
func TimeoutMiddleware(d time.Duration) MessageHandlerMiddleware {
	return func(next MessageHandler[Message]) MessageHandler[Message] {
		return MessageHandlerFn[Message](func(ctx context.Context, msg Message) error {
			if _, has := ctx.Deadline(); has || d <= 0 {
				return next.Handle(ctx, msg)
			}
			ctx, cancel := context.WithTimeout(ctx, d)
			defer cancel()
			return next.Handle(ctx, msg)
		})
	}
}

// RetryBackoffMiddleware retries the handler with exponential backoff on error.
func RetryBackoffMiddleware(attempts int, initialDelay time.Duration) MessageHandlerMiddleware {
	return func(next MessageHandler[Message]) MessageHandler[Message] {
		return MessageHandlerFn[Message](func(ctx context.Context, msg Message) error {
			var err error
			delay := initialDelay
			for range attempts {
				err = next.Handle(ctx, msg)
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
