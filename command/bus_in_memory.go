package command

import (
	"context"
	"errors"
	"fmt"
	"sync"
)

var _ Bus = (*InMemoryBus)(nil)

var (
	// ErrNoSubscribers is returned when a command is dispatched but no subscribers are registered.
	ErrNoSubscribers = errors.New("no subscribers")
)

type InMemoryBusOption func(*InMemoryBus)

// WithBufferSize sets the buffer size for the in-memory bus.
// The default buffer size is 10 messages per subscriber.
func WithBufferSize(size uint) InMemoryBusOption {
	return func(b *InMemoryBus) {
		b.bufferSize = size
	}
}

type InMemoryBus struct {
	mu            sync.RWMutex
	bufferSize    uint
	subscriptions map[string][]chan anyContext
}

func NewInMemoryBus(opts ...InMemoryBusOption) (*InMemoryBus, error) {
	b := &InMemoryBus{
		bufferSize: 100,
	}
	for _, opt := range opts {
		opt(b)
	}
	return b, nil
}

func (b *InMemoryBus) Dispatch(ctx context.Context, topic string, cmd Command[any]) error {
	b.mu.RLock()
	defer b.mu.RUnlock()

	if b.subscriptions == nil || len(b.subscriptions[topic]) == 0 {
		return fmt.Errorf("%w: %s", ErrNoSubscribers, topic)
	}

	subscriptions := b.subscriptions[topic]

	for _, sub := range subscriptions {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			sub <- NewContext(ctx, cmd)
		}
	}

	return nil
}

func (b *InMemoryBus) Subscribe(ctx context.Context, commandName string) (<-chan anyContext, error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	ch := make(chan anyContext, b.bufferSize)

	if b.subscriptions == nil {
		b.subscriptions = make(map[string][]chan anyContext)
	}

	b.subscriptions[commandName] = append(b.subscriptions[commandName], ch)

	go func() {
		<-ctx.Done()
		b.mu.Lock()
		defer b.mu.Unlock()

		for i, h := range b.subscriptions[commandName] {
			if h == ch {
				b.subscriptions[commandName] = append(b.subscriptions[commandName][:i], b.subscriptions[commandName][i+1:]...)
				break
			}
		}
	}()

	return ch, nil
}
