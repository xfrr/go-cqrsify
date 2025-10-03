package messaginghttp_test

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/xfrr/go-cqrsify/messaging"
	"github.com/xfrr/go-cqrsify/pkg/apix"

	messaginghttp "github.com/xfrr/go-cqrsify/messaging/http"
)

type fakeMessageHTTPValidator struct {
	problem *apix.Problem
}

func (v *fakeMessageHTTPValidator) Validate(_ context.Context, _ *http.Request) *apix.Problem {
	return v.problem
}

func makeJSONAPIMessageBody(typ string, attrs string) []byte {
	if attrs == "" {
		attrs = `"name":"test"`
	}
	body := fmt.Sprintf(`{"data":{"type":"%s","attributes":{%s}}}`, typ, attrs)
	return []byte(body)
}

func makeJSONAPIMessageRequest(t *testing.T, body []byte) *http.Request {
	t.Helper()
	req := httptest.NewRequest(http.MethodPost, "/messages", bytes.NewReader(body))
	req.Header.Set(apix.ContentTypeHeaderKey, apix.ContentTypeJSONAPI.String())
	return req
}

func recordHTTPResponse(handler http.Handler, r *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, r)
	return rr
}

type countingDecoder struct {
	count int
	err   error
	msg   messaging.Message
}

func (d *countingDecoder) fn(_ *http.Request) (messaging.Message, error) {
	d.count++
	if d.err != nil {
		return nil, d.err
	}

	msg := d.msg
	if msg == nil {
		msg = messaging.NewBaseMessage("counting-msg")
	}
	return msg, nil
}

func TestHTTPMessageServerSuite(t *testing.T) {
	suite.Run(t, &HTTPMessageServerSuite{})
}

type HTTPMessageServerSuite struct {
	suite.Suite

	cmdbus *messaging.InMemoryMessageBus
}

func (st *HTTPMessageServerSuite) SetupTest() {
	st.cmdbus = messaging.NewInMemoryMessageBus()
}

func (st *HTTPMessageServerSuite) AfterTest(_, _ string) {
	if st.cmdbus != nil {
		st.cmdbus.Close()
		st.cmdbus = nil
	}
}

func (st *HTTPMessageServerSuite) Test_ServeHTTP_NoBusConfigured_Returns500() {
	s := messaginghttp.NewMessageHTTPHandler(nil)

	req := makeJSONAPIMessageRequest(st.T(), makeJSONAPIMessageBody("X", ""))

	rr := recordHTTPResponse(s, req)
	st.Require().Equal(http.StatusInternalServerError, rr.Code)
	st.Contains(rr.Body.String(), "no message bus configured")
}

func (st *HTTPMessageServerSuite) Test_ServeHTTP_UnsupportedContentType_Returns415() {
	srv := messaginghttp.NewMessageHTTPHandler(st.cmdbus)

	req := httptest.NewRequest(http.MethodPost, "/messages", strings.NewReader(`{}`))
	req.Header.Set(apix.ContentTypeHeaderKey, "text/plain; charset=utf-8")

	rr := recordHTTPResponse(srv, req)
	st.Require().Contains(rr.Body.String(), "unsupported content type")
	st.Require().Equal(http.StatusUnsupportedMediaType, rr.Code)
}

func (st *HTTPMessageServerSuite) Test_JSONAPI_MissingType_Returns400() {
	srv := messaginghttp.NewMessageHTTPHandler(st.cmdbus)

	req := makeJSONAPIMessageRequest(st.T(), []byte(`{"data":{"attributes":{"x":1}}}`))
	rr := recordHTTPResponse(srv, req)
	st.Require().Equal(http.StatusBadRequest, rr.Code)
	st.Contains(rr.Body.String(), "missing type")
}

func (st *HTTPMessageServerSuite) Test_JSONAPI_UnknownMessageType_Returns400() {
	srv := messaginghttp.NewMessageHTTPHandler(st.cmdbus)

	req := makeJSONAPIMessageRequest(st.T(), makeJSONAPIMessageBody("unknown-msg", `"x":1`))
	rr := recordHTTPResponse(srv, req)
	st.Require().Equal(http.StatusBadRequest, rr.Code)
	st.Contains(rr.Body.String(), "no decoder registered for message type")
}

func (st *HTTPMessageServerSuite) Test_JSONAPI_NoDecoderForEncoding_Returns415() {
	s := messaginghttp.NewMessageHTTPHandler(st.cmdbus)

	// Register a decoder under a different message type, so lookup for "createUser" fails.
	err := messaginghttp.RegisterJSONAPIMessageDecoder(s, "otherType", func(_ apix.SingleDocument[any]) (messaging.Message, error) {
		fakeMsg := struct {
			messaging.BaseMessage
			X int `json:"x"`
		}{
			BaseMessage: messaging.NewBaseMessage("other.Type"),
			X:           1,
		}
		return &fakeMsg, nil
	})
	st.Require().NoError(err)

	req := makeJSONAPIMessageRequest(st.T(), makeJSONAPIMessageBody("createUser", `"x":1`))
	rr := recordHTTPResponse(s, req)
	st.Require().Equal(http.StatusBadRequest, rr.Code)
	st.Contains(rr.Body.String(), "no decoder registered for message type")
}

