package messagingnats

import (
	"github.com/nats-io/nats.go/jetstream"
	"github.com/xfrr/go-cqrsify/messaging"
)

// MessageBusOption defines a function type for configuring MessageBus options.
type MessageBusOption func(*MessageBusOptions)

// MessageBusOptions holds configuration options for MessageBus.
type MessageBusOptions struct {
	subjectBuilder SubjectBuilderFunc
	errorHandler   messaging.ErrorHandler
}

// WithSubjectBuilder sets a custom SubjectBuilder for the PubSubMessageBus.
func WithSubjectBuilder(sb SubjectBuilderFunc) MessageBusOption {
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

type PubSubMessageBusOptions struct {
	MessageBusOptions

	queueName string
}

type PubSubMessageBusOption func(p *PubSubMessageBusOptions)

func WithQueueName(name string) PubSubMessageBusOption {
	return func(p *PubSubMessageBusOptions) {
		p.queueName = name
	}
}

type JetStreamMessageBusOptions struct {
	MessageBusOptions

	streamCfg jetstream.StreamConfig
}

type JetStreamMessageBusOption func(p *JetStreamMessageBusOptions)

func WithStreamConfig(cfg jetstream.StreamConfig) JetStreamMessageBusOption {
	return func(p *JetStreamMessageBusOptions) {
		p.streamCfg = cfg
	}
}
