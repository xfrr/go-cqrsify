package messagingnats

import "github.com/xfrr/go-cqrsify/messaging"

// MessageBusOption defines a function type for configuring MessageBus options.
type MessageBusOption func(*MessageBusOptions)

// MessageBusOptions holds configuration options for MessageBus.
type MessageBusOptions struct {
	subjectBuilder SubjectBuilder
	errorHandler   messaging.ErrorHandler
}

// WithSubjectBuilder sets a custom SubjectBuilder for the PubSubMessageBus.
func WithSubjectBuilder(sb SubjectBuilder) MessageBusOption {
	return func(p *MessageBusOptions) {
		p.subjectBuilder = sb
	}
}

// WithErrorHandler sets a custom ErrorHandler for the PubSubMessageBus.
// Note that the Message could be nil if the error is not related to a specific message.
func WithErrorHandler(eh messaging.ErrorHandler) MessageBusOption {
	return func(p *MessageBusOptions) {
		p.errorHandler = eh
	}
}
