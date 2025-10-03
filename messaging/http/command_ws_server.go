package messaginghttp

import "github.com/xfrr/go-cqrsify/messaging"

// NewCommandWebsocketServer creates a new CommandWebsocketServer with the given CommandBus and options.
func NewCommandWebsocketServer(dispatcher messaging.CommandDispatcher, handler *CommandHandler) *MessageWebsocketServer {
	return NewMessageWebsocketServer(&cmdbusWrapper{dispatcher}, handler)
}
