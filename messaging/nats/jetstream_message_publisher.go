package messagingnats

import (
	"context"
	"fmt"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/xfrr/go-cqrsify/messaging"
	"go.opentelemetry.io/otel/propagation"
)

var _ messaging.MessagePublisher = (*JetstreamMessagePublisher)(nil)

// JetstreamMessagePublisher is a publisher that uses NATS JetStream.
type JetstreamMessagePublisher struct {
	streamName string
	js         jetstream.JetStream
	cfg        JetStreamMessagePublisherConfig
}

func NewJetStreamMessagePublisher(
	js jetstream.JetStream,
	streamName string,
	opts ...JetStreamMessagePublisherConfiger,
) (*JetstreamMessagePublisher, error) {
	cfg := NewJetStreamMessagePublisherConfig(opts...)

	p := &JetstreamMessagePublisher{
		streamName: streamName,
		js:         js,
		cfg:        cfg,
	}
	return p, nil
}

// Publish implements messaging.MessageBus.
func (p *JetstreamMessagePublisher) Publish(ctx context.Context, msg ...messaging.Message) error {
	for _, m := range msg {
		data, err := p.cfg.Serializer.Serialize(m)
		if err != nil {
			return err
		}

		opts := []jetstream.PublishOpt{
			jetstream.WithExpectStream(p.streamName),
			jetstream.WithRetryAttempts(p.getRetryAttempts(m)),
			jetstream.WithRetryWait(p.getRetryWaitDuration(m)),
			jetstream.WithMsgTTL(p.getMessageTTL(m)),
		}
		if m.MessageID() != "" {
			opts = append(opts, jetstream.WithMsgID(m.MessageID()))
		}

		subject := p.cfg.SubjectBuilder.Build(m)
		if subject == "" {
			return fmt.Errorf("no subject configured for message type '%s'", m.MessageType())
		}

		headers := nats.Header{}
		// Inject tracing headers
		if p.cfg.OTELPropagator != nil {
			p.cfg.OTELPropagator.Inject(ctx, propagation.HeaderCarrier(headers))
		}

		natsMsg := &nats.Msg{
			Subject: subject,
			Data:    data,
			Header:  headers,
		}

		_, err = p.js.PublishMsg(ctx, natsMsg, opts...)
		if err != nil {
			return err
		}
	}

	return nil
}

// PublishRequest sends a request message and waits for a single reply.
func (p *JetstreamMessagePublisher) PublishRequest(ctx context.Context, msg messaging.Message) (messaging.Message, error) {
	msgSubject := p.cfg.SubjectBuilder.Build(msg)
	if msgSubject == "" {
		return nil, fmt.Errorf("no subject configured for message type '%s'", msg.MessageType())
	}

	data, err := p.cfg.Serializer.Serialize(msg)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize message: %w", err)
	}

	replySubject := p.cfg.ReplySubjectBuilder.Build(msg)
	if replySubject == "" {
		return nil, fmt.Errorf("no reply subject configured for message type '%s'", msg.MessageType())
	}

	headers := nats.Header{
		replyHeaderKey: []string{replySubject},
	}
	// Inject tracing headers
	p.cfg.OTELPropagator.Inject(ctx, propagation.HeaderCarrier(headers))

	jsMsg := &nats.Msg{
		Subject: msgSubject,
		Data:    data,
		Header:  headers,
	}

	opts := []jetstream.PublishOpt{
		jetstream.WithExpectStream(p.streamName),
		jetstream.WithRetryAttempts(p.getRetryAttempts(msg)),
		jetstream.WithRetryWait(p.getRetryWaitDuration(msg)),
		jetstream.WithMsgTTL(p.getMessageTTL(msg)),
	}
	if msg.MessageID() != "" {
		opts = append(opts, jetstream.WithMsgID(msg.MessageID()))
	}

	pubAck, err := p.js.PublishMsg(ctx, jsMsg, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to publish request message: %w", err)
	}

	// Create a temporary consumer to receive the reply
	// TODO: make consumer configuration customizable
	consumerCfg := jetstream.ConsumerConfig{
		Name:          consumerNameFromMessageType(msg.MessageType()) + fmt.Sprintf("_reply_%d", pubAck.Sequence),
		DeliverPolicy: jetstream.DeliverAllPolicy,
		AckPolicy:     jetstream.AckExplicitPolicy,
		MaxDeliver:    5,
		FilterSubject: replySubject,
		BackOff:       []time.Duration{time.Second, 2 * time.Second, 5 * time.Second},
	}

	consumer, err := p.js.CreateConsumer(ctx, p.streamName, consumerCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create consumer for reply: %w", err)
	}

	// Receive the reply message
	replyMsg, err := consumer.Next(jetstream.FetchMaxWait(p.cfg.MaxReplyWait))
	if err != nil {
		return nil, fmt.Errorf("failed to receive reply message: %w", err)
	}

	// Deserialize the reply message
	reply, err := p.cfg.Deserializer.Deserialize(replyMsg.Data())
	if err != nil {
		termErr := replyMsg.TermWithReason("deserialization_failed")
		if termErr != nil {
			return nil, fmt.Errorf("failed to terminate message after deserialization failure: %w", termErr)
		}
		return nil, fmt.Errorf("failed to deserialize reply message: %w", err)
	}

	// Acknowledge the reply message
	if err = replyMsg.Ack(); err != nil {
		return nil, fmt.Errorf("failed to ack reply message: %w", err)
	}

	return reply, nil
}

func (p *JetstreamMessagePublisher) getRetryAttempts(_ messaging.Message) int {
	// TODO: make it configurable per message type if needed
	if p.cfg.RetryAttempts > 0 {
		return p.cfg.RetryAttempts
	}

	return defaultPublishRetryAttempts
}

func (p *JetstreamMessagePublisher) getRetryWaitDuration(_ messaging.Message) time.Duration {
	// TODO: make it configurable per message type if needed
	if p.cfg.RetryDelay > 0 {
		return p.cfg.RetryDelay
	}

	return defaultPublishRetryDelay
}

// Determine the effective TTL for the message.
// If both StreamTTL and MessageTTL are set, the shorter duration takes precedence.
func (p *JetstreamMessagePublisher) getMessageTTL(m messaging.Message) time.Duration {
	if len(p.cfg.MessageTTLMapping) == 0 {
		return p.cfg.StreamTTL
	}

	if ttl, ok := p.cfg.MessageTTLMapping[m.MessageType()]; ok && ttl > 0 {
		if p.cfg.StreamTTL > 0 {
			if ttl < p.cfg.StreamTTL {
				return ttl
			}
			return p.cfg.StreamTTL
		}
		return ttl
	}
	return p.cfg.StreamTTL
}
