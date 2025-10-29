package messaginghttp

type MessageWebsocketServerOption func(*MessageWebsocketServer)

// WithWebsocketErrorHandler sets a custom error handler for websocket errors.
func WithWebsocketErrorHandler(h func(error)) MessageWebsocketServerOption {
	return func(s *MessageWebsocketServer) {
		s.errorHandler = h
	}
}
