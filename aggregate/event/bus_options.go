package event

type BusOption func(*InMemoryBus)

// WithBufferSize sets the buffer size for the bus.
func WithBufferSize(size uint) BusOption {
	return func(b *InMemoryBus) {
		b.bufferSize = size
	}
}
