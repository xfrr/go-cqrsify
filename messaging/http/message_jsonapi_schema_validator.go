package messaginghttp

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/xfrr/go-cqrsify/pkg/apix"
	"github.com/xfrr/go-cqrsify/pkg/apix/jsonschema"
)

var _ apix.HTTPRequestValidator = (*MessageJSONAPISchemaValidator)(nil)

// MessageJSONAPISchemaValidator is an HTTP request validator
// that uses JSON Schema to validate message requests.
type MessageJSONAPISchemaValidator struct {
	httpValidator apix.HTTPRequestValidator
}

// NewMessageJSONAPISchemaValidator creates a new MessageSchemaValidator
// with the given HTTP request validator.
//
// The baseURI is used to construct the schema file path
// based on the message type found in the request body.
func NewMessageJSONAPISchemaValidator(options ...jsonschema.ValidatorOption) *MessageJSONAPISchemaValidator {
	return &MessageJSONAPISchemaValidator{
		httpValidator: jsonschema.NewValidator(options...),
	}
}

// Validate implements apix.HTTPRequestValidator.
func (v *MessageJSONAPISchemaValidator) Validate(ctx context.Context, r *http.Request) *apix.Problem {
	if v.httpValidator == nil {
		problem := apix.NewInternalServerErrorProblem("no HTTP request validator configured")
		return &problem
	}
	return v.httpValidator.Validate(ctx, r)
}

// JSONAPISchemaFilepathResolver returns a FilepathResolver
// that constructs the JSON Schema file path based on the message type
// found in the JSON:API request body.
//
// The baseURI is used as the base path for the schema files.
// For example, if baseURI is "/schemas" and the message type is "user.created",
// the resulting schema file path will be "/schemas/user.created.schema.json".
//
// Note: This resolver reads the entire request body (up to a reasonable limit)
// to extract the message type. It restores the body for downstream processing.
func JSONAPISchemaFilepathResolver(baseURI string) func(r *http.Request) (string, error) {
	return func(r *http.Request) (string, error) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			return "", fmt.Errorf("failed to read request body: %w", err)
		}

		// Restore body for downstream decoder
		r.Body = io.NopCloser(bytes.NewReader(body))

		type peekDoc struct {
			Data struct {
				Type string `json:"type"`
			} `json:"data"`
		}

		peek, unmarshallErr := apix.UnmarshalSingleDocument[peekDoc](body)
		if unmarshallErr != nil {
			return "", fmt.Errorf("failed to unmarshal JSON:API document: %w", unmarshallErr)
		}

		msgType := peek.Data.Type
		if msgType == "" {
			return "", errors.New("missing type in JSON:API document data")
		}

		path := baseURI + "/" + msgType + ".schema.json"
		return path, nil
	}
}
