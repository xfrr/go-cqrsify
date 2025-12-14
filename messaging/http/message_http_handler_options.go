package messaginghttp

import (
	"github.com/xfrr/go-cqrsify/pkg/apix"
)

type MessageHandlerOption interface {
	apply(*MessageHandlerOptions)
}

type messageHandlerOptionFunc func(*MessageHandlerOptions)

func (f messageHandlerOptionFunc) apply(s *MessageHandlerOptions) {
	f(s)
}

type MessageHandlerOptions struct {
	// errorMapper maps handler errors to HTTP Problems.
	errorMapper func(error) apix.Problem

	// decoderRegistry: messageType -> encoding -> decode
	decoderRegistry *MessageDecoderRegistry

	// messageValidator validates incoming HTTP requests.
	messageValidator apix.HTTPRequestValidator

	// maxBodyBytes is the maximum allowed request body size in bytes.
	// If zero or negative, no limit is applied.
	maxBodyBytes int64
}

// WithErrorMapper sets a custom domain-error -> Problem mapper.
func WithErrorMapper(mapper func(error) apix.Problem) MessageHandlerOption {
	return messageHandlerOptionFunc(func(s *MessageHandlerOptions) { s.errorMapper = mapper })
}

// WithMaxBodyBytes sets the maximum allowed request body size (defaults to 1MiB).
func WithMaxBodyBytes(n int64) MessageHandlerOption {
	return messageHandlerOptionFunc(func(s *MessageHandlerOptions) { s.maxBodyBytes = n })
}

// WithValidator sets the HTTP request validator (required).
func WithValidator(validator apix.HTTPRequestValidator) MessageHandlerOption {
	return messageHandlerOptionFunc(func(s *MessageHandlerOptions) { s.messageValidator = validator })
}

// WithDecoderRegistry sets the decoder registry to use.
func WithDecoderRegistry(registry *MessageDecoderRegistry) MessageHandlerOption {
	return messageHandlerOptionFunc(func(s *MessageHandlerOptions) { s.decoderRegistry = registry })
}
