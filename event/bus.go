package event

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

// DefaultBusBufferSize is the default number of events
// that can be queued for each subscriber.
const DefaultBusBufferSize = 1

var (
	// ErrNoSubscribers is returned when a event is published but no subscribers are registered.
	ErrNoSubscribers = errors.New("no subscribers")
)

// Bus represents a event bus that can publish events and subscribe to them.
// The behaviour of the bus is defined by the implementation of the Publisher and Subscriber interfaces.
type Bus interface {
	Publisher
	Subscriber
}

// A Publisher publishes events to subscribed handlers based on the topic.
type Publisher interface {
	// Publish publishes the provided event to the subscribers.
	// The behavior of this method depends on the implementation.
	Publish(ctx context.Context, reason string, evt Event[any]) error
}

// A Subscriber subscribes to events with a given topic.
type Subscriber interface {
	// Subscribe subscribes to the event with the provided topic.
	// The returned channels are closed when the context is canceled.
	// The behavior of this method depends on the implementation.
	Subscribe(ctx context.Context, topic string) (<-chan Context[any], error)
}

var _ Bus = (*bus)(nil)

// bus is a simple in-memory implementation of the Bus interface.
type bus struct {
	mu            sync.RWMutex
	bufferSize    uint
	subscriptions map[string][]chan anyContext

	// publishTimeout is the timeout for the publish operation.
	publishTimeout         time.Duration
	publishTimeoutFallback func(context.Context, string, Event[any])
}

// Publish publishes the provided event to the subscribers.
// If the context is canceled, the method returns an error.
// If no subscribers are registered for the provided event reason, the method returns an error.
// The method blocks until all events are published.
func (b *bus) Publish(ctx context.Context, reason string, evt Event[any]) error {
	b.mu.RLock()
	defer b.mu.RUnlock()

	if !b.hasSubscribers(reason) {
		return fmt.Errorf("%w: %s", ErrNoSubscribers, reason)
	}

	for _, sub := range b.subscriptions[reason] {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			evtctx := WithContext(ctx, evt)

			if b.publishTimeout > 0 {
				err := b.publishWithTimeout(evtctx, sub, reason, evt)
				if err != nil {
					b.timeoutFallback(ctx, reason, evt)
				}
			} else {
				sub <- evtctx
			}
		}
	}

	return nil
}

// Subscribe subscribes to the event with the provided topic.
// The returned channels are closed when the context is canceled.
// The method returns an error if the context is canceled.
func (b *bus) Subscribe(ctx context.Context, eventName string) (<-chan anyContext, error) {
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

func (b *bus) publishWithTimeout(evtctx anyContext, sub chan Context[any], topic string, evt Event[any]) error {
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
func (b *bus) timeoutFallback(ctx context.Context, topic string, evt Event[any]) {
	if b.publishTimeoutFallback != nil {
		b.publishTimeoutFallback(ctx, topic, evt)
	}
}

func (b *bus) newSubscription() chan anyContext {
	if b.subscriptions == nil {
		b.subscriptions = make(map[string][]chan anyContext)
	}
	return make(chan anyContext, b.bufferSize)
}

func (b *bus) hasSubscribers(reason string) bool {
	return b.subscriptions != nil && len(b.subscriptions[reason]) > 0
}

func (b *bus) addSubscription(ctx context.Context, reason string, ch chan anyContext) {
	b.subscriptions[reason] = append(b.subscriptions[reason], ch)
}

func (b *bus) removeSubscription(reason string, ch chan anyContext) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if subs, ok := b.subscriptions[reason]; ok {
		for i, sub := range subs {
			if sub == ch {
				close(sub)
				b.subscriptions[reason] = append(subs[:i], subs[i+1:]...)
				return
			}
		}
	}
}

// NewBus returns a new bus with the provided options.
// The default buffer size is 10 messages per subscriber.
func NewBus(opts ...BusOption) (*bus, error) {
	b := &bus{
		bufferSize: 10,
	}
	for _, opt := range opts {
		opt(b)
	}
	return b, nil
}
