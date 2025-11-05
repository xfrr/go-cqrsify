package messagingnats

import "time"

const (
	defaultRetryAttempts = 3
	defaultRetryDelay    = 100 * time.Millisecond
	defaultStreamTTL     = 0 * time.Second
	defaultMaxReplyWait  = 10 * time.Second
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
}

func NewJetStreamMessagePublisherConfig(opts ...JetstreamMessagePublisherConfiger) JetStreamMessagePublisherConfig {
	cfg := defaultJetStreamMessagePublisherConfig()
	for _, opt := range opts {
		opt(cfg)
	}
	return *cfg
}

func defaultJetStreamMessagePublisherConfig() *JetStreamMessagePublisherConfig {
	return &JetStreamMessagePublisherConfig{
		RetryAttempts:       defaultRetryAttempts,
		RetryDelay:          defaultRetryDelay,
		StreamTTL:           defaultStreamTTL,
		MaxReplyWait:        defaultMaxReplyWait,
		SubjectBuilder:      defaultSubjectBuilder,
		ReplySubjectBuilder: defaultReplySubjectBuilder,
		MessageTTLMapping:   make(map[string]time.Duration),
	}
}

type JetstreamMessagePublisherConfiger func(*JetStreamMessagePublisherConfig)

// WithJSPublishRetryAttempts sets the number of retry attempts for publishing a message.
func WithJSPublishRetryAttempts(attempts int) JetstreamMessagePublisherConfiger {
	return func(cfg *JetStreamMessagePublisherConfig) {
		cfg.RetryAttempts = attempts
	}
}

// WithJSPublishRetryDelay sets the delay between retry attempts.
func WithJSPublishRetryDelay(delay time.Duration) JetstreamMessagePublisherConfiger {
	return func(cfg *JetStreamMessagePublisherConfig) {
		cfg.RetryDelay = delay
	}
}

// WithJSPublishStreamTTL sets the global time-to-live for all messages.
func WithJSPublishStreamTTL(ttl time.Duration) JetstreamMessagePublisherConfiger {
	return func(cfg *JetStreamMessagePublisherConfig) {
		cfg.StreamTTL = ttl
	}
}

// WithJSPublishMaxReplyWait sets the maximum time to wait for a reply.
func WithJSPublishMaxReplyWait(maxWait time.Duration) JetstreamMessagePublisherConfiger {
	return func(cfg *JetStreamMessagePublisherConfig) {
		cfg.MaxReplyWait = maxWait
	}
}

// WithJSPublishMessageTTL sets the time-to-live for a specific message type.
func WithJSPublishMessageTTL(messageType string, ttl time.Duration) JetstreamMessagePublisherConfiger {
	return func(cfg *JetStreamMessagePublisherConfig) {
		if cfg.MessageTTLMapping == nil {
			cfg.MessageTTLMapping = make(map[string]time.Duration)
		}
		cfg.MessageTTLMapping[messageType] = ttl
	}
}

// WithJSPublishSubjectBuilder sets the subject builder function for the publisher.
func WithJSPublishSubjectBuilder(builder SubjectBuilderFunc) JetstreamMessagePublisherConfiger {
	return func(cfg *JetStreamMessagePublisherConfig) {
		cfg.SubjectBuilder = builder
	}
}

// WithJSPublishReplySubjectBuilder sets the reply subject builder function for the publisher.
func WithJSPublishReplySubjectBuilder(builder SubjectBuilderFunc) JetstreamMessagePublisherConfiger {
	return func(cfg *JetStreamMessagePublisherConfig) {
		cfg.ReplySubjectBuilder = builder
	}
}
