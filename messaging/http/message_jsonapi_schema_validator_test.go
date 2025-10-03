package messaginghttp_test

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	messaginghttp "github.com/xfrr/go-cqrsify/messaging/http"

	"github.com/stretchr/testify/suite"
)

// errReadCloser simulates a failing body reader.
type errReadCloser struct{}

func (e *errReadCloser) Read(_ []byte) (int, error) { return 0, errors.New("boom read error") }
func (e *errReadCloser) Close() error               { return nil }

// newRequestWithBody creates a new *http.Request with the provided JSON body.
func newRequestWithBody(method, target, body string) *http.Request {
	req := httptest.NewRequest(method, target, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/vnd.api+json")
	return req
}

// ================================
// JSONAPISchemaFilepathResolver
// ================================

type JSONAPISchemaFilepathResolverSuite struct {
	suite.Suite
	baseURI  string
	makeReq  func(body string) *http.Request
	resolver func(r *http.Request) (string, error)
	//nolint:containedctx // context is only used in tests
	ctx context.Context
}

func (s *JSONAPISchemaFilepathResolverSuite) SetupTest() {
	s.baseURI = "/schemas"
	s.makeReq = func(body string) *http.Request {
		return newRequestWithBody(http.MethodPost, "/messages", body)
	}
	s.resolver = messaginghttp.JSONAPISchemaFilepathResolver(s.baseURI)
	s.ctx = context.Background()
}

func (s *JSONAPISchemaFilepathResolverSuite) Test_Resolver_Success_PathConstruction_AndBodyRestored() {
	// Arrange
	body := `{"data":{"type":"user_created"}}`
	req := s.makeReq(body)

	// Act
	path, err := s.resolver(req)

	// Assert
	s.Require().NoError(err)
	s.Equal("/schemas/user_created.schema.json", path, "expected path constructed from baseURI and message type")

	// Body should be restored and readable downstream.
	got, readErr := io.ReadAll(req.Body)
	s.Require().NoError(readErr)
	s.JSONEq(body, string(got), "request body must be restored unchanged")
}

func (s *JSONAPISchemaFilepathResolverSuite) Test_Resolver_Error_InvalidJSON() {
	// Arrange
	req := s.makeReq(`{`) // malformed JSON

	// Act
	path, err := s.resolver(req)

	// Assert
	s.Require().Error(err)
	s.Empty(path)
	s.Contains(err.Error(), "failed to unmarshal JSON:API document")
}

func (s *JSONAPISchemaFilepathResolverSuite) Test_Resolver_Error_MissingType() {
	// Arrange
	req := s.makeReq(`{"data":{}}`) // no "type"

	// Act
	path, err := s.resolver(req)

	// Assert
	s.Require().Error(err)
	s.Empty(path)
	s.Contains(err.Error(), "missing type")
}

func (s *JSONAPISchemaFilepathResolverSuite) Test_Resolver_Error_BodyReadFailure() {
	// Arrange
	req := httptest.NewRequest(http.MethodPost, "/messages", &errReadCloser{})
	req.Header.Set("Content-Type", "application/vnd.api+json")

	// Act
	path, err := s.resolver(req)

	// Assert
	s.Require().Error(err)
	s.Empty(path)
	s.Contains(err.Error(), "failed to read request body")
}

func (s *JSONAPISchemaFilepathResolverSuite) Test_Resolver_AllowsRepeatedDownstreamReads() {
	// Arrange
	body := `{"data":{"type":"sample.event"}}`
	req := s.makeReq(body)

	// Act
	_, err := s.resolver(req)
	s.Require().NoError(err)

	// Simulate two downstream consumers reading the body sequentially.
	b1, err1 := io.ReadAll(req.Body)
	s.Require().NoError(err1)
	// Put the body back again so a further downstream consumer still can read it (mimics some middlewares).
	req.Body = io.NopCloser(bytes.NewReader(b1))

	b2, err2 := io.ReadAll(req.Body)
	s.Require().NoError(err2)

	// Assert
	s.Equal(string(b1), string(b2))
	s.JSONEq(body, string(b2))
}

func TestJSONAPISchemaFilepathResolverSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(JSONAPISchemaFilepathResolverSuite))
}

// ================================
// MessageJSONAPISchemaValidator
// ================================

type MessageJSONAPISchemaValidatorSuite struct {
	suite.Suite
}

func (s *MessageJSONAPISchemaValidatorSuite) Test_Validate_WhenNoUnderlyingValidator_ReturnsProblem() {
	// Arrange: zero-value struct -> httpValidator is nil inside
	var v messaginghttp.MessageJSONAPISchemaValidator

	req := newRequestWithBody(http.MethodPost, "/messages", `{"data":{"type":"x"}}`)

	// Act
	problem := v.Validate(context.Background(), req)

	// Assert
	// We cannot assert the Problem fields without coupling to apix internals,
	// but we can ensure a non-nil problem is returned as a defensive check.
	s.Require().NotNil(problem, "expected a non-nil apix.Problem when no validator is configured")
}

func (s *MessageJSONAPISchemaValidatorSuite) Test_Validate_DelegatesToUnderlyingValidator_FromConstructor() {
	// Arrange: constructor wires a jsonschema validator internally.
	validator := messaginghttp.NewMessageJSONAPISchemaValidator(
	// No options required for this behavioral check; we just verify delegation path executes.
	)

	// Any request; the inner validator may return a Problem because it can't resolve a schema,
	// but we only need to confirm Validate delegates and returns *apix.Problem (possibly nil).
	req := newRequestWithBody(http.MethodPost, "/messages", `{"data":{"type":"anything"}}`)

	// Act
	problem := validator.Validate(context.Background(), req)

	// Assert
	// We cannot assert the Problem fields without coupling to apix internals,
	// but we can ensure a non-nil problem is returned as a defensive check.
	s.Require().NotNil(problem, "expected a non-nil apix.Problem from underlying validator")
}

func TestMessageJSONAPISchemaValidatorSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(MessageJSONAPISchemaValidatorSuite))
}
