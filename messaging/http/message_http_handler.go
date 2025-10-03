package messaginghttp

import (
	"bytes"
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
	validator   apix.HTTPRequestValidator
	errorMapper func(error) apix.Problem

	// decoders: messageType -> encoding -> decode
	decoders map[string]map[HTTPMessageEncoding]func(*http.Request) (messaging.Message, error)

	messagePublisher messaging.MessagePublisher
	allowEncodings   map[HTTPMessageEncoding]struct{}

	// maxBodyBytes is the maximum allowed request body size in bytes.
	// If zero or negative, no limit is applied.
	maxBodyBytes int64
}

// NewMessageHTTPHandler creates a new HTTPMessageServer with the given MessageBus and options.
// If no decoders are registered, the server will return 500 Internal Server Error.
func NewMessageHTTPHandler(msgPublisher messaging.MessagePublisher, opts ...HTTPMessageServerOption) *MessageHandler {
	s := &MessageHandler{
		messagePublisher: msgPublisher,
		maxBodyBytes:     defaultMaxBodyBytes,
		errorMapper:      nil,
		validator:        nil,
		decoders:         make(map[string]map[HTTPMessageEncoding]func(*http.Request) (messaging.Message, error)),
		allowEncodings:   map[HTTPMessageEncoding]struct{}{HTTPMessageEncodingJSONAPI: {}},
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

	if handler.validator != nil {
		if problem := handler.validator.Validate(r.Context(), r); problem != nil {
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
	if handler.decoders == nil {
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

	msgDecoders, ok := handler.decoders[msgType]
	if !ok {
		problem := apix.NewBadRequestProblem(fmt.Sprintf("no decoder registered for message type: %s", msgType))
		return nil, &problem
	}

	decodeFunc, ok := msgDecoders[encoding]
	if !ok {
		problem := apix.NewUnsupportedMediaTypeProblem(fmt.Sprintf("no decoder registered for message type %q and encoding %q", msgType, encoding))
		return nil, &problem
	}

	// Ensure the decoder sees the same body
	r.Body = io.NopCloser(bytes.NewReader(body))

	msg, err := decodeFunc(r)
	if err != nil {
		problem := apix.NewBadRequestProblem(fmt.Sprintf("failed to decode message %q: %v", msgType, err))
		return nil, &problem
	}
	return msg, nil
}

func makeMessageDecoder[P any](decodeFunc func(apix.SingleDocument[P]) (messaging.Message, error)) func(*http.Request) (messaging.Message, error) {
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

		return decodeFunc(doc)
	}
}

// RegisterJSONAPIMessageDecoder registers a JSON:API message decoder for the given message type.
// If a decoder for the same message type and encoding already exists, an error is returned.
func RegisterJSONAPIMessageDecoder[A any](handler *MessageHandler, msgType string, decodeFunc func(apix.SingleDocument[A]) (messaging.Message, error)) error {
	if handler.decoders == nil {
		handler.decoders = make(map[string]map[HTTPMessageEncoding]func(*http.Request) (messaging.Message, error))
	}

	msgDecoders, ok := handler.decoders[msgType]
	if !ok {
		msgDecoders = make(map[HTTPMessageEncoding]func(*http.Request) (messaging.Message, error))
		handler.decoders[msgType] = msgDecoders
	}

	if _, exists := msgDecoders[HTTPMessageEncodingJSONAPI]; exists {
		return fmt.Errorf("message decoder for %q and encoding %q already exists", msgType, HTTPMessageEncodingJSONAPI)
	}

	msgDecoders[HTTPMessageEncodingJSONAPI] = makeMessageDecoder[A](decodeFunc)
	return nil
}

// RegisterJSONAPICommandDecoder registers a JSON:API command decoder for the given command type.
// If a decoder for the same command type and encoding already exists, an error is returned.
func RegisterJSONAPICommandDecoder[A any](handler *MessageHandler, msgType string, decodeFunc func(apix.SingleDocument[A]) (messaging.Command, error)) error {
	if handler.decoders == nil {
		handler.decoders = make(map[string]map[HTTPMessageEncoding]func(*http.Request) (messaging.Message, error))
	}

	msgDecoders, ok := handler.decoders[msgType]
	if !ok {
		msgDecoders = make(map[HTTPMessageEncoding]func(*http.Request) (messaging.Message, error))
		handler.decoders[msgType] = msgDecoders
	}
	if _, exists := msgDecoders[HTTPMessageEncodingJSONAPI]; exists {
		return fmt.Errorf("command decoder for %q and encoding %q already exists", msgType, HTTPMessageEncodingJSONAPI)
	}

	msgDecoders[HTTPMessageEncodingJSONAPI] = makeMessageDecoder[A](func(sd apix.SingleDocument[A]) (messaging.Message, error) {
		return decodeFunc(sd)
	})
	return nil
}
