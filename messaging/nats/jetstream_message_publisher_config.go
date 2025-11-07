package messagingnats

import (
	"time"

	"github.com/xfrr/go-cqrsify/messaging"
	"go.opentelemetry.io/otel/propagation"
)

const (
	defaultPublishRetryAttempts = 3
	defaultPublishRetryDelay    = 100 * time.Millisecond
	defaultStreamTTL            = 0 * time.Second
	defaultMaxReplyWait         = 10 * time.Second
)

type JetStreamMessagePublisherConfig struct {
	// RetryAttempts is the number of times to retry publishing a message in case of failure.
	// If zero, no retries are performed.
	RetryAttempts int
	// RetryDelay is the delay between retry attempts. If zero, no delay is applied.
	RetryDelay time.Duration
	// MaxReplyWait is the maximum time to wait for a reply. If zero, a default of 2 seconds is used.
	MaxReplyWait time.Duration
	// StreamTTL is the global time-to-live for messages. If zero, messages do not expire.
	// Note: This requires the stream to be configured with a MaxAge policy that respects message TTL.
	// If the stream is not configured accordingly, this setting will have no effect.
	// See https://docs.nats.io/nats-concepts/jetstream/js_walkthrough#id-1.-creating-a-stream for more details.
	StreamTTL time.Duration
	// MessageTTLMapping is the time-to-live for individual messages. If zero, messages do not expire.
	// Note: If both StreamTTL and MessageTTLMapping are set, the shorter duration takes precedence.
	MessageTTLMapping map[string]time.Duration // message type -> ttl
	// SubjectBuilder is a function that builds the subject for a given message.
	// If nil, DefaultSubjectBuilder is used.
	SubjectBuilder SubjectBuilder
	// ReplySubjectBuilder is a function that builds the reply subject for a given message.
	// If nil, DefaultSubjectBuilder is used.
	ReplySubjectBuilder SubjectBuilder
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

func NewJetStreamMessagePublisherConfig(opts ...JetStreamMessagePublisherConfiger) JetStreamMessagePublisherConfig {
	cfg := defaultJetStreamMessagePublisherConfig()
	for _, opt := range opts {
		opt(cfg)
	}
	return *cfg
}

func defaultJetStreamMessagePublisherConfig() *JetStreamMessagePublisherConfig {
	return &JetStreamMessagePublisherConfig{
		RetryAttempts:       defaultPublishRetryAttempts,
		RetryDelay:          defaultPublishRetryDelay,
		StreamTTL:           defaultStreamTTL,
		MaxReplyWait:        defaultMaxReplyWait,
		SubjectBuilder:      defaultSubjectBuilder,
		ReplySubjectBuilder: defaultReplySubjectBuilder,
		Serializer:          messaging.DefaultJSONSerializer,
		Deserializer:        messaging.DefaultJSONDeserializer,
		MessageTTLMapping:   make(map[string]time.Duration),
		OTELPropagator:      propagation.NewCompositeTextMapPropagator(),
	}
}

type JetStreamMessagePublisherConfiger func(*JetStreamMessagePublisherConfig)

// WithJetStreamPublishRetryAttempts sets the number of retry attempts for publishing a message.
func WithJetStreamPublishRetryAttempts(attempts int) JetStreamMessagePublisherConfiger {
	return func(cfg *JetStreamMessagePublisherConfig) {
		cfg.RetryAttempts = attempts
	}
}

// WithJetStreamPublishRetryDelay sets the delay between retry attempts.
func WithJetStreamPublishRetryDelay(delay time.Duration) JetStreamMessagePublisherConfiger {
	return func(cfg *JetStreamMessagePublisherConfig) {
		cfg.RetryDelay = delay
	}
}

// WithJetStreamPublishStreamTTL sets the global time-to-live for all messages.
func WithJetStreamPublishStreamTTL(ttl time.Duration) JetStreamMessagePublisherConfiger {
	return func(cfg *JetStreamMessagePublisherConfig) {
		cfg.StreamTTL = ttl
	}
}

// WithJetStreamPublishMaxReplyWait sets the maximum time to wait for a reply.
func WithJetStreamPublishMaxReplyWait(maxWait time.Duration) JetStreamMessagePublisherConfiger {
	return func(cfg *JetStreamMessagePublisherConfig) {
		cfg.MaxReplyWait = maxWait
	}
}

// WithJetStreamPublishMessageTTL sets the time-to-live for a specific message type.
func WithJetStreamPublishMessageTTL(messageType string, ttl time.Duration) JetStreamMessagePublisherConfiger {
	return func(cfg *JetStreamMessagePublisherConfig) {
		if cfg.MessageTTLMapping == nil {
			cfg.MessageTTLMapping = make(map[string]time.Duration)
		}
		cfg.MessageTTLMapping[messageType] = ttl
	}
}

// WithJetStreamPublishSubjectBuilder sets the subject builder function for the publisher.
func WithJetStreamPublishSubjectBuilder(builder SubjectBuilderFunc) JetStreamMessagePublisherConfiger {
	return func(cfg *JetStreamMessagePublisherConfig) {
		cfg.SubjectBuilder = builder
	}
}

// WithJetStreamPublishReplySubjectBuilder sets the reply subject builder function for the publisher.
func WithJetStreamPublishReplySubjectBuilder(builder SubjectBuilderFunc) JetStreamMessagePublisherConfiger {
	return func(cfg *JetStreamMessagePublisherConfig) {
		cfg.ReplySubjectBuilder = builder
	}
}

// WithJetStreamPublishMessageSerializer sets the message serializer for the publisher.
func WithJetStreamPublishMessageSerializer(serializer messaging.MessageSerializer) JetStreamMessagePublisherConfiger {
	return func(cfg *JetStreamMessagePublisherConfig) {
		cfg.Serializer = serializer
	}
}

// WithJetStreamPublishMessageDeserializer sets the message deserializer for the publisher.
func WithJetStreamPublishMessageDeserializer(deserializer messaging.MessageDeserializer) JetStreamMessagePublisherConfiger {
	return func(cfg *JetStreamMessagePublisherConfig) {
		cfg.Deserializer = deserializer
	}
}

// WithJetStreamPublishOTELPropagator sets the OpenTelemetry propagator for the publisher.
func WithJetStreamPublishOTELPropagator(propagator propagation.TextMapPropagator) JetStreamMessagePublisherConfiger {
	return func(cfg *JetStreamMessagePublisherConfig) {
		cfg.OTELPropagator = propagator
	}
}
