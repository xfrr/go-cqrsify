package messaginghttp

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"mime"
	"net/http"

	"github.com/xfrr/go-cqrsify/messaging"
	"github.com/xfrr/go-cqrsify/pkg/apix"
)

// MessageWithReplyHandler is an HTTP server for receiving messages.
type MessageWithReplyHandler struct {
	inner MessageHandler

	// messagePublisher publishes messages and waits for replies.
	messagePublisher messaging.MessagePublisherReplier

	// encoderRegistry: encoding -> encode
	encoderRegistry *MessageEncoderRegistry
}

// NewMessageWithReplyHandler creates a new MessageWithReplyHandler with the given MessagePublisher and options.
func NewMessageWithReplyHandler(msgPublisher messaging.MessagePublisherReplier, opts ...MessageHandlerWithReplyOption) *MessageWithReplyHandler {
	cfg := &MessageHandlerWithReplyOptions{
		encoderRegistry: NewMessageEncoderRegistry(),
		MessageHandlerOptions: MessageHandlerOptions{
			maxBodyBytes:    defaultMaxBodyBytes,
			decoderRegistry: NewMessageDecoderRegistry(),
		},
	}
	for _, opt := range opts {
		opt.apply(cfg)
	}

	return &MessageWithReplyHandler{
		inner: MessageHandler{
			maxBodyBytes:     cfg.maxBodyBytes,
			decoderRegistry:  cfg.decoderRegistry,
			errorMapper:      cfg.errorMapper,
			messageValidator: cfg.messageValidator,
		},
		messagePublisher: msgPublisher,
		encoderRegistry:  cfg.encoderRegistry,
	}
}

// ServeHTTP implements http.Handler.
func (h *MessageWithReplyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if h.messagePublisher == nil {
		apix.WriteProblem(w, apix.NewInternalServerErrorProblem("no message bus configured"))
		return
	}

	if h.inner.maxBodyBytes > 0 {
		r.Body = http.MaxBytesReader(w, r.Body, h.inner.maxBodyBytes)
	}
	defer r.Body.Close()

	if h.inner.messageValidator != nil {
		if problem := h.inner.messageValidator.Validate(r.Context(), r); problem != nil {
			apix.WriteProblem(w, *problem)
			return
		}
	}

	msg, problem := h.decodeMessageFromHTTPRequest(r)
	if problem != nil {
		apix.WriteProblem(w, *problem)
		return
	}

	replyMsg, err := h.messagePublisher.PublishRequest(r.Context(), msg)
	if err != nil {
		h.handleError(w, err)
		return
	}

	encodedReply, err := h.encodeMessage(r.Context(), replyMsg)
	if err != nil {
		h.handleError(w, err)
		return
	}

	// TODO: Allows to use PRG (Post/Redirect/Get) pattern or similar to avoid caching issues with POST responses.
	_, _ = apix.Write(
		w,
		encodedReply.Data,
		apix.WithHeader(apix.ContentTypeHeaderKey, encodedReply.ContentType.String()),
		// Redirect to GET to avoid POST response caching issues
		apix.WithStatusCode(http.StatusOK),
	)
}

func (h *MessageWithReplyHandler) handleError(w http.ResponseWriter, err error) {
	if h.inner.errorMapper != nil {
		apix.WriteProblem(w, h.inner.errorMapper(err))
		return
	}
	apix.WriteProblem(w, apix.NewInternalServerErrorProblem(err.Error()))
}

func (h *MessageWithReplyHandler) decodeMessageFromHTTPRequest(r *http.Request) (messaging.Message, *apix.Problem) {
	if h.inner.decoderRegistry == nil {
		problem := apix.NewInternalServerErrorProblem("no message decoders configured")
		return nil, &problem
	}

	mediaType, _, err := mime.ParseMediaType(r.Header.Get(apix.ContentTypeHeaderKey))
	if err != nil {
		problem := apix.NewUnsupportedMediaTypeProblem(fmt.Sprintf("invalid Content-Type header: %v", err))
		return nil, &problem
	}

	switch mediaType {
	case apix.ContentTypeJSONAPI.String():
		return h.decodeJSONAPIMessage(r, HTTPMessageEncodingJSONAPI)
	default:
		problem := apix.NewUnsupportedMediaTypeProblem(fmt.Sprintf("unsupported content type: %s", mediaType))
		return nil, &problem
	}
}

func (h *MessageWithReplyHandler) encodeMessage(ctx context.Context, msg messaging.Message) (*EncodedMessage, error) {
	if h.encoderRegistry == nil {
		return nil, errors.New("no message encoders configured")
	}

	encoder, err := h.encoderRegistry.Get(msg.MessageType())
	if err != nil {
		return nil, fmt.Errorf("failed to get encoder for message reply type %q: %w", msg.MessageType(), err)
	}

	encodedMsg, err := encoder.encodeFn(ctx, msg)
	if err != nil {
		return nil, fmt.Errorf("failed to encode message: %w", err)
	}

	return &EncodedMessage{
		ContentType: encodedMsg.ContentType,
		Data:        encodedMsg.Data,
	}, nil
}

func (h *MessageWithReplyHandler) decodeJSONAPIMessage(r *http.Request, encoding HTTPMessageEncoding) (messaging.Message, *apix.Problem) {
	// Read entire (bounded) body so we can unmarshal multiple times.
	body, err := io.ReadAll(r.Body)
	if err != nil {
		problem := apix.NewBadRequestProblem(fmt.Sprintf("failed to read request body: %s", err))
		return nil, &problem
	}

	// Restore body for downstream decoder
	r.Body = io.NopCloser(bytes.NewReader(body))

	// Minimal struct to peek data.type
	type peekDoc struct {
		Data struct {
			Type string `json:"type"`
		} `json:"data"`
	}

	peek, unmarshallErr := apix.UnmarshalSingleDocument[peekDoc](body)
	if unmarshallErr != nil {
		problem := apix.NewBadRequestProblem(fmt.Sprintf("failed to unmarshal JSON:API document: %s", unmarshallErr))
		return nil, &problem
	}

	msgType := peek.Data.Type
	if msgType == "" {
		problem := apix.NewBadRequestProblem("missing type in JSON:API document data")
		return nil, &problem
	}

	decoder, err := h.inner.decoderRegistry.Get(msgType, encoding)
	if err != nil {
		switch {
		case errors.Is(err, ErrMessageDecoderNotFoundForType):
			problem := apix.NewNotFoundProblem(fmt.Sprintf("no decoder registered for message type: %s", msgType))
			return nil, &problem
		case errors.Is(err, ErrMessageDecoderNotFoundForEncoding):
			problem := apix.NewUnsupportedMediaTypeProblem(fmt.Sprintf("no decoder registered for message type %q and encoding %q", msgType, encoding))
			return nil, &problem
		default:
			problem := apix.NewInternalServerErrorProblem(fmt.Sprintf("failed to get decoder for message type %q and encoding %q: %v", msgType, encoding, err))
			return nil, &problem
		}
	}

	// Ensure the decoder sees the same body
	r.Body = io.NopCloser(bytes.NewReader(body))

	msg, err := decoder(r)
	if err != nil {
		problem := apix.NewBadRequestProblem(fmt.Sprintf("failed to decode message %q: %v", msgType, err))
		return nil, &problem
	}
	return msg, nil
}
