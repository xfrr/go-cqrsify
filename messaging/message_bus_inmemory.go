package messaging

import (
	"context"
	"sync"
)

const (
	// defaultQueueSize is the default size of the async delivery queue.
	// Increase if you expect bursts of messages or slow handlers.
	// Decrease to limit memory usage.
	defaultQueueSize = 100
)

var _ MessageBus = (*InMemoryMessageBus)(nil)

// InMemoryMessageBus is a simple, fast, process-local message bus.
// Great for tests and single-process apps; swap for a distributed bus in prod if needed.
type InMemoryMessageBus struct {
	opts        MessageBusConfig
	mu          sync.RWMutex
	subscribers map[string][]MessageHandler[Message]

	// async pipeline (enabled if opts.AsyncWorkers > 0)
	queue   chan queued
	workers []worker

	// composed middleware chain applied to handlers
	mw []MessageHandlerMiddleware

	// lifecycle
	closed  bool
	closeMu sync.Mutex
	wg      sync.WaitGroup
}

type queued struct {
	ctx context.Context //nolint:containedctx // required for async delivery
	msg Message
	h   MessageHandler[Message]
}

type worker struct {
	id int
}

func NewInMemoryMessageBus(optFns ...MessageBusConfigModifier) *InMemoryMessageBus {
	cfg := MessageBusConfig{
		AsyncWorkers: 0,
		QueueSize:    defaultQueueSize,
		ErrorHandler: nil,
	}
	for _, fn := range optFns {
		fn(&cfg)
	}

	b := &InMemoryMessageBus{
		opts:        cfg,
		subscribers: make(map[string][]MessageHandler[Message]),
	}

	if cfg.AsyncWorkers > 0 {
		b.queue = make(chan queued, max(1, cfg.QueueSize))
		for i := range cfg.AsyncWorkers {
			b.addWorker(i)
		}
	}
	return b
}

func (b *InMemoryMessageBus) Publish(ctx context.Context, msgs ...Message) error {
	b.mu.RLock()
	defer b.mu.RUnlock()

	if b.closed {
		return ErrPublishOnClosedBus
	}

	for _, msg := range msgs {
		handlers := b.subscribers[msg.MessageType()]
		if len(handlers) == 0 {
			return &NoSubscribersForMessageError{MessageType: msg.MessageType()}
		}

		// Note: we do not short-circuit on handler errors; we attempt to
		// dispatch to all handlers (sync or async)
		for _, h := range handlers {
			if b.queue == nil {
				// Synchronous inline dispatch
				if err := b.wrap(h).Handle(ctx, msg); err != nil {
					if b.opts.ErrorHandler != nil {
						b.opts.ErrorHandler(msg.MessageType(), err)
					} else {
						return err
					}
				}
				continue
			}
			// Async: enqueue delivery (non-blocking if buffered; otherwise may block).
			select {
			case b.queue <- queued{ctx: ctx, msg: msg, h: h}:
			case <-ctx.Done():
				return ctx.Err()
			}
		}
	}

	return nil
}

func (b *InMemoryMessageBus) Subscribe(_ context.Context, messageName string, h MessageHandler[Message]) (UnsubscribeFunc, error) {
	b.mu.Lock()
	b.subscribers[messageName] = append(b.subscribers[messageName], h)
	idx := len(b.subscribers[messageName]) - 1
	b.mu.Unlock()

	return func() {
		b.mu.Lock()
		defer b.mu.Unlock()
		hs := b.subscribers[messageName]
		if idx < 0 || idx >= len(hs) {
			// handler already removed
			return
		}

		// Remove handler by swapping with last and truncating slice.
		hs[idx] = hs[len(hs)-1]
		b.subscribers[messageName] = hs[:len(hs)-1]
	}, nil
}

func (b *InMemoryMessageBus) Close() error {
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

func (b *InMemoryMessageBus) Use(mw ...MessageHandlerMiddleware) {
	b.mw = append(b.mw, mw...)
}

func (b *InMemoryMessageBus) addWorker(id int) {
	// Each worker consumes queued deliveries; failures go to ErrorHandler.
	b.workers = append(b.workers, worker{id: id})
	b.wg.Go(func() {
		defer b.wg.Done()
		for q := range b.queue {
			// Compose middleware chain around handler for each delivery.
			h := b.wrap(q.h)
			if err := h.Handle(q.ctx, q.msg); err != nil {
				if b.opts.ErrorHandler != nil {
					b.opts.ErrorHandler(q.msg.MessageType(), err)
				}
			}
		}
	})
}

func (b *InMemoryMessageBus) wrap(h MessageHandler[Message]) MessageHandler[Message] {
	for i := len(b.mw) - 1; i >= 0; i-- {
		h = b.mw[i](h)
	}
	return h
}
