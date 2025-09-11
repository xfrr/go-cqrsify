package domain

import (
	"context"
	"errors"
	"sync"
)

var _ EventBus = (*InMemoryEventBus)(nil)

// ErrNoSubscribersForEvent is returned by Publish if no subscribers exist for an event.
type ErrNoSubscribersForEvent struct {
	EventName string
}

func (e ErrNoSubscribersForEvent) Error() string {
	return "eventbus: no subscribers for event " + e.EventName
}

// InMemoryEventBus is a simple, fast, process-local event bus.
// Great for tests and single-process apps; swap for a distributed bus in prod if needed.
type InMemoryEventBus struct {
	opts        EventBusConfig
	mu          sync.RWMutex
	subscribers map[string][]EventHandler[Event]

	// async pipeline (enabled if opts.AsyncWorkers > 0)
	queue   chan queued
	workers []worker

	// composed middleware chain applied to handlers
	mw []EventHandlerMiddleware[Event]

	closed  bool
	closeMu sync.Mutex
	wg      sync.WaitGroup
}

type queued struct {
	ctx context.Context
	evt Event
	h   EventHandler[Event]
}

type worker struct {
	id int
}

func NewInMemoryEventBus(optFns ...EventBusConfigModifier) *InMemoryEventBus {
	cfg := EventBusConfig{
		AsyncWorkers: 0,
		QueueSize:    1024,
	}
	for _, fn := range optFns {
		fn(&cfg)
	}

	b := &InMemoryEventBus{
		opts:        cfg,
		subscribers: make(map[string][]EventHandler[Event]),
	}

	if cfg.AsyncWorkers > 0 {
		b.queue = make(chan queued, max(1, cfg.QueueSize))
		for i := 0; i < cfg.AsyncWorkers; i++ {
			b.addWorker(i)
		}
	}
	return b
}

func (b *InMemoryEventBus) addWorker(id int) {
	// Each worker consumes queued deliveries; failures go to ErrorHandler.
	b.workers = append(b.workers, worker{id: id})
	b.wg.Add(1)
	go func() {
		defer b.wg.Done()
		for q := range b.queue {
			// Compose middleware chain around handler for each delivery.
			h := b.wrap(q.h)
			if err := h.Handle(q.ctx, q.evt); err != nil {
				if b.opts.ErrorHandler != nil {
					b.opts.ErrorHandler(q.evt.Name(), err)
				}
			}
		}
	}()
}

func (b *InMemoryEventBus) Use(mw ...EventHandlerMiddleware[Event]) { b.mw = append(b.mw, mw...) }

func (b *InMemoryEventBus) wrap(h EventHandler[Event]) EventHandler[Event] {
	// Apply middlewares in registration order
	for i := len(b.mw) - 1; i >= 0; i-- {
		h = b.mw[i](h)
	}
	return h
}

func (b *InMemoryEventBus) Publish(ctx context.Context, evts ...Event) error {
	b.mu.RLock()
	defer b.mu.RUnlock()

	if b.closed {
		return errors.New("eventbus: publish on closed bus")
	}

	for _, evt := range evts {
		handlers := append([]EventHandler[Event](nil), b.subscribers[evt.Name()]...)
		if len(handlers) == 0 {
			return ErrNoSubscribersForEvent{EventName: evt.Name()}
		}

		// dispatch to all handlers (sync or async)
		for _, h := range handlers {
			if b.queue == nil {
				// Synchronous inline dispatch
				if err := b.wrap(h).Handle(ctx, evt); err != nil {
					if b.opts.ErrorHandler != nil {
						b.opts.ErrorHandler(evt.Name(), err)
					}
				}
				continue
			}
			// Async: enqueue delivery (non-blocking if buffered; otherwise may block).
			b.queue <- queued{ctx: ctx, evt: evt, h: h}
		}
	}
	return nil
}

func (b *InMemoryEventBus) Subscribe(eventName string, h EventHandler[Event]) (unsubscribe func()) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.subscribers[eventName] = append(b.subscribers[eventName], h)

	return func() {
		b.mu.Lock()
		defer b.mu.Unlock()
		hs := b.subscribers[eventName]
		for i := range hs {
			if &hs[i] == &h {
				b.subscribers[eventName] = append(hs[:i], hs[i+1:]...)
				break
			}
		}
	}
}

func (b *InMemoryEventBus) Close() error {
	b.closeMu.Lock()
	defer b.closeMu.Unlock()

	if b.closed {
		return nil
	}
	b.closed = true

	if b.queue != nil {
		close(b.queue)
	}
	b.wg.Wait()
	return nil
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
