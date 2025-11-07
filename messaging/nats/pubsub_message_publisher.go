package messagingnats

import (
	"context"
	"fmt"

	"github.com/nats-io/nats.go"
	"github.com/xfrr/go-cqrsify/messaging"
	"go.opentelemetry.io/otel/propagation"
)

// PubSubMessagePublisher is a publisher that uses NATS JetStream.
type PubSubMessagePublisher struct {
	conn *nats.Conn
	cfg  PubSubMessagePublisherConfig
}

func NewPubSubMessagePublisher(
	conn *nats.Conn,
	opts ...PubSubMessagePublisherConfiger,
) (*PubSubMessagePublisher, error) {
	cfg := NewPubSubMessagePublisherConfig(opts...)

	p := &PubSubMessagePublisher{
		conn: conn,
		cfg:  cfg,
	}

	return p, nil
}

// Publish implements messaging.MessageBus.
func (p *PubSubMessagePublisher) Publish(ctx context.Context, messages ...messaging.Message) error {
	for _, msg := range messages {
		data, err := p.cfg.Serializer.Serialize(msg)
		if err != nil {
			p.cfg.ErrorHandler.Handle(msg, fmt.Errorf("failed to serialize message: %w", err))
			continue
		}

		subject := p.cfg.SubjectBuilder.Build(msg)
		if subject == "" {
			return fmt.Errorf("no subject configured for message type '%s'", msg.MessageType())
		}

		// Inject tracing context into message headers
		headers := nats.Header{}
		p.cfg.OTELPropagator.Inject(ctx, propagation.HeaderCarrier(headers))

		msg := &nats.Msg{
			Subject: subject,
			Data:    data,
			Header:  headers,
		}

		err = p.conn.PublishMsg(msg)
		if err != nil {
			return err
		}
	}

	return nil
}

// PublishRequest sends a request message and waits for a single reply.
func (p *PubSubMessagePublisher) PublishRequest(ctx context.Context, msg messaging.Message) (messaging.Message, error) {
	msgSubject := p.cfg.SubjectBuilder.Build(msg)
	if msgSubject == "" {
		return nil, fmt.Errorf("no subject configured for message type '%s'", msg.MessageType())
	}

	// Publish the message with a header indicating the reply subject
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

	nmsg := &nats.Msg{
		Subject: msgSubject,
		Data:    data,
		Header:  headers,
	}

	// check if context has a deadline, otherwise set a timeout based on MaxReplyWait
	ctxWithTimeout := ctx
	if _, ok := ctx.Deadline(); !ok {
		timeout := p.cfg.MaxReplyWait
		if timeout <= 0 {
			timeout = defaultMaxReplyWait
		}

		var cancel context.CancelFunc
		ctxWithTimeout, cancel = context.WithTimeout(ctx, timeout)
		defer cancel()
	}

	natsMsg, err := p.conn.RequestMsgWithContext(ctxWithTimeout, nmsg)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	replyMsg, err := p.cfg.Deserializer.Deserialize(natsMsg.Data)
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize reply message: %w", err)
	}

	return replyMsg, nil
}
