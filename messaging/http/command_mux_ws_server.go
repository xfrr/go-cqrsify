package messaginghttp

import (
	"github.com/xfrr/go-cqrsify/messaging"
)

// NewCommandWebsocketServer creates a new CommandWebsocketServer with the given CommandBus and options.
func NewCommandWebsocketServer(cmdbus messaging.CommandBus) *MessageWebsocketServer {
	handler := NewCommandHandler(cmdbus)
	return newMessageWebsocketServer(handler)
}
