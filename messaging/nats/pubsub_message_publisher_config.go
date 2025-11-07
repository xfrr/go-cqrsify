package messagingnats

import (
	"time"

	"github.com/xfrr/go-cqrsify/messaging"
	"go.opentelemetry.io/otel/propagation"
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
	// Serializer is the message serializer to use for publishing messages.
	// If nil, a default JSON serializer is used.
	Serializer messaging.MessageSerializer
	// Deserializer is the message deserializer to use for receiving messages.
	// If nil, a default JSON deserializer is used.
	Deserializer messaging.MessageDeserializer
	// OTELPropagator is the OpenTelemetry propagator for trace
	// propagation using message headers and context.
	OTELPropagator propagation.TextMapPropagator
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
		Serializer:          messaging.DefaultJSONSerializer,
		Deserializer:        messaging.DefaultJSONDeserializer,
		OTELPropagator:      propagation.NewCompositeTextMapPropagator(),
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

// WithPubSubPublisherSerializer sets a custom MessageSerializer for the PubSubMessagePublisher.
func WithPubSubPublisherSerializer(serializer messaging.MessageSerializer) PubSubMessagePublisherConfiger {
	return func(cfg *PubSubMessagePublisherConfig) {
		cfg.Serializer = serializer
	}
}

// WithPubSubPublisherDeserializer sets a custom MessageDeserializer for the PubSubMessagePublisher.
func WithPubSubPublisherDeserializer(deserializer messaging.MessageDeserializer) PubSubMessagePublisherConfiger {
	return func(cfg *PubSubMessagePublisherConfig) {
		cfg.Deserializer = deserializer
	}
}

// WithPubSubPublisherOTELPropagator sets a custom OpenTelemetry propagator for the PubSubMessagePublisher.
func WithPubSubPublisherOTELPropagator(propagator propagation.TextMapPropagator) PubSubMessagePublisherConfiger {
	return func(cfg *PubSubMessagePublisherConfig) {
		cfg.OTELPropagator = propagator
	}
}
