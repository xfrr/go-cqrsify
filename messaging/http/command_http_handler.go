package messaginghttp

import (
	"github.com/xfrr/go-cqrsify/messaging"
	"github.com/xfrr/go-cqrsify/pkg/apix"
)

// CommandHandler is an alias to HTTPMessageServer to keep external API surface familiar.
type CommandHandler = MessageHandler

// NewCommandHandler creates a new CommandHTTPServer with the given CommandBus and options.
// If no decoders are registered, the server will return 500 Internal Server Error.
func NewCommandHandler(dispatcher messaging.CommandDispatcher, opts ...HTTPMessageServerOption) *CommandHandler {
	return NewMessageHandler(&commandDispatcherWrapper{dispatcher}, opts...)
}

// RegisterJSONAPICommandDecoder registers a JSON:API command decoder for the given command type.
// If a decoder for the same command type and encoding already exists, an error is returned.
func RegisterJSONAPICommandDecoder[A any](handler *MessageHandler, msgType string, decodeFunc func(apix.SingleDocument[A]) (messaging.Command, error)) error {
	return RegisterJSONAPIMessageDecoder(handler, msgType, func(sd apix.SingleDocument[A]) (messaging.Message, error) {
		return decodeFunc(sd)
	})
}
