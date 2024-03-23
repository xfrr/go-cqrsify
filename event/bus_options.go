package event

import (
	"context"
	"time"
)

type BusOption func(*bus)

// WithBufferSize sets the buffer size for the bus.
func WithBufferSize(size uint) BusOption {
	return func(b *bus) {
		b.bufferSize = size
	}
}

// WithPublishTimeout sets the timeout for the publish operation.
func WithPublishTimeout(timeout time.Duration) BusOption {
	return func(b *bus) {
		b.publishTimeout = timeout
	}
}

// WithPublishTimeoutFallback sets the function that is called
// when publishing a event to a subscriber times out.
func WithPublishTimeoutFallback(fb func(context.Context, string, Event[any, any])) BusOption {
	return func(b *bus) {
		b.publishTimeoutFallback = fb
	}
}
