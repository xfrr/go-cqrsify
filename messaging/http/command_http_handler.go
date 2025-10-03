package messaginghttp

import (
	"context"
	"fmt"

	"github.com/xfrr/go-cqrsify/messaging"
)

// CommandHTTPServer is an alias to HTTPMessageServer to keep external API surface familiar.
type CommandHTTPServer = HTTPMessageServer

// NewCommandHTTPHandler creates a new CommandHTTPServer with the given CommandBus and options.
// If no decoders are registered, the server will return 500 Internal Server Error.
func NewCommandHTTPHandler(messageBus messaging.CommandBus, opts ...HTTPMessageServerOption) *CommandHTTPServer {
	return NewMessageHTTPHandler(&cmdbusWrapper{messageBus}, opts...)
}

var _ messaging.MessageBus = (*cmdbusWrapper)(nil)

type cmdbusWrapper struct {
	cmdbus messaging.CommandBus
}

func (w *cmdbusWrapper) Subscribe(ctx context.Context, subject string, h messaging.MessageHandler[messaging.Message]) (messaging.UnsubscribeFunc, error) {
	return w.cmdbus.Subscribe(ctx, subject, messaging.CommandHandlerFn[messaging.Command](func(ctx context.Context, cmd messaging.Command) error {
		return h.Handle(ctx, cmd)
	}))
}

func (w *cmdbusWrapper) Publish(ctx context.Context, msgs ...messaging.Message) error {
	commands := make([]messaging.Command, len(msgs))
	for i, m := range msgs {
		cmd, ok := m.(messaging.Command)
		if !ok {
			return fmt.Errorf("message at index %d is not a Command", i)
		}
		commands[i] = cmd
	}
	return w.cmdbus.Dispatch(ctx, commands...)
}
