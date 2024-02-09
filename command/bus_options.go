package command

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

// WithDispatchTimeout sets the timeout for the dispatch operation.
func WithDispatchTimeout(timeout time.Duration) BusOption {
	return func(b *bus) {
		b.dispatchTimeout = timeout
	}
}

// WithDispatchTimeoutFallback sets the function that is called
// when dispatching a command to a subscriber times out.
func WithDispatchTimeoutFallback(fb func(context.Context, string, Command[any])) BusOption {
	return func(b *bus) {
		b.dispatchTimeoutFallback = fb
	}
}
