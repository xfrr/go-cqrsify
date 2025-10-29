package messaginghttp

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/xfrr/go-cqrsify/messaging"
)

var (
	// ErrMessageDecoderAlreadyExists is returned when trying to register a decoder for a message type and encoding that already has a decoder.
	ErrMessageDecoderAlreadyExists = errors.New("decoder for message type and encoding already exists")
	// ErrMessageDecoderNotFound is returned when no decoder is found.
	ErrMessageDecoderNotFound = errors.New("decoder not found")
	// ErrMessageDecoderNotFoundForType is returned when no decoder is found for the given message type.
	ErrMessageDecoderNotFoundForType = errors.New("decoder not found for message type")
	// ErrMessageDecoderNotFoundForEncoding is returned when no decoder is found for the given encoding.
	ErrMessageDecoderNotFoundForEncoding = errors.New("decoder not found for encoding")
)

type MessageDecoder func(*http.Request) (messaging.Message, error)

type MessageDecoderRegistry struct {
	messageDecoders map[string]map[HTTPMessageEncoding]MessageDecoder
}

func NewMessageDecoderRegistry() *MessageDecoderRegistry {
	return &MessageDecoderRegistry{
		messageDecoders: make(map[string]map[HTTPMessageEncoding]MessageDecoder),
	}
}

func (r *MessageDecoderRegistry) Register(msgType string, encoding HTTPMessageEncoding, decodeFunc MessageDecoder) error {
	if r.messageDecoders == nil {
		r.messageDecoders = make(map[string]map[HTTPMessageEncoding]MessageDecoder)
	}
	if _, exists := r.messageDecoders[msgType]; !exists {
		r.messageDecoders[msgType] = make(map[HTTPMessageEncoding]MessageDecoder)
	}
	if _, exists := r.messageDecoders[msgType][encoding]; exists {
		return fmt.Errorf("%w: type=%s, encoding=%s", ErrMessageDecoderAlreadyExists, msgType, encoding)
	}

	r.messageDecoders[msgType][encoding] = decodeFunc
	return nil
}

func (r *MessageDecoderRegistry) Get(msgType string, encoding HTTPMessageEncoding) (MessageDecoder, error) {
	if r.messageDecoders == nil {
		return nil, ErrMessageDecoderNotFound
	}

	encodingDecoders, ok := r.messageDecoders[msgType]
	if !ok {
		return nil, ErrMessageDecoderNotFoundForType
	}

	messageDecoder, ok := encodingDecoders[encoding]
	if !ok {
		return nil, ErrMessageDecoderNotFoundForEncoding
	}

	return messageDecoder, nil
}
