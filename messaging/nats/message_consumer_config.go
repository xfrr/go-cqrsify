package messagingnats

import (
	"github.com/xfrr/go-cqrsify/messaging"
)

// MessageConsumerConfigModifier defines a function type for configuring MessageBus options.
type MessageConsumerConfigModifier func(*MessageConsumerConfig)

// MessageConsumerConfig holds configuration options for MessageBus.
type MessageConsumerConfig struct {
	subjectBuilder SubjectBuilder
	errorHandler   messaging.ErrorHandler
	replyBuilder   ReplySubjectBuilder
}

// WithSubjectBuilder sets a custom SubjectBuilder for the PubSubMessageBus.
func WithSubjectBuilder(sb SubjectBuilder) MessageConsumerConfigModifier {
	return func(p *MessageConsumerConfig) {
		p.subjectBuilder = sb
	}
}

// WithErrorHandler sets a custom ErrorHandler for the PubSubMessageBus.
// Note that the Message could be nil if the error is not related to a specific message.
func WithErrorHandler(eh messaging.ErrorHandler) MessageConsumerConfigModifier {
	return func(p *MessageConsumerConfig) {
		p.errorHandler = eh
	}
}

type PubSubMessageBusOptions struct {
	MessageConsumerConfig

	queueName string
}

type PubSubMessageBusOption func(p *PubSubMessageBusOptions)

func WithQueueName(name string) PubSubMessageBusOption {
	return func(p *PubSubMessageBusOptions) {
		p.queueName = name
	}
}
