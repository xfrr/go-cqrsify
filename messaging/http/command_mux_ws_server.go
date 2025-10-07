package messaginghttp

// NewCommandWebsocketServer creates a new CommandWebsocketServer with the given CommandBus and options.
func NewCommandWebsocketServer(handler *CommandHandler) *MessageMUXWebsocketServer {
	return NewMessageWebsocketServer(handler)
}
