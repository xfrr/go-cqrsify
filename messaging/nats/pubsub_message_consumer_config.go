package messagingnats

import (
	"time"

	"github.com/xfrr/go-cqrsify/messaging"
)

type PubSubMessageConsumerConfig struct {
	// Subject is the NATS subject to subscribe to. If empty, subscribe to all subjects.
	Subject string
	// MaxReplyWait is the maximum time to wait for a reply. If zero, a default of 2 seconds is used.
	MaxReplyWait time.Duration
	// SubjectBuilder is a function that builds the subject for a given message.
	// If nil, DefaultSubjectBuilder is used.
	SubjectBuilder SubjectBuilder
	// ErrorHandler is a custom error handler subscribe for errors occurring during message consumption.
	ErrorHandler messaging.ErrorHandler
	// Serializer is the message serializer to use for receiving messages.
	// If nil, a default JSON serializer is used.
	Serializer messaging.MessageSerializer
	// Deserializer is the message deserializer to use for receiving messages.
	// If nil, a default JSON deserializer is used.
	Deserializer messaging.MessageDeserializer
}

func NewPubSubMessageConsumerConfig(opts ...PubSubMessageConsumerConfiger) PubSubMessageConsumerConfig {
	cfg := defaultPubSubMessageConsumerConfig()
	for _, opt := range opts {
		opt(cfg)
	}
	return *cfg
}

func defaultPubSubMessageConsumerConfig() *PubSubMessageConsumerConfig {
	return &PubSubMessageConsumerConfig{
		MaxReplyWait:   defaultMaxReplyWait,
		SubjectBuilder: defaultSubjectBuilder,
		ErrorHandler:   messaging.DefaultErrorHandler,
		Serializer:     messaging.DefaultJSONSerializer,
		Deserializer:   messaging.DefaultJSONDeserializer,
	}
}

type PubSubMessageConsumerConfiger func(*PubSubMessageConsumerConfig)

func WithPubSubConsumerMaxReplyWait(maxWait time.Duration) PubSubMessageConsumerConfiger {
	return func(cfg *PubSubMessageConsumerConfig) {
		cfg.MaxReplyWait = maxWait
	}
}

// WithPubSubConsumerSubjectBuilder sets a custom SubjectBuilder for the PubSubMessageConsumer.
func WithPubSubConsumerSubjectBuilder(sb SubjectBuilder) PubSubMessageConsumerConfiger {
	return func(cfg *PubSubMessageConsumerConfig) {
		cfg.SubjectBuilder = sb
	}
}

// WithPubSubConsumerErrorHandler sets a custom ErrorHandler for the PubSubMessageConsumer.
func WithPubSubConsumerErrorHandler(eh messaging.ErrorHandler) PubSubMessageConsumerConfiger {
	return func(cfg *PubSubMessageConsumerConfig) {
		cfg.ErrorHandler = eh
	}
}

// WithPubSubConsumerSubject sets the NATS subject to subscribe to.
func WithPubSubConsumerSubject(subject string) PubSubMessageConsumerConfiger {
	return func(cfg *PubSubMessageConsumerConfig) {
		cfg.Subject = subject
	}
}

// WithPubSubConsumerMessageSerializer sets a custom Serializer for the PubSubMessageConsumer.
func WithPubSubConsumerMessageSerializer(serializer messaging.MessageSerializer) PubSubMessageConsumerConfiger {
	return func(cfg *PubSubMessageConsumerConfig) {
		cfg.Serializer = serializer
	}
}

// WithPubSubConsumerMessageDeserializer sets a custom Deserializer for the PubSubMessageConsumer.
func WithPubSubConsumerMessageDeserializer(deserializer messaging.MessageDeserializer) PubSubMessageConsumerConfiger {
	return func(cfg *PubSubMessageConsumerConfig) {
		cfg.Deserializer = deserializer
	}
}
