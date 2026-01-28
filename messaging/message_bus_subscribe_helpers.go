package messaging

import (
	"context"
	"fmt"
)

// RegisterMessageHandler is a shorthand for handling messages.
func RegisterMessageHandler[E Message](
	ctx context.Context,
	consumer MessageConsumer,
	handlerFn MessageHandler[E],
) (UnsubscribeFunc, error) {
	return consumer.Subscribe(
		ctx,
		MessageHandlerFn[Message](func(ctx context.Context, evt Message) error {
			castMessage, ok := evt.(E)
			if !ok {
				return InvalidMessageTypeError{
					Actual:   fmt.Sprintf("%T", evt),
					Expected: fmt.Sprintf("%T", *new(E)),
				}
			}
			return handlerFn.Handle(ctx, castMessage)
		}),
	)
}

// RegisterMessageHandlerWithReply is a shorthand for handling messages with reply.
func RegisterMessageHandlerWithReply[E Message, R MessageReply](
	ctx context.Context,
	consumer MessageConsumerReplier,
	handler MessageHandlerWithReply[E, R],
) (UnsubscribeFunc, error) {
	return consumer.SubscribeWithReply(
		ctx,
		MessageHandlerWithReplyFn[Message, MessageReply](func(ctx context.Context, msg Message) (MessageReply, error) {
			castMsg, ok := msg.(E)
			if !ok {
				return nil, InvalidMessageTypeError{
					Actual:   fmt.Sprintf("%T", msg),
					Expected: fmt.Sprintf("%T", *new(E)),
				}
			}
			return handler.Handle(ctx, castMsg)
		}),
	)
}
