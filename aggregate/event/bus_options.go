package event

import (
	"context"
	"time"
)

type BusOption func(*InMemoryBus)

// WithBufferSize sets the buffer size for the bus.
func WithBufferSize(size uint) BusOption {
	return func(b *InMemoryBus) {
		b.bufferSize = size
	}
}

// WithPublishTimeout sets the timeout for the publish operation.
func WithPublishTimeout(timeout time.Duration) BusOption {
	return func(b *InMemoryBus) {
		b.publishTimeout = timeout
	}
}

// WithPublishTimeoutFallback sets the function that is called
// when publishing a event to a subscriber times out.
func WithPublishTimeoutFallback(fb func(context.Context, string, Event[any, any])) BusOption {
	return func(b *InMemoryBus) {
		b.publishTimeoutFallback = fb
	}
}
