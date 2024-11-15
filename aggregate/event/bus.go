package event

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

const (
	defaultPublishTimeout = 5 * time.Second
	defaultBufferSize     = 10
)

var (
	// ErrBusHasNoSubscribers is returned when a event is published but no subscribers are registered.
	ErrBusHasNoSubscribers = errors.New("no subscribers")
)

// Bus represents a event bus that can publish events and subscribe to them.
// The behaviour of the bus is defined by the implementation of the Publisher and Subscriber interfaces.
type Bus interface {
	Publisher
	Subscriber
}

// A Publisher publishes events to subscribed handlers based on the subject.
type Publisher interface {
	// Publish publishes the provided event to the subscribers.
	// The behavior of this method depends on the implementation.
	Publish(ctx context.Context, evt Event[any, any]) error
}

// A Subscriber subscribes to events with a given subject.
type Subscriber interface {
	// Subscribe subscribes to the event with the provided subject.
	// The returned channels are closed when the context is canceled.
	// The behavior of this method depends on the implementation.
	Subscribe(ctx context.Context, subject string) (<-chan Context[any, any], error)
}

var _ Bus = (*InMemoryBus)(nil)

// InMemoryBus is a simple in-memory implementation of the Bus interface.
type InMemoryBus struct {
	mu                     sync.RWMutex
	bufferSize             uint
	subscriptions          map[string][]chan Context[any, any]
	publishTimeout         time.Duration
	publishTimeoutFallback func(context.Context, string, Event[any, any])
}

// Publish publishes the provided event to the subscribers.
// If the context is canceled, the method returns an error.
// If no subscribers are registered for the provided event name, the method returns an error.
// The method blocks until all events are published.
func (b *InMemoryBus) Publish(ctx context.Context, evt Event[any, any]) error {
	b.mu.RLock()
	defer b.mu.RUnlock()

	if !b.hasSubscribers(evt.Name()) {
		return fmt.Errorf("%w: %s", ErrBusHasNoSubscribers, evt.Name())
	}

	for _, sub := range b.subscriptions[evt.Name()] {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			evtctx := WithContext(ctx, evt)

			if b.publishTimeout > 0 {
				err := b.publishWithTimeout(evtctx, sub)
				if err != nil {
					b.timeoutFallback(ctx, evt.Name(), evt)
				}
			} else {
				sub <- evtctx
			}
		}
	}

	return nil
}

// Subscribe subscribes to the event with the provided subject.
// The returned channels are closed when the context is canceled.
// The method returns an error if the context is canceled.
func (b *InMemoryBus) Subscribe(ctx context.Context, eventName string) (<-chan Context[any, any], error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	ch := b.newSubscription()
	b.addSubscription(ctx, eventName, ch)

	go func() {
		<-ctx.Done()
		b.removeSubscription(eventName, ch)
	}()

	return ch, nil
}

func (b *InMemoryBus) publishWithTimeout(evtctx Context[any, any], sub chan Context[any, any]) error {
	publishCtx, cancel := context.WithTimeout(evtctx, b.publishTimeout)
	defer cancel()

	select {
	case sub <- evtctx:
	case <-publishCtx.Done():
		return publishCtx.Err()
	}

	return nil
}

// timeoutFallback calls the fallback function if the publish times out.
func (b *InMemoryBus) timeoutFallback(ctx context.Context, subject string, evt Event[any, any]) {
	if b.publishTimeoutFallback != nil {
		b.publishTimeoutFallback(ctx, subject, evt)
	}
}

func (b *InMemoryBus) newSubscription() chan Context[any, any] {
	if b.subscriptions == nil {
		b.subscriptions = make(map[string][]chan Context[any, any])
	}
	return make(chan Context[any, any], b.bufferSize)
}

func (b *InMemoryBus) hasSubscribers(name string) bool {
	return b.subscriptions != nil && len(b.subscriptions[name]) > 0
}

func (b *InMemoryBus) addSubscription(_ context.Context, name string, ch chan Context[any, any]) {
	b.subscriptions[name] = append(b.subscriptions[name], ch)
}

func (b *InMemoryBus) removeSubscription(name string, ch chan Context[any, any]) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if subs, ok := b.subscriptions[name]; ok {
		for i, sub := range subs {
			if sub == ch {
				close(sub)
				b.subscriptions[name] = append(subs[:i], subs[i+1:]...)
				return
			}
		}
	}
}

// NewInMemoryBus returns a new bus with the provided options.
// The default buffer size is 10 events per subscriber.
func NewInMemoryBus(opts ...BusOption) (*InMemoryBus, error) {
	b := &InMemoryBus{
		mu:                     sync.RWMutex{},
		bufferSize:             defaultBufferSize,
		publishTimeout:         defaultPublishTimeout,
		subscriptions:          make(map[string][]chan Context[any, any]),
		publishTimeoutFallback: nil,
	}
	for _, opt := range opts {
		opt(b)
	}
	return b, nil
}
