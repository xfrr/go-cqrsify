package messaginghttp

import (
	"context"

	"github.com/xfrr/go-cqrsify/messaging"
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
