package messagingnats

import "time"

const (
	defaultRetryAttempts = 3
	defaultRetryDelay    = 100 * time.Millisecond
	defaultStreamTTL     = 0 * time.Second
)

type JetStreamMessagePublisherConfig struct {
	// RetryAttempts is the number of times to retry publishing a message in case of failure.
	// If zero, no retries are performed.
	RetryAttempts int
	// RetryDelay is the delay between retry attempts. If zero, no delay is applied.
	RetryDelay time.Duration
	// StreamTTL is the global time-to-live for messages. If zero, messages do not expire.
	// Note: This requires the stream to be configured with a MaxAge policy that respects message TTL.
	// If the stream is not configured accordingly, this setting will have no effect.
	// See https://docs.nats.io/nats-concepts/jetstream/js_walkthrough#id-1.-creating-a-stream for more details.
	StreamTTL time.Duration
	// MessageTTL is the time-to-live for individual messages. If zero, messages do not expire.
	// Note: If both StreamTTL and MessageTTL are set, the shorter duration takes precedence.
	MessageTTL map[string]time.Duration // message type -> ttl
	// SubjectBuilder is a function that builds the subject for a given message.
	// If nil, DefaultSubjectBuilder is used.
	SubjectBuilder SubjectBuilderFunc
}

func NewJetStreamMessagePublisherConfig(opts ...JetStreamMessagePublisherOption) *JetStreamMessagePublisherConfig {
	cfg := defaultJetStreamMessagePublisherConfig()
	for _, opt := range opts {
		opt(cfg)
	}
	return cfg
}

type JetStreamMessagePublisherOption func(*JetStreamMessagePublisherConfig)

// WithJSPublishRetryAttempts sets the number of retry attempts for publishing a message.
func WithJSPublishRetryAttempts(attempts int) JetStreamMessagePublisherOption {
	return func(cfg *JetStreamMessagePublisherConfig) {
		cfg.RetryAttempts = attempts
	}
}

// WithJSPublishRetryDelay sets the delay between retry attempts.
func WithJSPublishRetryDelay(delay time.Duration) JetStreamMessagePublisherOption {
	return func(cfg *JetStreamMessagePublisherConfig) {
		cfg.RetryDelay = delay
	}
}

// WithJSPublishStreamTTL sets the global time-to-live for all messages.
func WithJSPublishStreamTTL(ttl time.Duration) JetStreamMessagePublisherOption {
	return func(cfg *JetStreamMessagePublisherConfig) {
		cfg.StreamTTL = ttl
	}
}

// WithJSPublishMessageTTL sets the time-to-live for a specific message type.
func WithJSPublishMessageTTL(messageType string, ttl time.Duration) JetStreamMessagePublisherOption {
	return func(cfg *JetStreamMessagePublisherConfig) {
		if cfg.MessageTTL == nil {
			cfg.MessageTTL = make(map[string]time.Duration)
		}
		cfg.MessageTTL[messageType] = ttl
	}
}

func defaultJetStreamMessagePublisherConfig() *JetStreamMessagePublisherConfig {
	return &JetStreamMessagePublisherConfig{
		RetryAttempts: defaultRetryAttempts,
		RetryDelay:    defaultRetryDelay,
		StreamTTL:     defaultStreamTTL,
		MessageTTL:    make(map[string]time.Duration),
	}
}
