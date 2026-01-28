package messaging

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
)

const (
	// defaultQueueSize is the default size of the async delivery queue.
	defaultQueueSize = 100
)

var _ MessageBus = (*InMemoryMessageBus)(nil)
var _ MessageBusReplier = (*InMemoryMessageBus)(nil)

// InMemoryMessageBus is a simple, fast, process-local message bus.
type InMemoryMessageBus struct {
	opts MessageBusConfig

	mu       sync.RWMutex
	handlers map[string][]handlerEntry // subject -> handlers
	nextID   uint64                    // atomic incremental id for handlers

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

type handlerEntry struct {
	id uint64
	h  MessageHandler[Message]
}

type queued struct {
	//nolint:containedctx // context is passed from publisher to worker
	ctx context.Context
	msg Message
	h   MessageHandler[Message]
}

type worker struct {
	id int
}

type replyEnvelope struct {
	MessageReply
	replyCh chan MessageReply
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
		handlers: make(map[string][]handlerEntry),
	}

	if cfg.AsyncWorkers > 0 {
		qSize := cfg.QueueSize
		if qSize < 1 {
			qSize = 1
		}
		b.queue = make(chan queued, qSize)
		for i := range cfg.AsyncWorkers {
			b.addWorker(i)
		}
	}
	return b
}

func (b *InMemoryMessageBus) Publish(ctx context.Context, msgs ...Message) error {
	if len(msgs) == 0 {
		return errors.New("no messages to publish")
	}

	// Snapshot handlers outside of lock to avoid running user code while locked.
	b.mu.RLock()
	if b.closed {
		b.mu.RUnlock()
		return ErrPublishOnClosedBus
	}
	snap := make(map[string][]MessageHandler[Message], len(b.handlers))
	for mt, entries := range b.handlers {
		cp := make([]MessageHandler[Message], len(entries))
		for i := range entries {
			cp[i] = entries[i].h
		}
		snap[mt] = cp
	}
	b.mu.RUnlock()

	for _, msg := range msgs {
		handlers := snap[msg.MessageType()]
		if len(handlers) == 0 {
			return NoHandlersForMessageError{MessageType: msg.MessageType()}
		}
		for _, h := range handlers {
			if b.queue == nil {
				if err := b.deliverSync(ctx, h, msg); err != nil {
					return err
				}
				continue
			}
			if err := b.enqueue(ctx, h, msg); err != nil {
				return err
			}
		}
	}
	return nil
}

func (b *InMemoryMessageBus) PublishRequest(ctx context.Context, msg Message) (Message, error) {
	b.mu.RLock()
	if b.closed {
		b.mu.RUnlock()
		return nil, ErrPublishOnClosedBus
	}
	entries := b.handlers[msg.MessageType()]
	if len(entries) == 0 {
		b.mu.RUnlock()
		return nil, NoHandlersForMessageError{MessageType: msg.MessageType()}
	}
	h := entries[0].h
	b.mu.RUnlock()

	if h == nil {
		return nil, NoHandlersForMessageError{MessageType: msg.MessageType()}
	}

	wrapper, ok := h.(*inMemoryMessageBusHandlerWithReplyWrapper)
	if !ok {
		return nil, &InvalidMessageTypeError{
			Expected: "inMemoryMessageBusHandlerWithReplyWrapper",
			Actual:   fmt.Sprintf("%T", h),
		}
	}

	replyCh := make(chan MessageReply, 1)
	envelopedMsg := replyEnvelope{
		MessageReply: msg,
		replyCh:      replyCh,
	}

	// Use the wrapped handler to ensure middleware is applied.
	wrappedHandler := b.wrap(wrapper)
	if err := wrappedHandler.Handle(ctx, envelopedMsg); err != nil {
		return nil, err
	}

	select {
	case reply := <-replyCh:
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
	refs := make([]subRef, 0, len(b.opts.Subjects))

	b.mu.Lock()
	for _, subject := range b.opts.Subjects {
		id := b.addHandlerLocked(subject, h)
		refs = append(refs, subRef{subject: subject, id: id})
	}
	b.mu.Unlock()

	return func() error {
		return b.unsubscribeByRefs(refs)
	}, nil
}

func (b *InMemoryMessageBus) SubscribeWithReply(_ context.Context, h MessageHandlerWithReply[Message, MessageReply]) (UnsubscribeFunc, error) {
	wrapped := wrapMessageHandlerWithReply(h)
	refs := make([]subRef, 0, len(b.opts.Subjects))

	b.mu.Lock()
	for _, subject := range b.opts.Subjects {
		if hs := b.handlers[subject]; len(hs) > 0 {
			b.mu.Unlock()
			return nil, errors.New("reply handler already exists for message type: " + subject)
		}
		id := b.addHandlerLocked(subject, wrapped)
		refs = append(refs, subRef{subject: subject, id: id})
	}
	b.mu.Unlock()

	return func() error {
		return b.unsubscribeByRefs(refs)
	}, nil
}

func (b *InMemoryMessageBus) Close() error {
	b.closeMu.Lock()
	defer b.closeMu.Unlock()

	b.mu.Lock()
	defer b.mu.Unlock()

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
	b.workers = append(b.workers, worker{id: id})

	b.wg.Add(1)
	go func() {
		defer b.wg.Done()
		for q := range b.queue {
			h := b.wrap(q.h)
			if err := h.Handle(q.ctx, q.msg); err != nil && b.opts.ErrorHandler != nil {
				b.opts.ErrorHandler(q.msg.MessageType(), err)
			}
		}
	}()
}

func (b *InMemoryMessageBus) wrap(h MessageHandler[Message]) MessageHandler[Message] {
	for i := len(b.mw) - 1; i >= 0; i-- {
		h = b.mw[i](h)
	}
	return h
}

// internal subscription reference
type subRef struct {
	subject string
	id      uint64
}

func (b *InMemoryMessageBus) addHandlerLocked(subject string, h MessageHandler[Message]) uint64 {
	id := atomic.AddUint64(&b.nextID, 1)
	b.handlers[subject] = append(b.handlers[subject], handlerEntry{id: id, h: h})
	return id
}

func (b *InMemoryMessageBus) unsubscribeByRefs(refs []subRef) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	for _, ref := range refs {
		hs, ok := b.handlers[ref.subject]
		if !ok || len(hs) == 0 {
			continue
		}
		for i := range hs {
			if hs[i].id == ref.id {
				last := len(hs) - 1
				if i != last {
					hs[i] = hs[last]
				}
				hs[last] = handlerEntry{}
				hs = hs[:last]
				if len(hs) == 0 {
					delete(b.handlers, ref.subject)
				} else {
					b.handlers[ref.subject] = hs
				}
				break
			}
		}
	}
	return nil
}

type inMemoryMessageBusHandlerWithReplyWrapper struct {
	h MessageHandlerWithReply[Message, MessageReply]
}

func wrapMessageHandlerWithReply(h MessageHandlerWithReply[Message, MessageReply]) MessageHandler[Message] {
	return &inMemoryMessageBusHandlerWithReplyWrapper{
		h: h,
	}
}

func (w *inMemoryMessageBusHandlerWithReplyWrapper) Handle(ctx context.Context, msg Message) error {
	env, ok := msg.(replyEnvelope)
	if !ok {
		return &InvalidMessageTypeError{
			Expected: "replyEnvelope",
			Actual:   fmt.Sprintf("%T", msg),
		}
	}

	reply, err := w.h.Handle(ctx, env.MessageReply)
	if err != nil {
		return err
	}

	select {
	case env.replyCh <- reply:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
