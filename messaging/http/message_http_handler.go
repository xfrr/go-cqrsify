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

// HTTPMessageServer is an HTTP server for receiving messages.
type HTTPMessageServer struct {
	validator   apix.HTTPRequestValidator
	errorMapper func(error) apix.Problem

	// decoders: messageType -> encoding -> decode
	decoders map[string]map[HTTPMessageEncoding]func(*http.Request) (messaging.Message, error)

	messageBus     messaging.MessageBus
	allowEncodings map[HTTPMessageEncoding]struct{}

	// maxBodyBytes is the maximum allowed request body size in bytes.
	// If zero or negative, no limit is applied.
	maxBodyBytes int64
}

// NewMessageHTTPHandler creates a new HTTPMessageServer with the given MessageBus and options.
// If no decoders are registered, the server will return 500 Internal Server Error.
func NewMessageHTTPHandler(messageBus messaging.MessageBus, opts ...HTTPMessageServerOption) *HTTPMessageServer {
	s := &HTTPMessageServer{
		messageBus:     messageBus,
		maxBodyBytes:   defaultMaxBodyBytes,
		errorMapper:    nil,
		validator:      nil,
		decoders:       make(map[string]map[HTTPMessageEncoding]func(*http.Request) (messaging.Message, error)),
		allowEncodings: map[HTTPMessageEncoding]struct{}{HTTPMessageEncodingJSONAPI: {}},
	}
	for _, o := range opts {
		o(s)
	}
	return s
}

// ServeHTTP implements http.Handler.
func (s *HTTPMessageServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if s.messageBus == nil {
		apix.WriteProblem(w, apix.NewInternalServerErrorProblem("no message bus configured"))
		return
	}

	if s.maxBodyBytes > 0 {
		r.Body = http.MaxBytesReader(w, r.Body, s.maxBodyBytes)
	}
	defer r.Body.Close()

	if s.validator != nil {
		// Validate request (body, headers, method, etc.)
		if problem := s.validator.Validate(r.Context(), r); problem != nil {
			apix.WriteProblem(w, *problem)
			return
		}
	}

	msg, problem := s.decodeMessageFromHTTPRequest(r)
	if problem != nil {
		apix.WriteProblem(w, *problem)
		return
	}

	if err := s.messageBus.Publish(r.Context(), msg); err != nil {
		s.handleDispatchError(w, err)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

func (s *HTTPMessageServer) handleDispatchError(w http.ResponseWriter, err error) {
	if s.errorMapper != nil {
		apix.WriteProblem(w, s.errorMapper(err))
		return
	}
	apix.WriteProblem(w, apix.NewInternalServerErrorProblem(err.Error()))
}

func (s *HTTPMessageServer) decodeMessageFromHTTPRequest(r *http.Request) (messaging.Message, *apix.Problem) {
	if s.decoders == nil {
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
		return s.decodeJSONAPIMessage(r, HTTPMessageEncodingJSONAPI)
	default:
		problem := apix.NewUnsupportedMediaTypeProblem(fmt.Sprintf("unsupported content type: %s", mediaType))
		return nil, &problem
	}
}

func (s *HTTPMessageServer) decodeJSONAPIMessage(r *http.Request, encoding HTTPMessageEncoding) (messaging.Message, *apix.Problem) {
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

	msgDecoders, ok := s.decoders[msgType]
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

// RegisterJSONAPIMessageDecoder registers a decoder for a given message type (method form).
func (s *HTTPMessageServer) RegisterJSONAPIMessageDecoder(msgType string, decodeFunc func(apix.SingleDocument[any]) (messaging.Message, error)) error {
	if s.decoders == nil {
		s.decoders = make(map[string]map[HTTPMessageEncoding]func(*http.Request) (messaging.Message, error))
	}
	msgDecoders, ok := s.decoders[msgType]
	if !ok {
		msgDecoders = make(map[HTTPMessageEncoding]func(*http.Request) (messaging.Message, error))
		s.decoders[msgType] = msgDecoders
	}
	if _, exists := msgDecoders[HTTPMessageEncodingJSONAPI]; exists {
		return fmt.Errorf("message decoder for %q and encoding %q already exists", msgType, HTTPMessageEncodingJSONAPI)
	}
	msgDecoders[HTTPMessageEncodingJSONAPI] = makeMessageDecoder(decodeFunc)
	return nil
}

// Backwards-compatible free function (optional; can be removed if not needed).
func RegisterJSONAPIMessageDecoder[P any](server *HTTPMessageServer, msgType string, decodeFunc func(apix.SingleDocument[P]) (messaging.Message, error)) error {
	return server.RegisterJSONAPIMessageDecoder(msgType, func(sd apix.SingleDocument[any]) (messaging.Message, error) {
		var converted apix.SingleDocument[P]
		attr, ok := sd.Data.Attributes.(P)
		if !ok {
			return nil, fmt.Errorf("failed to convert attributes to %T: %v", converted.Data, sd.Data.Attributes)
		}

		converted.Included = sd.Included
		converted.Meta = sd.Meta
		converted.Links = sd.Links
		converted.Data = apix.Resource[P]{
			Type:          sd.Data.Type,
			ID:            sd.Data.ID,
			Attributes:    attr,
			Relationships: sd.Data.Relationships,
			Meta:          sd.Data.Meta,
		}

		return decodeFunc(converted)
	})
}
