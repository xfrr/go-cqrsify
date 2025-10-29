package messaginghttp_test

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	messaginghttp "github.com/xfrr/go-cqrsify/messaging/http"
	messagingmock "github.com/xfrr/go-cqrsify/messaging/mock"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/xfrr/go-cqrsify/messaging"
	"github.com/xfrr/go-cqrsify/pkg/apix"
)

type mockMessage struct {
	messaging.BaseMessage
}

func newJSONAPIRequest(t *testing.T, body []byte) *http.Request {
	t.Helper()
	req := httptest.NewRequest(http.MethodPost, "/messages", bytes.NewReader(body))
	req.Header.Set(apix.ContentTypeHeaderKey, apix.ContentTypeJSONAPI.String())
	return req
}

func readBody(t *testing.T, rr *httptest.ResponseRecorder) []byte {
	t.Helper()
	b, err := io.ReadAll(rr.Body)
	require.NoError(t, err)
	return b
}

type ServeHTTPSuite struct {
	suite.Suite
}

func TestServeHTTPSuite(t *testing.T) {
	suite.Run(t, new(ServeHTTPSuite))
}

func (s *ServeHTTPSuite) Test_NoPublisherConfigured_Returns500() {
	h := messaginghttp.NewMessageHandler(nil) // publisher is nil
	rr := httptest.NewRecorder()
	req := newJSONAPIRequest(s.T(), []byte(`{"data":{"type":"any"}}`))

	h.ServeHTTP(rr, req)

	s.Equal(http.StatusInternalServerError, rr.Code)
	s.NotEmpty(readBody(s.T(), rr))
}

func (s *ServeHTTPSuite) Test_NoDecoderRegistered_Returns400() {
	pub := &messagingmock.MessagePublisher{}
	h := messaginghttp.NewMessageHandler(pub)

	rr := httptest.NewRecorder()
	req := newJSONAPIRequest(s.T(), []byte(`{"data":{"type":"unknown","attributes":{}}}`))

	h.ServeHTTP(rr, req)

	s.Equal(http.StatusBadRequest, rr.Code, "no decoder for message type should yield 400")
	s.Zero(pub.PublishCalls(), "publisher must not be called")
}

func (s *ServeHTTPSuite) Test_InvalidContentTypeHeader_Returns415() {
	pub := &messagingmock.MessagePublisher{}
	h := messaginghttp.NewMessageHandler(pub)

	req := httptest.NewRequest(http.MethodPost, "/x", bytes.NewReader(nil))
	req.Header.Set(apix.ContentTypeHeaderKey, "not-a-valid-media-type")

	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)

	s.Equal(http.StatusUnsupportedMediaType, rr.Code)
}

func (s *ServeHTTPSuite) Test_UnsupportedContentType_Returns415() {
	pub := &messagingmock.MessagePublisher{}
	h := messaginghttp.NewMessageHandler(pub)

	req := httptest.NewRequest(http.MethodPost, "/x", bytes.NewReader(nil))
	req.Header.Set(apix.ContentTypeHeaderKey, "application/xml")

	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)

	s.Equal(http.StatusUnsupportedMediaType, rr.Code)
}

func (s *ServeHTTPSuite) Test_Success_Returns202_AndPublishes() {
	pub := &messagingmock.MessagePublisher{
		PublishFunc: func(_ context.Context, _ ...messaging.Message) error {
			return nil
		},
	}

	h := messaginghttp.NewMessageHandler(pub)

	err := messaginghttp.RegisterJSONAPIMessageDecoder(h, "ping",
		func(_ apix.SingleDocument[struct{}]) (messaging.Message, error) {
			return &mockMessage{
				BaseMessage: messaging.NewMessage("ping"),
			}, nil
		})
	s.Require().NoError(err)

	rr := httptest.NewRecorder()
	req := newJSONAPIRequest(s.T(), []byte(`{"data":{"type":"ping","attributes":{}}}`))

	h.ServeHTTP(rr, req)

	s.Equal(http.StatusAccepted, rr.Code)
	s.Len(pub.PublishCalls(), 1)
	s.NotNil(pub.PublishCalls()[0].Messages)
	s.Len(pub.PublishCalls()[0].Messages, 1)
	s.IsType(&mockMessage{}, pub.PublishCalls()[0].Messages[0])
	s.Equal("ping", pub.PublishCalls()[0].Messages[0].MessageType())
}

func (s *ServeHTTPSuite) Test_PublishError_DefaultsTo500() {
	pub := &messagingmock.MessagePublisher{
		PublishFunc: func(_ context.Context, _ ...messaging.Message) error {
			return errors.New("something went wrong")
		},
	}

	h := messaginghttp.NewMessageHandler(pub)

	err := messaginghttp.RegisterJSONAPIMessageDecoder(h, "ok",
		func(_ apix.SingleDocument[struct{}]) (messaging.Message, error) {
			return &mockMessage{}, nil
		})
	s.Require().NoError(err)

	rr := httptest.NewRecorder()
	req := newJSONAPIRequest(s.T(), []byte(`{"data":{"type":"ok","attributes":{}}}`))

	h.ServeHTTP(rr, req)

	s.Equal(http.StatusInternalServerError, rr.Code)
	s.Len(pub.PublishCalls(), 1)
}

func TestRegisterJSONAPIMessageDecoder_DuplicateReturnsError(t *testing.T) {
	pub := &messagingmock.MessagePublisher{}
	h := messaginghttp.NewMessageHandler(pub)

	require.NoError(t, messaginghttp.RegisterJSONAPIMessageDecoder(h, "dup",
		func(_ apix.SingleDocument[struct{}]) (messaging.Message, error) { return &mockMessage{}, nil }))

	err := messaginghttp.RegisterJSONAPIMessageDecoder(h, "dup",
		func(_ apix.SingleDocument[struct{}]) (messaging.Message, error) { return &mockMessage{}, nil })

	require.Error(t, err)
	assert.Contains(t, err.Error(), "already exists")
}

func TestServeHTTP_MissingType_Returns400(t *testing.T) {
	pub := &messagingmock.MessagePublisher{}
	h := messaginghttp.NewMessageHandler(pub)

	rr := httptest.NewRecorder()
	req := newJSONAPIRequest(t, []byte(`{"data":{"attributes":{"a":1}}}`)) // no data.type

	h.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Zero(t, pub.PublishCalls(), "publisher must not be called")
}
