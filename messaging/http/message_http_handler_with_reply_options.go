package messaginghttp

type MessageHandlerWithReplyOption interface {
	apply(*MessageHandlerWithReplyOptions)
}

type messageHandlerWithReplyOptionFunc func(*MessageHandlerWithReplyOptions)

func (f messageHandlerWithReplyOptionFunc) apply(s *MessageHandlerWithReplyOptions) {
	f(s)
}

type MessageHandlerWithReplyOptions struct {
	MessageHandlerOptions

	// encoderRegistry: encoding -> encode
	encoderRegistry *MessageEncoderRegistry
}

// WithEncoderRegistry sets the encoder registry to use.
func WithEncoderRegistry(registry *MessageEncoderRegistry) MessageHandlerWithReplyOption {
	return messageHandlerWithReplyOptionFunc(func(s *MessageHandlerWithReplyOptions) { s.encoderRegistry = registry })
}
