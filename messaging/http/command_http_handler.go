package messaginghttp

import (
	"context"
	"fmt"
	"net/http"

	"github.com/xfrr/go-cqrsify/messaging"
	"github.com/xfrr/go-cqrsify/pkg/apix"
)

// CommandHandler is an alias to HTTPMessageServer to keep external API surface familiar.
type CommandHandler = MessageHandler

// NewCommandHandler creates a new CommandHTTPServer with the given CommandBus and options.
// If no decoders are registered, the server will return 500 Internal Server Error.
func NewCommandHandler(dispatcher messaging.CommandDispatcher, opts ...HTTPMessageServerOption) *CommandHandler {
	return NewMessageHTTPHandler(&cmdbusWrapper{dispatcher}, opts...)
}

var _ messaging.MessagePublisher = (*cmdbusWrapper)(nil)

type cmdbusWrapper struct {
	dispatcher messaging.CommandDispatcher
}

func (w *cmdbusWrapper) Publish(ctx context.Context, msgs ...messaging.Message) error {
	commands := make([]messaging.Command, len(msgs))
	for i, m := range msgs {
		cmd, _ := m.(messaging.Command)
		commands[i] = cmd
	}
	return w.dispatcher.Dispatch(ctx, commands...)
}

// RegisterJSONAPICommandDecoder registers a JSON:API command decoder for the given command type.
// If a decoder for the same command type and encoding already exists, an error is returned.
func RegisterJSONAPICommandDecoder[A any](handler *MessageHandler, msgType string, decodeFunc func(apix.SingleDocument[A]) (messaging.Command, error)) error {
	if handler.decoders == nil {
		handler.decoders = make(map[string]map[HTTPMessageEncoding]func(*http.Request) (messaging.Message, error))
	}

	msgDecoders, ok := handler.decoders[msgType]
	if !ok {
		msgDecoders = make(map[HTTPMessageEncoding]func(*http.Request) (messaging.Message, error))
		handler.decoders[msgType] = msgDecoders
	}
	if _, exists := msgDecoders[HTTPMessageEncodingJSONAPI]; exists {
		return fmt.Errorf("command decoder for %q and encoding %q already exists", msgType, HTTPMessageEncodingJSONAPI)
	}

	msgDecoders[HTTPMessageEncodingJSONAPI] = makeMessageDecoder[A](func(sd apix.SingleDocument[A]) (messaging.Message, error) {
		return decodeFunc(sd)
	})
	return nil
}
