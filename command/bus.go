package command

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

// DefaultBusBufferSize is the default number of commands
// that can be queued for each subscriber.
const DefaultBusBufferSize = 1

var (
	// ErrNoSubscribers is returned when a command is dispatched but no subscribers are registered.
	ErrNoSubscribers = errors.New("no subscribers")
)

// Bus represents a command bus that can dispatch commands and subscribe to them.
// The behaviour of the bus is defined by the implementation of the Dispatcher and Subscriber interfaces.
type Bus interface {
	Dispatcher
	Subscriber
}

// A Dispatcher dispatches commands to subscribed handlers based on the subject.
type Dispatcher interface {
	// Dispatch dispatches the provided command to the subscribers.
	// The behavior of this method depends on the implementation.
	Dispatch(ctx context.Context, subject string, cmd Command[any]) error
}

// A Subscriber subscribes to commands with a given subject.
type Subscriber interface {
	// Subscribe subscribes to the command with the provided subject.
	// The returned channels are closed when the context is canceled.
	// The behavior of this method depends on the implementation.
	Subscribe(ctx context.Context, subject string) (<-chan Context[any], error)
}

var _ Bus = (*bus)(nil)

// bus is a simple in-memory implementation of the Bus interface.
type bus struct {
	mu            sync.RWMutex
	bufferSize    uint
	subscriptions map[string][]chan anyContext

	// dispatchTimeout is the timeout for the dispatch operation.
	dispatchTimeout         time.Duration
	dispatchTimeoutFallback func(context.Context, string, Command[any])
}

// Dispatch dispatches the provided command to the subscribers.
// If the context is canceled, the method returns an error.
// If no subscribers are registered for the provided subject, the method returns an error.
// The method blocks until all commands are dispatched.
func (b *bus) Dispatch(ctx context.Context, subject string, cmd Command[any]) error {
	b.mu.RLock()
	defer b.mu.RUnlock()

	if !b.hasSubscribers(subject) {
		return fmt.Errorf("%w: %s", ErrNoSubscribers, subject)
	}

	for _, sub := range b.subscriptions[subject] {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			cmdctx := WithContext(ctx, cmd)

			if b.dispatchTimeout > 0 {
				err := b.dispatchWithTimeout(cmdctx, sub, subject, cmd)
				if err != nil {
					b.timeoutFallback(ctx, subject, cmd)
				}
			} else {
				sub <- cmdctx
			}
		}
	}

	return nil
}

// Subscribe subscribes to the command with the provided subject.
// The returned channels are closed when the context is canceled.
// The method returns an error if the context is canceled.
func (b *bus) Subscribe(ctx context.Context, commandName string) (<-chan anyContext, error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	ch := b.newSubscription()
	b.addSubscription(ctx, commandName, ch)

	go func() {
		<-ctx.Done()
		b.removeSubscription(commandName, ch)
	}()

	return ch, nil
}

func (b *bus) dispatchWithTimeout(cmdctx anyContext, sub chan Context[any], subject string, cmd Command[any]) error {
	dispatchCtx, cancel := context.WithTimeout(cmdctx, b.dispatchTimeout)
	defer cancel()

	select {
	case sub <- cmdctx:
	case <-dispatchCtx.Done():
		return dispatchCtx.Err()
	}

	return nil
}

// timeoutFallback calls the fallback function if the dispatch times out.
func (b *bus) timeoutFallback(ctx context.Context, subject string, cmd Command[any]) {
	if b.dispatchTimeoutFallback != nil {
		b.dispatchTimeoutFallback(ctx, subject, cmd)
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
// The default buffer size is 10 payloads per subscriber.
func NewBus(opts ...BusOption) (*bus, error) {
	b := &bus{
		bufferSize: 10,
	}
	for _, opt := range opts {
		opt(b)
	}
	return b, nil
}
