package messagingnats

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/xfrr/go-cqrsify/messaging"
)

// Ensure PubSubMessageBus implements the MessageBus interface.
var _ messaging.MessageBus = (*PubSubMessageBus)(nil)

// PubSubMessageBus is a NATS-based implementation of the MessageBus interface.
// It provides methods for publishing and subscribing to messages using NATS as the underlying message bus.
type PubSubMessageBus struct {
	mu sync.Mutex

	conn *nats.Conn

	subjectBuilder SubjectBuilder
	serializer     messaging.MessageSerializer
	deserializer   messaging.MessageDeserializer

	handlers     map[string][]messaging.MessageHandler[messaging.Message]
	errorHandler messaging.ErrorHandler
}

func NewPubSubMessageBus(
	conn *nats.Conn,
	serializer messaging.MessageSerializer,
	deserializer messaging.MessageDeserializer,
	opts ...PubSubMessageBusOption,
) *PubSubMessageBus {
	busOptions := PubSubMessageBusOptions{
		MessageBusOptions: MessageBusOptions{
			subjectBuilder: DefaultSubjectBuilder,
			errorHandler:   messaging.DefaultErrorHandler,
		},
	}
	for _, opt := range opts {
		opt(&busOptions)
	}

	return &PubSubMessageBus{
		conn:           conn,
		serializer:     serializer,
		deserializer:   deserializer,
		subjectBuilder: busOptions.subjectBuilder,
		errorHandler:   busOptions.errorHandler,
		handlers:       make(map[string][]messaging.MessageHandler[messaging.Message]),
	}
}

// Publish implements messaging.MessageBus.
func (p *PubSubMessageBus) Publish(_ context.Context, msg ...messaging.Message) error {
	for i, m := range msg {
		data, err := p.serializer.Serialize(m)
		if err != nil {
			return err
		}

		subject := p.subjectBuilder(m)
		if err = p.conn.Publish(subject, data); err != nil {
			return fmt.Errorf("failed to publish message %d to subject %s: %w", i, subject, err)
		}
	}

	return nil
}

// PublishRequest sends a request message and waits for a single reply.
func (p *PubSubMessageBus) PublishRequest(ctx context.Context, msg messaging.Message) (messaging.Message, error) {
	if msg == nil {
		return nil, errors.New("nil message provided")
	}

	data, err := p.serializer.Serialize(msg)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize message: %w", err)
	}

	subject := p.subjectBuilder(msg)
	natsMsg, err := p.conn.RequestWithContext(ctx, subject, data)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	replyMsg, err := p.deserializer.Deserialize(natsMsg.Data)
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize reply message: %w", err)
	}

	return replyMsg, nil
}

func (p *PubSubMessageBus) Subscribe(ctx context.Context, msgType string, handler messaging.MessageHandler[messaging.Message]) (messaging.UnsubscribeFunc, error) {
	if p.handlers == nil {
		p.handlers = make(map[string][]messaging.MessageHandler[messaging.Message])
	}

	p.mu.Lock()
	p.handlers[msgType] = append(p.handlers[msgType], handler)
	p.mu.Unlock()

	subject := msgType
	sub, err := p.conn.Subscribe(subject, p.natsMsgHandler(ctx, handler))
	if err != nil {
		return nil, fmt.Errorf("failed to subscribe to subject %s: %w", subject, err)
	}

	return p.unsubscribeFn(subject, sub, handler), nil
}

func (p *PubSubMessageBus) unsubscribeFn(
	subject string,
	sub *nats.Subscription,
	handler messaging.MessageHandler[messaging.Message],
) func() {
	return func() {
		err := sub.Unsubscribe()
		if err != nil {
			p.errorHandler(nil, fmt.Errorf("failed to unsubscribe from subject %s: %w", subject, err))
		}

		p.mu.Lock()
		defer p.mu.Unlock()

		handlers := p.handlers[subject]
		for i := range handlers {
			if &handlers[i] == &handler {
				handlers = append(handlers[:i], handlers[i+1:]...)
				break
			}
		}

		if len(handlers) == 0 {
			delete(p.handlers, subject)
		}
	}
}

func (p *PubSubMessageBus) natsMsgHandler(ctx context.Context, handler messaging.MessageHandler[messaging.Message]) nats.MsgHandler {
	return func(m *nats.Msg) {
		msg, err := p.deserializer.Deserialize(m.Data)
		if err != nil {
			p.errorHandler(nil, err)
			return
		}

		if err = handler.Handle(ctx, msg); err != nil {
			p.errorHandler(msg, err)
			return
		}

		if msgReplier, ok := msg.(messaging.ReplyableMessage); ok && m.Reply != "" {
			replyCtx, cancel := context.WithTimeout(ctx, messaging.DefaultReplyTimeoutSeconds*time.Second)
			defer cancel()

			replyMsg, replyErr := msgReplier.GetReply(replyCtx)
			if replyErr != nil {
				p.errorHandler(msg, fmt.Errorf("failed to get reply message: %w", replyErr))
				return
			}

			replyData, serializeErr := p.serializer.Serialize(replyMsg)
			if serializeErr != nil {
				p.errorHandler(replyMsg, fmt.Errorf("failed to serialize reply message: %w", serializeErr))
				return
			}

			if err = m.Respond(replyData); err != nil {
				p.errorHandler(replyMsg, fmt.Errorf("failed to send reply message: %w", err))
				return
			}
		}
	}
}
