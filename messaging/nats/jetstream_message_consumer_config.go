package messagingnats

import (
	"time"

	"github.com/nats-io/nats.go/jetstream"
	"github.com/xfrr/go-cqrsify/messaging"
)

type JetStreamMessageConsumerConfig struct {
	// ErrorHandler is a custom error handler for the JetStreamMessageConsumer.
	ErrorHandler messaging.ErrorHandler
	// MaxReplyWait is the maximum time to wait for a reply message.
	MaxReplyWait time.Duration
	// ConsumerConfig is the JetStream consumer configuration.
	ConsumerConfig jetstream.ConsumerConfig
}

func defaultJetStreamMessageConsumerConfig() *JetStreamMessageConsumerConfig {
	return &JetStreamMessageConsumerConfig{
		ErrorHandler: messaging.DefaultErrorHandler,
		MaxReplyWait: defaultMaxReplyWait,
		ConsumerConfig: jetstream.ConsumerConfig{
			DeliverPolicy: jetstream.DeliverAllPolicy,
			AckPolicy:     jetstream.AckExplicitPolicy,
		},
	}
}

func NewJetStreamMessageConsumerConfig(opts ...MessageConsumerConfiger) JetStreamMessageConsumerConfig {
	cfg := defaultJetStreamMessageConsumerConfig()
	for _, opt := range opts {
		opt(cfg)
	}
	return *cfg
}

type MessageConsumerConfiger func(*JetStreamMessageConsumerConfig)

func WithJSConsumerErrorHandler(handler messaging.ErrorHandler) MessageConsumerConfiger {
	return func(cfg *JetStreamMessageConsumerConfig) {
		cfg.ErrorHandler = handler
	}
}

func WithJSConsumerMaxReplyWait(maxWait time.Duration) MessageConsumerConfiger {
	return func(cfg *JetStreamMessageConsumerConfig) {
		cfg.MaxReplyWait = maxWait
	}
}

func WithConsumerConfig(jsCfg jetstream.ConsumerConfig) MessageConsumerConfiger {
	return func(cfg *JetStreamMessageConsumerConfig) {
		cfg.ConsumerConfig = jsCfg
	}
}
