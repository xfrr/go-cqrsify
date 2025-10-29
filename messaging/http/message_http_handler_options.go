package messaginghttp

import (
	"github.com/xfrr/go-cqrsify/pkg/apix"
)

type HTTPMessageServerOption func(*MessageHandler)

// WithErrorMapper sets a custom domain-error -> Problem mapper.
func WithErrorMapper(mapper func(error) apix.Problem) HTTPMessageServerOption {
	return func(s *MessageHandler) { s.errorMapper = mapper }
}

// WithMaxBodyBytes sets the maximum allowed request body size (defaults to 1MiB).
func WithMaxBodyBytes(n int64) HTTPMessageServerOption {
	return func(s *MessageHandler) { s.maxBodyBytes = n }
}

// WithValidator sets the HTTP request validator (required).
func WithValidator(validator apix.HTTPRequestValidator) HTTPMessageServerOption {
	return func(s *MessageHandler) { s.messageValidator = validator }
}

// WithDecoderRegistry sets the decoder registry to use.
func WithDecoderRegistry(registry *MessageDecoderRegistry) HTTPMessageServerOption {
	return func(s *MessageHandler) { s.decoderRegistry = registry }
}
