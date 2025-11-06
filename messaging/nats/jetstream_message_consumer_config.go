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
	// OrderedConsumerConfig is the JetStream ordered consumer configuration.
	OrderedConsumerConfig jetstream.OrderedConsumerConfig
	// IsOrdered indicates whether the consumer is ordered.
	IsOrdered bool
}

func NewJetStreamMessageConsumerConfig(opts ...JetStreamMessageConsumerConfiger) JetStreamMessageConsumerConfig {
	cfg := &JetStreamMessageConsumerConfig{
		ErrorHandler:   messaging.DefaultErrorHandler,
		MaxReplyWait:   defaultMaxReplyWait,
		ConsumerConfig: jetstream.ConsumerConfig{},
	}
	for _, opt := range opts {
		opt.apply(cfg)
	}
	return *cfg
}

func NewJetStreamOrderedMessageConsumerConfig(opts ...JetStreamMessageConsumerConfiger) JetStreamMessageConsumerConfig {
	cfg := &JetStreamMessageConsumerConfig{
		IsOrdered:             true,
		ErrorHandler:          messaging.DefaultErrorHandler,
		MaxReplyWait:          defaultMaxReplyWait,
		OrderedConsumerConfig: jetstream.OrderedConsumerConfig{},
	}
	for _, opt := range opts {
		opt.apply(cfg)
	}
	return *cfg
}

type JetStreamMessageConsumerConfiger interface {
	apply(*JetStreamMessageConsumerConfig)
}

type jetStreamMessageConsumerConfigFunc func(*JetStreamMessageConsumerConfig)

func (f jetStreamMessageConsumerConfigFunc) apply(cfg *JetStreamMessageConsumerConfig) {
	f(cfg)
}

// WithJetStreamConsumerErrorHandler sets a custom error handler for the consumer.
func WithJetStreamConsumerErrorHandler(handler messaging.ErrorHandler) JetStreamMessageConsumerConfiger {
	return jetStreamMessageConsumerConfigFunc(func(cfg *JetStreamMessageConsumerConfig) {
		cfg.ErrorHandler = handler
	})
}

// WithJetStreamConsumerMaxReplyWait sets the maximum time to wait for a reply message.
func WithJetStreamConsumerMaxReplyWait(d time.Duration) JetStreamMessageConsumerConfiger {
	return jetStreamMessageConsumerConfigFunc(func(cfg *JetStreamMessageConsumerConfig) {
		cfg.MaxReplyWait = d
	})
}

// WithJetStreamConsumerConfig sets the JetStream consumer configuration.
func WithJetStreamConsumerConfig(consumerConfig jetstream.ConsumerConfig) JetStreamMessageConsumerConfiger {
	return jetStreamMessageConsumerConfigFunc(func(cfg *JetStreamMessageConsumerConfig) {
		cfg.ConsumerConfig = consumerConfig
	})
}

// WithJetStreamOrderedConsumerConfig sets the JetStream ordered consumer configuration.
func WithJetStreamOrderedConsumerConfig(orderedConsumerConfig jetstream.OrderedConsumerConfig) JetStreamMessageConsumerConfiger {
	return jetStreamMessageConsumerConfigFunc(func(cfg *JetStreamMessageConsumerConfig) {
		cfg.OrderedConsumerConfig = orderedConsumerConfig
	})
}
