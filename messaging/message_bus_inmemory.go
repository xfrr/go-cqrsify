package messaging

import (
	"context"
	"errors"
	"sync"
)

const (
	// defaultQueueSize is the default size of the async delivery queue.
	// Increase if you expect bursts of messages or slow handlers.
	// Decrease to limit memory usage.
	defaultQueueSize = 100
)

var _ MessageBus = (*InMemoryMessageBus)(nil)
var _ MessageBusReplier = (*InMemoryMessageBus)(nil)

// InMemoryMessageBus is a simple, fast, process-local message bus.
// Great for tests and single-process apps; swap for a distributed bus in prod if needed.
type InMemoryMessageBus struct {
	opts     MessageBusConfig
	mu       sync.RWMutex
	handlers map[string][]MessageHandler[Message]

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

func NewInMemoryMessageBus(optFns ...MessageBusConfigConfiger) *InMemoryMessageBus {
	cfg := MessageBusConfig{
		AsyncWorkers: 0,
		QueueSize:    defaultQueueSize,
		ErrorHandler: nil,
	}
	for _, fn := range optFns {
		fn(&cfg)
	}

	b := &InMemoryMessageBus{
		opts:     cfg,
		handlers: make(map[string][]MessageHandler[Message]),
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

	if len(msgs) == 0 {
		return errors.New("no messages to publish")
	}

	if b.closed {
		return ErrPublishOnClosedBus
	}

	for _, msg := range msgs {
		handlers := b.handlers[msg.MessageType()]
		if len(handlers) == 0 {
			return NoHandlersForMessageError{MessageType: msg.MessageType()}
		}

		// Deliver to all handlers (sync or async).
		for _, h := range handlers {
			if b.queue == nil {
				// Synchronous inline dispatch
				if err := b.deliverSync(ctx, h, msg); err != nil {
					return err
				}
				continue
			}

			// Async: enqueue delivery (non-blocking if buffered; otherwise may block).
			if err := b.enqueue(ctx, h, msg); err != nil {
				return err
			}
		}
	}

	return nil
}

func (b *InMemoryMessageBus) PublishRequest(ctx context.Context, msg Message) (Message, error) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	if b.closed {
		return nil, ErrPublishOnClosedBus
	}

	handlers := b.handlers[msg.MessageType()]
	if len(handlers) == 0 {
		return nil, NoHandlersForMessageError{MessageType: msg.MessageType()}
	}

	// For request/reply, we only support a single handler.
	h := handlers[0]
	if h == nil {
		return nil, NoHandlersForMessageError{MessageType: msg.MessageType()}
	}

	// Extract reply from handler wrapper.
	wrapper, ok := h.(*inMemoryMessageBusHandlerWithReplyWrapper)
	if !ok {
		return nil, &InvalidMessageTypeError{Expected: "inMemoryMessageBusHandlerWithReplyWrapper", Actual: ""}
	}

	// Compose middleware chain around handler for delivery.
	wrappedHandler := b.wrap(wrapper)

	// Deliver and get reply.
	if err := wrappedHandler.Handle(ctx, msg); err != nil {
		return nil, err
	}

	select {
	case reply := <-wrapper.out:
		return reply, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func (b *InMemoryMessageBus) deliverSync(ctx context.Context, h MessageHandler[Message], msg Message) error {
	if err := b.wrap(h).Handle(ctx, msg); err != nil {
		if b.opts.ErrorHandler != nil {
			b.opts.ErrorHandler(msg.MessageType(), err)
			return nil
		}
		return err
	}
	return nil
}

func (b *InMemoryMessageBus) enqueue(ctx context.Context, h MessageHandler[Message], msg Message) error {
	select {
	case b.queue <- queued{ctx: ctx, msg: msg, h: h}:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (b *InMemoryMessageBus) Subscribe(_ context.Context, h MessageHandler[Message]) (UnsubscribeFunc, error) {
	b.mu.Lock()
	// register handler for all subjects
	idxs := make([]int, 0, len(b.opts.Subjects))
	for _, messageName := range b.opts.Subjects {
		b.handlers[messageName] = append(b.handlers[messageName], h)
		idxs = append(idxs, len(b.handlers[messageName])-1)
	}
	b.mu.Unlock()

	return UnsubscribeFunc(func() error {
		return b.unsubscribeFunc(idxs)
	}), nil
}

func (b *InMemoryMessageBus) SubscribeWithReply(_ context.Context, h MessageHandlerWithReply[Message, MessageReply]) (UnsubscribeFunc, error) {
	b.mu.Lock()
	// register handler for all subjects
	idxs := make([]int, 0, len(b.opts.Subjects))
	for _, subject := range b.opts.Subjects {
		if _, exists := b.handlers[subject]; exists {
			// only one reply handler per message type
			b.mu.Unlock()
			return nil, errors.New("reply handler already exists for message type: " + subject)
		}

		b.handlers[subject] = append(b.handlers[subject], wrapMessageHandlerWithReply(h))
		idxs = append(idxs, len(b.handlers[subject])-1)
	}
	b.mu.Unlock()

	return UnsubscribeFunc(func() error {
		return b.unsubscribeFunc(idxs)
	}), nil
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

func (b *InMemoryMessageBus) unsubscribeFunc(idxs []int) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	for i, subject := range b.opts.Subjects {
		if i >= len(idxs) {
			continue
		}

		hs, ok := b.handlers[subject]
		if !ok || len(hs) == 0 {
			continue
		}

		idx := idxs[i]
		if idx < 0 || idx >= len(hs) {
			continue
		}

		last := len(hs) - 1
		if idx != last {
			hs[idx] = hs[last]
		}
		hs[last] = nil
		hs = hs[:last]

		if len(hs) == 0 {
			delete(b.handlers, subject)
		} else {
			b.handlers[subject] = hs
		}
	}
	return nil
}

type inMemoryMessageBusHandlerWithReplyWrapper struct {
	h   MessageHandlerWithReply[Message, MessageReply]
	out chan MessageReply
}

func wrapMessageHandlerWithReply(h MessageHandlerWithReply[Message, MessageReply]) MessageHandler[Message] {
	return &inMemoryMessageBusHandlerWithReplyWrapper{
		h:   h,
		out: make(chan MessageReply, 1),
	}
}

func (w *inMemoryMessageBusHandlerWithReplyWrapper) Handle(ctx context.Context, msg Message) error {
	reply, err := w.h.Handle(ctx, msg)
	if err != nil {
		return err
	}

	select {
	case w.out <- reply:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
