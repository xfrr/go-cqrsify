package messagingnats

import (
	"context"
	"errors"
	"fmt"

	"github.com/nats-io/nats.go"
	"github.com/xfrr/go-cqrsify/messaging"
	"go.opentelemetry.io/otel/propagation"
)

var _ messaging.MessageConsumer = (*PubSubMessageConsumer)(nil)

// PubSubMessageConsumer consumes messages from a NATS subject (core pub/sub).
type PubSubMessageConsumer struct {
	conn *nats.Conn
	cfg  PubSubMessageConsumerConfig
}

func NewPubSubMessageConsumer(
	conn *nats.Conn,
	opts ...PubSubMessageConsumerConfiger,
) (*PubSubMessageConsumer, error) {
	if conn == nil {
		return nil, errors.New("nats connection cannot be nil")
	}

	cfg := NewPubSubMessageConsumerConfig(opts...)

	return &PubSubMessageConsumer{
		conn: conn,
		cfg:  cfg,
	}, nil
}

// Subscribe implements messaging.MessageConsumer.
func (p *PubSubMessageConsumer) Subscribe(
	ctx context.Context,
	handler messaging.MessageHandler[messaging.Message],
) (messaging.UnsubscribeFunc, error) {
	if handler == nil {
		return nil, errors.New("handler cannot be nil")
	}

	sub, err := p.conn.Subscribe(p.cfg.Subject, p.handleMessage(ctx, handler))
	if err != nil {
		return nil, fmt.Errorf("failed to subscribe to subject %q: %w", p.cfg.Subject, err)
	}
	return p.unsubscribeFn(p.cfg.Subject, sub), nil
}

// SubscribeWithReply implements messaging.MessageConsumerReplier.
func (p *PubSubMessageConsumer) SubscribeWithReply(
	ctx context.Context,
	handler messaging.MessageHandlerWithReply[messaging.Message, messaging.Message],
) (messaging.UnsubscribeFunc, error) {
	if handler == nil {
		return nil, errors.New("handler cannot be nil")
	}

	sub, err := p.conn.Subscribe(p.cfg.Subject, p.handleMessageWithReply(ctx, handler))
	if err != nil {
		return nil, fmt.Errorf("failed to subscribe to subject %q: %w", p.cfg.Subject, err)
	}
	return p.unsubscribeFn(p.cfg.Subject, sub), nil
}

func (p *PubSubMessageConsumer) handleMessage(
	ctx context.Context,
	handler messaging.MessageHandler[messaging.Message],
) nats.MsgHandler {
	return func(nm *nats.Msg) {
		m := p.deserializeOrTerm(nm)
		if m == nil {
			return
		}

		msgCtx := p.propagateTracingContext(ctx, nm)
		if err := handler.Handle(msgCtx, m); err != nil {
			// TODO: decide whether to Term or Nak based on error type
			p.errHandle(m, fmt.Errorf("failed to handle message: %w", err))
			return
		}
	}
}

func (p *PubSubMessageConsumer) handleMessageWithReply(
	ctx context.Context,
	handler messaging.MessageHandlerWithReply[messaging.Message, messaging.Message],
) nats.MsgHandler {
	return func(nm *nats.Msg) {
		m := p.deserializeOrTerm(nm)
		if m == nil {
			return
		}

		msgCtx := p.propagateTracingContext(ctx, nm)
		reply, err := handler.Handle(msgCtx, m)
		if err != nil {
			// TODO: decide whether to Term or Nak based on error type
			p.errHandle(m, fmt.Errorf("failed to handle message: %w", err))
			return
		}
		if reply == nil {
			p.errHandle(m, errors.New("handler returned nil reply"))
			return
		}

		if replyErr := p.sendReply(msgCtx, nm, reply); replyErr != nil {
			// sendReply already handled error reporting and termination.
			return
		}
	}
}

// sendReply serializes the reply, injects tracing headers and sends the reply message;
// it reports errors via errAndTerm and returns a non-nil error when something failed.
func (p *PubSubMessageConsumer) sendReply(msgCtx context.Context, nm *nats.Msg, reply messaging.Message) error {
	data, serr := p.cfg.Serializer.Serialize(reply)
	if serr != nil {
		p.errHandle(reply, fmt.Errorf("failed to serialize reply message: %w", serr))
		return serr
	}

	if nm.Reply == "" {
		p.errHandle(reply, errors.New("no reply subject on incoming message"))
		return errors.New("no reply subject on incoming message")
	}

	headers := nats.Header{}
	p.cfg.OTELPropagator.Inject(msgCtx, propagation.HeaderCarrier(headers))

	if respondErr := nm.RespondMsg(&nats.Msg{
		Subject: nm.Reply,
		Data:    data,
		Header:  headers,
	}); respondErr != nil {
		p.errHandle(reply, fmt.Errorf("failed to send reply message: %w", respondErr))
		return respondErr
	}

	return nil
}

func (p *PubSubMessageConsumer) unsubscribeFn(subject string, sub *nats.Subscription) messaging.UnsubscribeFunc {
	return func() error {
		if err := sub.Drain(); err != nil {
			return fmt.Errorf("failed to unsubscribe from subject %q: %w", subject, err)
		}
		return nil
	}
}

func (p *PubSubMessageConsumer) deserializeOrTerm(nm *nats.Msg) messaging.Message {
	m, err := p.cfg.Deserializer.Deserialize(nm.Data)
	if err != nil {
		p.errHandle(nil, fmt.Errorf("failed to deserialize message: %w", err))
		if termErr := nm.Term(); termErr != nil {
			p.errHandle(nil, fmt.Errorf("failed to term message: %w", termErr))
		}
		return nil
	}
	if m == nil {
		p.errHandle(nil, errors.New("nil message after deserialization"))
		if termErr := nm.Term(); termErr != nil {
			p.errHandle(nil, fmt.Errorf("failed to term message: %w", termErr))
		}
		return nil
	}
	return m
}

func (p *PubSubMessageConsumer) propagateTracingContext(parent context.Context, nm *nats.Msg) context.Context {
	return p.cfg.OTELPropagator.Extract(parent, propagation.HeaderCarrier(nm.Header))
}

func (p *PubSubMessageConsumer) errHandle(msg messaging.Message, err error) {
	if p.cfg.ErrorHandler != nil {
		p.cfg.ErrorHandler.Handle(msg, err)
	}
}
