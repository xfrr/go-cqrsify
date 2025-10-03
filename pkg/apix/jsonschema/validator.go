package jsonschema

import (
	"context"
	"errors"
	"net/http"

	"github.com/xeipuuv/gojsonschema"

	apihttp "github.com/xfrr/go-cqrsify/pkg/apix"
)

var _ apihttp.HTTPRequestValidator = (*HTTPRequestValidator)(nil)

// HTTPRequestValidator allows validating http request against a JSON Schema.
type HTTPRequestValidator struct {
	filepathResolver FilepathResolver
	problemURL       string
}

// FilepathResolver is a function that retrieves the JSON Schema file path
// based on the incoming HTTP request.
type FilepathResolver func(r *http.Request) (string, error)

// NewValidator creates a new JSON Schema validator.
func NewValidator(options ...ValidatorOption) *HTTPRequestValidator {
	const defaultProblemURL = "https://example.com/problems"
	var defaultFilepathResolver FilepathResolver = func(r *http.Request) (string, error) {
		return "./schemas" + r.URL.Path + ".schema.json", nil
	}

	v := &HTTPRequestValidator{
		problemURL:       defaultProblemURL,
		filepathResolver: defaultFilepathResolver,
	}

	for _, opt := range options {
		opt(v)
	}
	return v
}

// Validate implements apihttp.Validator.
func (v *HTTPRequestValidator) Validate(_ context.Context, r *http.Request) *apihttp.Problem {
	schemaPath, err := v.filepathResolver(r)
	if err != nil {
		return v.failedToLoadSchema(err)
	}

	schema := gojsonschema.NewReferenceLoader(schemaPath)
	if schema == nil {
		return v.failedToLoadSchema(nil)
	}

	// Validate the request against the JSON Schema
	result, err := gojsonschema.Validate(schema, gojsonschema.NewGoLoader(r))
	if err != nil {
		return v.failedToValidateRequest(err)
	}
	if !result.Valid() {
		return v.failedToValidateRequest(errors.New("request does not match schema"))
	}
	return nil
}

func (v *HTTPRequestValidator) failedToLoadSchema(err error) *apihttp.Problem {
	return &apihttp.Problem{
		Type:   v.problemURL + "/schema-load-failure",
		Title:  "Failed to load JSON Schema",
		Status: http.StatusInternalServerError,
		Detail: err.Error(),
	}
}

func (v *HTTPRequestValidator) failedToValidateRequest(err error) *apihttp.Problem {
	return &apihttp.Problem{
		Type:   v.problemURL + "/validation-failure",
		Title:  "Request validation failed",
		Status: http.StatusBadRequest,
		Detail: err.Error(),
	}
}
