package messaginghttp

type MessageWebsocketServerOption func(*MessageMUXWebsocketServer)

// WithWebsocketErrorHandler sets a custom error handler for websocket errors.
func WithWebsocketErrorHandler(h func(error)) MessageWebsocketServerOption {
	return func(s *MessageMUXWebsocketServer) {
		s.errorHandler = h
	}
}
