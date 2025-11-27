package messaging

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/xfrr/go-cqrsify/pkg/multierror"
)

var (
	_ MessageHandler[Message]                        = (*MessageHandlerTypedRouter)(nil)
	_ MessageHandlerWithReply[Message, MessageReply] = (*MessageHandlerWithReplyTypedRouter)(nil)
)

// MessageHandlerTypedRouter routes messages to handlers based on message type.
type MessageHandlerTypedRouter struct {
	mu     sync.RWMutex
	byType map[string][]MessageHandler[Message]
}

// NewMessageHandlerTypedRouter creates a new MessageHandlerTypedRouter.
func NewMessageHandlerTypedRouter() *MessageHandlerTypedRouter {
	return &MessageHandlerTypedRouter{byType: make(map[string][]MessageHandler[Message])}
}

// Register adds a MessageHandler for the given message type.
func (r *MessageHandlerTypedRouter) Register(messageType string, h MessageHandler[Message]) {
	r.mu.Lock()
	r.byType[messageType] = append(r.byType[messageType], h)
	r.mu.Unlock()
}

// Handle routes the message to the appropriate handlers based on message type.
// It implements the MessageHandler interface.
func (r *MessageHandlerTypedRouter) Handle(ctx context.Context, msg Message) error {
	r.mu.RLock()
	handlers := append([]MessageHandler[Message](nil), r.byType[msg.MessageType()]...)
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
type MessageHandlerWithReplyTypedRouter struct {
	mu     sync.RWMutex
	byType map[string][]MessageHandlerWithReply[Message, MessageReply]
}

// NewMessageHandlerWithReplyTypedRouter creates a new MessageHandlerWithReplyTypedRouter.
func NewMessageHandlerWithReplyTypedRouter() *MessageHandlerWithReplyTypedRouter {
	return &MessageHandlerWithReplyTypedRouter{byType: make(map[string][]MessageHandlerWithReply[Message, MessageReply])}
}

// Register adds a MessageHandlerWithReply for the given message type.
//
// It only supports one handler per message type.
func (r *MessageHandlerWithReplyTypedRouter) Register(messageType string, h MessageHandlerWithReply[Message, MessageReply]) error {
	r.mu.Lock()
	if _, exists := r.byType[messageType]; exists {
		r.mu.Unlock()
		return errors.New("handler already exists for message type: " + messageType)
	}
	r.byType[messageType] = []MessageHandlerWithReply[Message, MessageReply]{h}
	r.mu.Unlock()
	return nil
}

// Handle routes the message to the appropriate handler based on message type.
// It implements the MessageHandlerWithReply interface.
func (r *MessageHandlerWithReplyTypedRouter) Handle(ctx context.Context, msg Message) (MessageReply, error) {
	r.mu.RLock()
	handlers := append([]MessageHandlerWithReply[Message, MessageReply](nil), r.byType[msg.MessageType()]...)
	r.mu.RUnlock()

	if len(handlers) == 0 {
		return nil, fmt.Errorf("%w for message type: %s", ErrHandlerNotFound, msg.MessageType())
	}

	return handlers[0].Handle(ctx, msg)
}

// RegisterCommandHandlerTypedRouter  is a helper function to register a CommandHandler in a MessageHandlerTypedRouter.
func RegisterCommandHandlerTypedRouter[T Command](router *MessageHandlerTypedRouter, commandType string, handler MessageHandler[T]) {
	router.Register(commandType, CommandHandlerFn(handler.Handle))
}

// RegisterEventHandlerTypedRouter is a helper function to register an EventHandler in a MessageHandlerTypedRouter.
func RegisterEventHandlerTypedRouter[T Event](router *MessageHandlerTypedRouter, eventType string, handler MessageHandler[T]) {
	router.Register(eventType, EventHandlerFn(handler.Handle))
}

// RegisterQueryHandlerTypedRouter is a helper function to register a QueryHandler in a MessageHandlerWithReplyTypedRouter.
func RegisterQueryHandlerTypedRouter[E Query, R QueryReply](router *MessageHandlerWithReplyTypedRouter, queryType string, handler MessageHandlerWithReply[E, R]) error {
	return router.Register(queryType, QueryHandlerFn(handler.Handle))
}
