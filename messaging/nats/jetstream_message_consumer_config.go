package messagingnats

import (
	"time"

	"github.com/nats-io/nats.go/jetstream"
	"github.com/xfrr/go-cqrsify/messaging"
	"go.opentelemetry.io/otel/propagation"
)

type jetStreamConsumerConfig interface {
	jetstream.ConsumerConfig | jetstream.OrderedConsumerConfig
}

type JetStreamMessageConsumerConfig[T jetstream.ConsumerConfig | jetstream.OrderedConsumerConfig] struct {
	// ErrorHandler is a custom error handler for the JetStreamMessageConsumer.
	ErrorHandler messaging.ErrorHandler
	// MaxReplyWait is the maximum time to wait for a reply message.
	MaxReplyWait time.Duration
	// ConsumerConfig is the JetStream consumer configuration.
	ConsumerConfig T
	// Serializer is the message Serializer.
	Serializer messaging.MessageSerializer
	// Deserializer is the message Deserializer.
	Deserializer messaging.MessageDeserializer
	// OTELPropagator is the OpenTelemetry propagator for trace
	// propagation using message headers and context.
	OTELPropagator propagation.TextMapPropagator
}

func NewJetStreamMessageConsumerConfig(opts ...JetStreamMessageConsumerConfiger[jetstream.ConsumerConfig]) JetStreamMessageConsumerConfig[jetstream.ConsumerConfig] {
	cfg := &JetStreamMessageConsumerConfig[jetstream.ConsumerConfig]{
		ErrorHandler:   messaging.DefaultErrorHandler,
		MaxReplyWait:   defaultMaxReplyWait,
		ConsumerConfig: jetstream.ConsumerConfig{},
		Serializer:     messaging.DefaultJSONSerializer,
		Deserializer:   messaging.DefaultJSONDeserializer,
		OTELPropagator: propagation.NewCompositeTextMapPropagator(),
	}
	for _, opt := range opts {
		opt.apply(cfg)
	}
	return *cfg
}

func NewJetStreamOrderedMessageConsumerConfig(opts ...JetStreamMessageConsumerConfiger[jetstream.OrderedConsumerConfig]) JetStreamMessageConsumerConfig[jetstream.OrderedConsumerConfig] {
	cfg := &JetStreamMessageConsumerConfig[jetstream.OrderedConsumerConfig]{
		ErrorHandler:   messaging.DefaultErrorHandler,
		MaxReplyWait:   defaultMaxReplyWait,
		ConsumerConfig: jetstream.OrderedConsumerConfig{},
		Serializer:     messaging.DefaultJSONSerializer,
		Deserializer:   messaging.DefaultJSONDeserializer,
		OTELPropagator: propagation.NewCompositeTextMapPropagator(),
	}
	for _, opt := range opts {
		opt.apply(cfg)
	}
	return *cfg
}

type JetStreamMessageConsumerConfiger[T jetstream.ConsumerConfig | jetstream.OrderedConsumerConfig] interface {
	apply(*JetStreamMessageConsumerConfig[T])
}

type jetStreamMessageConsumerConfigFunc[T jetstream.ConsumerConfig | jetstream.OrderedConsumerConfig] func(*JetStreamMessageConsumerConfig[T])

//nolint:unused // implements interface
func (f jetStreamMessageConsumerConfigFunc[T]) apply(cfg *JetStreamMessageConsumerConfig[T]) {
	f(cfg)
}

// WithJetStreamConsumerErrorHandler sets a custom error handler for the consumer.
func WithJetStreamConsumerErrorHandler[T jetStreamConsumerConfig](handler messaging.ErrorHandler) JetStreamMessageConsumerConfiger[T] {
	return jetStreamMessageConsumerConfigFunc[T](func(cfg *JetStreamMessageConsumerConfig[T]) {
		cfg.ErrorHandler = handler
	})
}

// WithJetStreamConsumerMaxReplyWait sets the maximum time to wait for a reply message.
func WithJetStreamConsumerMaxReplyWait[T jetStreamConsumerConfig](d time.Duration) JetStreamMessageConsumerConfiger[T] {
	return jetStreamMessageConsumerConfigFunc[T](func(cfg *JetStreamMessageConsumerConfig[T]) {
		cfg.MaxReplyWait = d
	})
}

// WithJetStreamConsumerConfig sets the JetStream consumer configuration.
func WithJetStreamConsumerConfig[T jetStreamConsumerConfig](consumerConfig T) JetStreamMessageConsumerConfiger[T] {
	return jetStreamMessageConsumerConfigFunc[T](func(cfg *JetStreamMessageConsumerConfig[T]) {
		cfg.ConsumerConfig = consumerConfig
	})
}

// WithJetStreamConsumerMessageSerializer sets a custom message serializer for the consumer.
func WithJetStreamConsumerMessageSerializer[T jetStreamConsumerConfig](serializer messaging.MessageSerializer) JetStreamMessageConsumerConfiger[T] {
	return jetStreamMessageConsumerConfigFunc[T](func(cfg *JetStreamMessageConsumerConfig[T]) {
		cfg.Serializer = serializer
	})
}

// WithJetStreamConsumerMessageDeserializer sets a custom message deserializer for the consumer.
func WithJetStreamConsumerMessageDeserializer[T jetStreamConsumerConfig](deserializer messaging.MessageDeserializer) JetStreamMessageConsumerConfiger[T] {
	return jetStreamMessageConsumerConfigFunc[T](func(cfg *JetStreamMessageConsumerConfig[T]) {
		cfg.Deserializer = deserializer
	})
}

// WithJetStreamConsumerOTELPropagator sets a custom OpenTelemetry propagator for the consumer.
func WithJetStreamConsumerOTELPropagator[T jetStreamConsumerConfig](propagator propagation.TextMapPropagator) JetStreamMessageConsumerConfiger[T] {
	return jetStreamMessageConsumerConfigFunc[T](func(cfg *JetStreamMessageConsumerConfig[T]) {
		cfg.OTELPropagator = propagator
	})
}
