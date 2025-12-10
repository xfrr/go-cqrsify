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

// Aliases to keep external API surface familiar.
type (
	HTTPMessageRequestValidator = apix.HTTPRequestValidator
)

type HTTPMessageEncoding string

const (
	HTTPMessageEncodingJSONAPI HTTPMessageEncoding = HTTPMessageEncoding(apix.ContentTypeJSONAPI)
	defaultMaxBodyBytes                            = int64(1 << 20) // 1 MiB sane default
)

// MessageHandler is an HTTP server for receiving messages.
type MessageHandler struct {
	// errorMapper maps handler errors to HTTP Problems.
	errorMapper func(error) apix.Problem

	// decoderRegistry: messageType -> encoding -> decode
	decoderRegistry *MessageDecoderRegistry

	// messagePublisher is the message bus used to publish messages.
	messagePublisher messaging.MessagePublisher

	// messageValidator validates incoming HTTP requests.
	messageValidator apix.HTTPRequestValidator

	// maxBodyBytes is the maximum allowed request body size in bytes.
	// If zero or negative, no limit is applied.
	maxBodyBytes int64
}

// NewMessageHandler creates a new MessageHandler with the given MessagePublisher and options.
func NewMessageHandler(msgPublisher messaging.MessagePublisher, opts ...HTTPMessageServerOption) *MessageHandler {
	s := &MessageHandler{
		messagePublisher: msgPublisher,
		maxBodyBytes:     defaultMaxBodyBytes,
		errorMapper:      nil,
		messageValidator: nil,
		decoderRegistry:  NewMessageDecoderRegistry(),
	}
	for _, o := range opts {
		o(s)
	}
	return s
}

// ServeHTTP implements http.Handler.
func (handler *MessageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if handler.messagePublisher == nil {
		apix.WriteProblem(w, apix.NewInternalServerErrorProblem("no message bus configured"))
		return
	}

	if handler.maxBodyBytes > 0 {
		r.Body = http.MaxBytesReader(w, r.Body, handler.maxBodyBytes)
	}
	defer r.Body.Close()

	if handler.messageValidator != nil {
		if problem := handler.messageValidator.Validate(r.Context(), r); problem != nil {
			apix.WriteProblem(w, *problem)
			return
		}
	}

	msg, problem := handler.decodeMessageFromHTTPRequest(r)
	if problem != nil {
		apix.WriteProblem(w, *problem)
		return
	}

	if err := handler.messagePublisher.Publish(r.Context(), msg); err != nil {
		handler.handleDispatchError(w, err)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

func (handler *MessageHandler) handleDispatchError(w http.ResponseWriter, err error) {
	if handler.errorMapper != nil {
		apix.WriteProblem(w, handler.errorMapper(err))
		return
	}
	apix.WriteProblem(w, apix.NewInternalServerErrorProblem(err.Error()))
}

func (handler *MessageHandler) decodeMessageFromHTTPRequest(r *http.Request) (messaging.Message, *apix.Problem) {
	if handler.decoderRegistry == nil {
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
		return handler.decodeJSONAPIMessage(r, HTTPMessageEncodingJSONAPI)
	default:
		problem := apix.NewUnsupportedMediaTypeProblem(fmt.Sprintf("unsupported content type: %s", mediaType))
		return nil, &problem
	}
}

func (handler *MessageHandler) decodeJSONAPIMessage(r *http.Request, encoding HTTPMessageEncoding) (messaging.Message, *apix.Problem) {
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

	decoder, err := handler.decoderRegistry.Get(msgType, encoding)
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

func makeMessageDecoder[P any](decodeFunc func(context.Context, apix.SingleDocument[P]) (messaging.Message, error)) func(*http.Request) (messaging.Message, error) {
	return func(r *http.Request) (messaging.Message, error) {
		defer r.Body.Close()

		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			return nil, err
		}

		doc, err := apix.UnmarshalSingleDocument[P](bodyBytes)
		if err != nil {
			return nil, err
		}

		return decodeFunc(r.Context(), doc)
	}
}

// RegisterJSONSingleDocumentMessageDecoder registers a JSON:API message decoder for the given message type.
// If a decoder for the same message type and encoding already exists, an error is returned.
func RegisterJSONSingleDocumentMessageDecoder[A any](handler *MessageHandler, msgType string, decodeFunc func(context.Context, apix.SingleDocument[A]) (messaging.Message, error)) error {
	return handler.decoderRegistry.Register(
		msgType,
		HTTPMessageEncodingJSONAPI,
		makeMessageDecoder(decodeFunc),
	)
}

// RegisterJSONManyDocumentMessageDecoder registers a JSON:API message decoder for the given message type that decodes ManyDocument.
// If a decoder for the same message type and encoding already exists, an error is returned.
func RegisterJSONManyDocumentMessageDecoder[A any](handler *MessageHandler, msgType string, decodeFunc func(apix.ManyDocument[A]) (messaging.Message, error)) error {
	return handler.decoderRegistry.Register(
		msgType,
		HTTPMessageEncodingJSONAPI,
		func(r *http.Request) (messaging.Message, error) {
			defer r.Body.Close()

			bodyBytes, err := io.ReadAll(r.Body)
			if err != nil {
				return nil, err
			}

			doc, err := apix.UnmarshalManyDocument[A](bodyBytes)
			if err != nil {
				return nil, err
			}

			return decodeFunc(doc)
		},
	)
}
