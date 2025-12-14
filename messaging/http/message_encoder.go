package messaginghttp

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/xfrr/go-cqrsify/messaging"
	"github.com/xfrr/go-cqrsify/pkg/apix"
)

var (
	ErrMessageEncoderAlreadyExists       = errors.New("message encoder already exists")
	ErrMessageEncoderNotFound            = errors.New("message encoder not found")
	ErrMessageEncoderNotFoundForEncoding = errors.New("message encoder not found for encoding")
)

type EncodedMessage struct {
	ContentType apix.ContentType
	Data        []byte
}

type MessageEncoder struct {
	contentType apix.ContentType
	encodeFn    func(ctx context.Context, msg messaging.Message) (*EncodedMessage, error)
}

type MessageEncoderRegistry struct {
	messageEncoders map[string]MessageEncoder
}

func NewMessageEncoderRegistry() *MessageEncoderRegistry {
	return &MessageEncoderRegistry{
		messageEncoders: make(map[string]MessageEncoder),
	}
}

func (r *MessageEncoderRegistry) Register(msgType string, encodeFunc MessageEncoder) error {
	if r.messageEncoders == nil {
		r.messageEncoders = make(map[string]MessageEncoder)
	}
	if _, exists := r.messageEncoders[msgType]; !exists {
		r.messageEncoders[msgType] = MessageEncoder{}
	}
	r.messageEncoders[msgType] = encodeFunc
	return nil
}

func (r *MessageEncoderRegistry) Get(msgType string) (*MessageEncoder, error) {
	if r.messageEncoders == nil {
		return nil, ErrMessageEncoderNotFound
	}

	messageEncoder, ok := r.messageEncoders[msgType]
	if !ok {
		return nil, ErrMessageEncoderNotFoundForEncoding
	}

	return &messageEncoder, nil
}

// RegisterSingleDocumentMessageEncoder registers a JSON:API single document encoder for the given message type.
// If an encoder for the same message type and encoding already exists, an error is returned.
func RegisterSingleDocumentMessageEncoder[T messaging.Message, A messaging.Message](
	handler *MessageWithReplyHandler,
	msgType string,
	encodeFunc func(context.Context, T) (apix.SingleDocument[A], error),
) error {
	return handler.encoderRegistry.Register(
		msgType,
		MessageEncoder{
			contentType: apix.ContentTypeJSONAPI,
			encodeFn: func(ctx context.Context, msg messaging.Message) (*EncodedMessage, error) {
				castedMsg, ok := msg.(T)
				if !ok {
					return nil, fmt.Errorf("invalid message type: expected %T, got %T", castedMsg, msg)
				}

				doc, err := encodeFunc(ctx, castedMsg)
				if err != nil {
					return nil, err
				}

				body, err := json.Marshal(doc)
				if err != nil {
					return nil, err
				}

				return &EncodedMessage{
					ContentType: apix.ContentTypeJSONAPI,
					Data:        body,
				}, nil
			},
		},
	)
}

// RegisterManyDocumentMessageEncoder registers a JSON:API many document encoder for the given message type.
// If an encoder for the same message type and encoding already exists, an error is returned.
func RegisterManyDocumentMessageEncoder[T messaging.Message, A any](
	handler *MessageWithReplyHandler,
	msgType string,
	encodeFunc func(context.Context, T) (apix.ManyDocument[A], error),
) error {
	return handler.encoderRegistry.Register(
		msgType,
		MessageEncoder{
			contentType: apix.ContentTypeJSONAPI,
			encodeFn: func(ctx context.Context, msg messaging.Message) (*EncodedMessage, error) {
				castedMsg, ok := msg.(T)
				if !ok {
					return nil, fmt.Errorf("invalid message type: expected %T, got %T", castedMsg, msg)
				}

				doc, err := encodeFunc(ctx, castedMsg)
				if err != nil {
					return nil, err
				}

				body, err := json.Marshal(doc)
				if err != nil {
					return nil, err
				}

				return &EncodedMessage{
					ContentType: apix.ContentTypeJSONAPI,
					Data:        body,
				}, nil
			},
		},
	)
}
