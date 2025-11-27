package messaging

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/xfrr/go-cqrsify/pkg/multierror"
)

var (
	_ MessageHandler[Message]                        = (*MessageHandlerTypedRouter[Message])(nil)
	_ MessageHandlerWithReply[Message, MessageReply] = (*MessageHandlerWithReplyTypedRouter[Message, MessageReply])(nil)
)

// MessageHandlerTypedRouter routes messages to handlers based on message type.
type MessageHandlerTypedRouter[T Message] struct {
	mu     sync.RWMutex
	byType map[string][]MessageHandler[T]
}

// NewMessageHandlerTypedRouter creates a new MessageHandlerTypedRouter.
func NewMessageHandlerTypedRouter[T Message]() *MessageHandlerTypedRouter[T] {
	return &MessageHandlerTypedRouter[T]{byType: make(map[string][]MessageHandler[T])}
}

// Register adds a MessageHandler for the given message type.
func (r *MessageHandlerTypedRouter[T]) Register(messageType string, h MessageHandler[T]) {
	r.mu.Lock()
	r.byType[messageType] = append(r.byType[messageType], h)
	r.mu.Unlock()
}

// Handle routes the message to the appropriate handlers based on message type.
// It implements the MessageHandler interface.
func (r *MessageHandlerTypedRouter[T]) Handle(ctx context.Context, msg T) error {
	r.mu.RLock()
	handlers := append([]MessageHandler[T](nil), r.byType[msg.MessageType()]...)
	r.mu.RUnlock()

	if len(handlers) == 0 {
		return nil
	}

	errs := multierror.New()
	for i := range handlers {
		if err := handlers[i].Handle(ctx, msg); err != nil {
			errs.Append(err)
		}
	}
	return errs.ErrorOrNil()
}

// MessageHandlerWithReplyTypedRouter routes messages to handlers with reply based on message type.
type MessageHandlerWithReplyTypedRouter[T Message, R MessageReply] struct {
	mu     sync.RWMutex
	byType map[string][]MessageHandlerWithReply[T, R]
}

// NewMessageHandlerWithReplyTypedRouter creates a new MessageHandlerWithReplyTypedRouter.
func NewMessageHandlerWithReplyTypedRouter[T Message, R MessageReply]() *MessageHandlerWithReplyTypedRouter[T, R] {
	return &MessageHandlerWithReplyTypedRouter[T, R]{byType: make(map[string][]MessageHandlerWithReply[T, R])}
}

// Register adds a MessageHandlerWithReply for the given message type.
//
// It only supports one handler per message type.
func (r *MessageHandlerWithReplyTypedRouter[T, R]) Register(messageType string, h MessageHandlerWithReply[T, R]) error {
	r.mu.Lock()
	if _, exists := r.byType[messageType]; exists {
		r.mu.Unlock()
		return errors.New("handler already exists for message type: " + messageType)
	}
	r.byType[messageType] = []MessageHandlerWithReply[T, R]{h}
	r.mu.Unlock()
	return nil
}

// Handle routes the message to the appropriate handler based on message type.
// It implements the MessageHandlerWithReply interface.
func (r *MessageHandlerWithReplyTypedRouter[T, R]) Handle(ctx context.Context, msg T) (R, error) {
	r.mu.RLock()
	handlers := append([]MessageHandlerWithReply[T, R](nil), r.byType[msg.MessageType()]...)
	r.mu.RUnlock()

	if len(handlers) == 0 {
		var zero R
		return zero, fmt.Errorf("%w for message type: %s", ErrHandlerNotFound, msg.MessageType())
	}

	return handlers[0].Handle(ctx, msg)
}

// RegisterCommandHandlerTypedRouter  is a helper function to register a CommandHandler in a MessageHandlerTypedRouter.
func RegisterCommandHandlerTypedRouter[T Command](router *MessageHandlerTypedRouter[T], commandType string, handler MessageHandler[T]) {
	router.Register(commandType, handler)
}

// RegisterEventHandlerTypedRouter is a helper function to register an EventHandler in a MessageHandlerTypedRouter.
func RegisterEventHandlerTypedRouter[T Event](router *MessageHandlerTypedRouter[T], eventType string, handler MessageHandler[T]) {
	router.Register(eventType, handler)
}

// RegisterQueryHandlerTypedRouter is a helper function to register a QueryHandler in a MessageHandlerWithReplyTypedRouter.
func RegisterQueryHandlerTypedRouter[T Query, R QueryReply](router *MessageHandlerWithReplyTypedRouter[T, R], queryType string, handler MessageHandlerWithReply[T, R]) error {
	return router.Register(queryType, handler)
}

// NewCommandHandlerTypedRouter creates a new MessageHandlerTypedRouter for Commands.
func NewCommandHandlerTypedRouter() *MessageHandlerTypedRouter[Command] {
	return NewMessageHandlerTypedRouter[Command]()
}

// NewEventHandlerTypedRouter creates a new MessageHandlerTypedRouter for Events.
func NewEventHandlerTypedRouter() *MessageHandlerTypedRouter[Event] {
	return NewMessageHandlerTypedRouter[Event]()
}

// NewQueryHandlerTypedRouter creates a new MessageHandlerWithReplyTypedRouter for Queries.
func NewQueryHandlerTypedRouter() *MessageHandlerWithReplyTypedRouter[Query, QueryReply] {
	return NewMessageHandlerWithReplyTypedRouter[Query, QueryReply]()
}