func (st *HTTPMessageServerSuite) Test_JSONAPI_DecoderError_Returns400() {
	sut := messaginghttp.NewMessageHTTPHandler(st.cmdbus)

	dec := &countingDecoder{err: errors.New("boom")}

	// Register via method; we can't inject the inner function directly, so register a wrapper
	err := messaginghttp.RegisterJSONAPIMessageDecoder(sut, "createUser", func(_ apix.SingleDocument[map[string]any]) (messaging.Message, error) {
		// mimic the failing decoder
		return dec.fn(nil)
	})
	st.Require().NoError(err)

	req := makeJSONAPIMessageRequest(st.T(), makeJSONAPIMessageBody("createUser", `"x":1`))
	rr := recordHTTPResponse(sut, req)
	st.Require().Equal(http.StatusBadRequest, rr.Code)
	st.Contains(rr.Body.String(), "failed to decode message")
	st.Equal(1, dec.count)
}

func (st *HTTPMessageServerSuite) Test_JSONAPI_HappyPath_Returns202() {
	type fakeCmd struct {
		messaging.BaseMessage
		X int `json:"x"`
	}

	messageHandlerCalls := 0
	unsub, err := messaging.SubscribeMessage(st.T().Context(), st.cmdbus, "createUser", messaging.MessageHandlerFn[fakeCmd](func(_ context.Context, msg fakeCmd) error {
		// basic assertion: message was decoded properly
		st.Equal("createUser", msg.MessageType())
		st.Equal(1, msg.X)
		messageHandlerCalls++
		return nil
	}))
	st.Require().NoError(err)
	defer unsub()

	srv := messaginghttp.NewMessageHTTPHandler(st.cmdbus)

	decoderCalls := 0
	decodeErr := messaginghttp.RegisterJSONAPIMessageDecoder(srv, "createUser", func(doc apix.SingleDocument[map[string]any]) (messaging.Message, error) {
		decoderCalls++
		// basic assertion: attributes were present
		st.Require().NotNil(doc.Data)
		st.Require().NotNil(doc.Data.Attributes)
		xVal, ok := doc.Data.Attributes["x"]
		st.Require().True(ok)
		xFloat, ok := xVal.(float64) // JSON numbers are float64
		st.Require().True(ok)
		st.Require().InEpsilon(1.0, xFloat, 0.0001)

		msg := fakeCmd{
			BaseMessage: messaging.NewBaseMessage("createUser"),
			X:           int(xFloat),
		}
		return msg, nil
	})
	st.Require().NoError(decodeErr)

	req := makeJSONAPIMessageRequest(st.T(), makeJSONAPIMessageBody("createUser", `"x":1`))
	rr := recordHTTPResponse(srv, req)
	st.Require().Empty(rr.Body.String()) // no content
	st.Require().Equal(http.StatusAccepted, rr.Code)
	st.Equal(1, decoderCalls)
	st.Equal(1, messageHandlerCalls)
}

func (st *HTTPMessageServerSuite) Test_BodyTooLarge_Returns400() {
	srv := messaginghttp.NewMessageHTTPHandler(
		st.cmdbus,
		messaginghttp.WithMaxBodyBytes(10),
		messaginghttp.WithErrorMapper(func(err error) apix.Problem {
			if errors.Is(err, http.ErrContentLength) {
				return apix.NewBadRequestProblem("request body too large")
			}
			return apix.NewInternalServerErrorProblem("internal server error")
		}),
	)

	// Body larger than 10 bytes
	req := makeJSONAPIMessageRequest(st.T(), makeJSONAPIMessageBody("createUser", `"x":"1234567890ABCDEF"`))
	rr := recordHTTPResponse(srv, req)

	// The server reports a BadRequest because reading the body fails via MaxBytesReader.
	st.Require().Equal(http.StatusBadRequest, rr.Code)
	st.Contains(rr.Body.String(), "failed to read request body")
}

func (st *HTTPMessageServerSuite) Test_DuplicateRegistration_ReturnsError() {
	srv := messaginghttp.NewMessageHTTPHandler(st.cmdbus)
	err := messaginghttp.RegisterJSONAPIMessageDecoder(srv, "createUser", func(_ apix.SingleDocument[map[string]any]) (messaging.Message, error) {
		return messaging.NewBaseMessage("createUser"), nil
	})
	st.Require().NoError(err)

	err = messaginghttp.RegisterJSONAPIMessageDecoder(srv, "createUser", func(_ apix.SingleDocument[map[string]any]) (messaging.Message, error) {
		return messaging.NewBaseMessage("createUser"), nil
	})
	st.Require().Error(err)
	st.Contains(err.Error(), "already exists")
}

func (st *HTTPMessageServerSuite) Test_ValidatorProblem_Returned() {
	v := &fakeMessageHTTPValidator{problem: ptr(apix.NewBadRequestProblem("bad headers"))}

	srv := messaginghttp.NewMessageHTTPHandler(st.cmdbus, messaginghttp.WithValidator(v))

	req := makeJSONAPIMessageRequest(st.T(), makeJSONAPIMessageBody("createUser", `"x":1`))
	rr := recordHTTPResponse(srv, req)
	st.Require().Equal(http.StatusBadRequest, rr.Code)
	st.Contains(rr.Body.String(), "bad headers")
}

// helper
func ptr[T any](v T) *T { return &v }
