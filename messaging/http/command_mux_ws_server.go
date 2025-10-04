package messaginghttp

// NewMUXCommandWebsocketServer creates a new CommandWebsocketServer with the given CommandBus and options.
func NewMUXCommandWebsocketServer(handler *CommandHandler) *MessageMUXWebsocketServer {
	return NewMUXMessageWebsocketServer(handler)
}
