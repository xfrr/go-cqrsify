package messaginghttp

import "github.com/xfrr/go-cqrsify/pkg/apix"

type HTTPMessageServerOption func(*HTTPMessageServer)

// WithErrorMapper sets a custom domain-error -> Problem mapper.
func WithErrorMapper(mapper func(error) apix.Problem) HTTPMessageServerOption {
	return func(s *HTTPMessageServer) { s.errorMapper = mapper }
}

// WithMaxBodyBytes sets the maximum allowed request body size (defaults to 1MiB).
func WithMaxBodyBytes(n int64) HTTPMessageServerOption {
	return func(s *HTTPMessageServer) { s.maxBodyBytes = n }
}

// WithValidator sets the HTTP request validator (required).
func WithValidator(validator apix.HTTPRequestValidator) HTTPMessageServerOption {
	return func(s *HTTPMessageServer) { s.validator = validator }
}
