package apix_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	apix "github.com/xfrr/go-cqrsify/pkg/apix"
)

type WriteJSONSuite struct {
	suite.Suite
}

type payload struct {
	A int    `json:"a"`
	B string `json:"b"`
}

func (s *WriteJSONSuite) mustMarshal(v any) []byte {
	b, err := json.Marshal(v)
	s.Require().NoError(err)
	return b
}

func (s *WriteJSONSuite) Test_StrongETag_DefaultAndBody_ContentType_Status() {
	rec := httptest.NewRecorder()
	p := payload{A: 1, B: "z"}

	// choose explicit content-type and status via options
	apix.WriteJSON(
		rec,
		p,
		apix.WithContentType(apix.ContentTypeJSON),
		apix.WithStatus(http.StatusOK),
	)

	// Status and Content-Type
	s.Equal(http.StatusOK, rec.Code)
	s.Equal("application/json", rec.Header().Get(apix.ContentTypeHeaderKey))

	// Body is JSON encoding of payload
	var got payload
	s.Require().NoError(json.Unmarshal(rec.Body.Bytes(), &got))
	s.Equal(p, got)

	// Strong ETag equals hash of encoded body
	expected := apix.StrongETagFromBytes(rec.Body.Bytes())
	s.Equal(expected, rec.Header().Get(apix.ETagHeaderKey))
}

func (s *WriteJSONSuite) Test_WeakETag_HeaderSet_WhenOptionEnabled() {
	rec := httptest.NewRecorder()
	p := payload{A: 7, B: "foo"}

	apix.WriteJSON(
		rec,
		p,
		apix.WithContentType(apix.ContentTypeJSON),
		apix.WithStatus(http.StatusAccepted),
		apix.WithWeakETag(true),
	)

	s.Equal(http.StatusAccepted, rec.Code)
	s.Equal("application/json", rec.Header().Get(apix.ContentTypeHeaderKey))

	// Weak ETag computed from body
	expected := apix.WeakETagFromBytes(rec.Body.Bytes())
	s.Equal(expected, rec.Header().Get(apix.ETagHeaderKey))
}

func (s *WriteJSONSuite) Test_ExplicitETag_Overrides_ComputedETags() {
	rec := httptest.NewRecorder()
	p := payload{A: 3, B: "bar"}

	apix.WriteJSON(
		rec,
		p,
		apix.WithContentType(apix.ContentTypeJSON),
		apix.WithStatus(http.StatusOK),
		apix.WithETag(`"my-custom-etag"`),
	)

	s.Equal(`"my-custom-etag"`, rec.Header().Get(apix.ETagHeaderKey))
}

func (s *WriteJSONSuite) Test_Cache_LastModified_Vary_IfModifiedSince_CustomHeaders() {
	rec := httptest.NewRecorder()
	p := payload{A: 5, B: "headers"}

	// Use times truncated to seconds to match http.TimeFormat precision
	lm := time.Now().UTC().Truncate(time.Second)
	ims := lm.Add(-time.Hour)

	apix.WriteJSON(
		rec,
		p,
		apix.WithContentType(apix.ContentTypeJSON),
		apix.WithStatus(http.StatusOK),
		apix.WithCacheControl("public, max-age=60"),
		apix.WithLastModified(lm),
		apix.WithVary("Accept, Accept-Language"),
		apix.WithIfModifiedSince(ims),
		apix.WithHeaders(map[string]string{
			"X-Custom": "abc",
			"Server":   "unit-test",
		}),
	)

	s.Equal("public, max-age=60", rec.Header().Get(apix.CacheControlHeaderKey))
	s.Equal(lm.Format(http.TimeFormat), rec.Header().Get(apix.LastModifiedHeaderKey))
	s.Equal("Accept, Accept-Language", rec.Header().Get(apix.VaryHeaderKey))
	s.Equal(ims.Format(http.TimeFormat), rec.Header().Get(apix.IfModifiedSinceHeaderKey))
	s.Equal("abc", rec.Header().Get("X-Custom"))
	s.Equal("unit-test", rec.Header().Get("Server"))
}

func (s *WriteJSONSuite) Test_EncodingError_WritesRFC7807_500() {
	// json.Marshal should fail on a channel -> WriteProblem(500)
	rec := httptest.NewRecorder()

	apix.WriteJSON(
		rec,
		make(chan int),
		apix.WithContentType(apix.ContentTypeJSON),
	)

	s.Equal(http.StatusInternalServerError, rec.Code)
	s.Equal(apix.ContentTypeProblemJSON.String(), rec.Header().Get(apix.ContentTypeHeaderKey))

	var prob map[string]any
	s.Require().NoError(json.Unmarshal(rec.Body.Bytes(), &prob))
	s.Equal("Encoding error", prob["title"])
	s.InDelta(float64(http.StatusInternalServerError), prob["status"], 0)
}

func TestWriteJSONSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(WriteJSONSuite))
}
