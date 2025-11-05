package messagingnats

import (
	"time"

	"github.com/xfrr/go-cqrsify/messaging"
)

type PubSubMessagePublisherConfig struct {
	// MaxReplyWait is the maximum time to wait for a reply. If zero, a default of 2 seconds is used.
	MaxReplyWait time.Duration
	// SubjectBuilder is a function that builds the subject for a given message.
	// If nil, DefaultSubjectBuilder is used.
	SubjectBuilder SubjectBuilder
	// ReplySubjectBuilder is a function that builds the reply subject for a given message.
	// If nil, DefaultSubjectBuilder is used.
	ReplySubjectBuilder SubjectBuilder
	// ErrorHandler is a custom error handler for the PubSubMessagePublisher.	ErrorHandler messaging.ErrorHandler
	// RetryAttempts is the number of retry attempts for publishing a message.
	ErrorHandler messaging.ErrorHandler
}

func NewPubSubMessagePublisherConfig(opts ...PubSubMessagePublisherConfiger) PubSubMessagePublisherConfig {
	cfg := defaultPubSubMessagePublisherConfig()
	for _, opt := range opts {
		opt(cfg)
	}
	return *cfg
}

func defaultPubSubMessagePublisherConfig() *PubSubMessagePublisherConfig {
	return &PubSubMessagePublisherConfig{
		MaxReplyWait:        defaultMaxReplyWait,
		SubjectBuilder:      defaultSubjectBuilder,
		ReplySubjectBuilder: defaultReplySubjectBuilder,
		ErrorHandler:        messaging.DefaultErrorHandler,
	}
}

type PubSubMessagePublisherConfiger func(*PubSubMessagePublisherConfig)

func WithPubSubPublisherMaxReplyWait(maxWait time.Duration) PubSubMessagePublisherConfiger {
	return func(cfg *PubSubMessagePublisherConfig) {
		cfg.MaxReplyWait = maxWait
	}
}

// WithPubSubPublisherSubjectBuilder sets a custom SubjectBuilder for the PubSubMessagePublisher.
func WithPubSubPublisherSubjectBuilder(sb SubjectBuilder) PubSubMessagePublisherConfiger {
	return func(cfg *PubSubMessagePublisherConfig) {
		cfg.SubjectBuilder = sb
	}
}

// WithPubSubPublisherReplySubjectBuilder sets a custom ReplySubjectBuilder for the PubSubMessagePublisher.
func WithPubSubPublisherReplySubjectBuilder(rsb SubjectBuilder) PubSubMessagePublisherConfiger {
	return func(cfg *PubSubMessagePublisherConfig) {
		cfg.ReplySubjectBuilder = rsb
	}
}

// WithPubSubPublisherErrorHandler sets a custom ErrorHandler for the PubSubMessagePublisher.
func WithPubSubPublisherErrorHandler(eh messaging.ErrorHandler) PubSubMessagePublisherConfiger {
	return func(cfg *PubSubMessagePublisherConfig) {
		cfg.ErrorHandler = eh
	}
}
