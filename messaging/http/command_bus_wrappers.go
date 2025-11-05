package messaginghttp

import (
	"context"

	"github.com/xfrr/go-cqrsify/messaging"
)

var _ messaging.MessagePublisher = (*commandDispatcherWrapper)(nil)

type commandDispatcherWrapper struct {
	dispatcher messaging.CommandDispatcher
}

func (w *commandDispatcherWrapper) Publish(ctx context.Context, msgs ...messaging.Message) error {
	commands := make([]messaging.Command, len(msgs))
	for i, m := range msgs {
		cmd, _ := m.(messaging.Command)
		commands[i] = cmd
	}
	return w.dispatcher.Dispatch(ctx, commands...)
}

var _ messaging.MessageConsumer = (*commandConsumerWrapper)(nil)

type commandConsumerWrapper struct {
	consumer messaging.CommandConsumer
}

func (w *commandConsumerWrapper) Subscribe(ctx context.Context, h messaging.MessageHandler[messaging.Message]) (messaging.UnsubscribeFunc, error) {
	return w.consumer.Subscribe(ctx, messaging.CommandHandlerFn[messaging.Command](func(ctx context.Context, cmd messaging.Command) error {
		return h.Handle(ctx, cmd)
	}))
}
