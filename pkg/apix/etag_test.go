package apix_test

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/suite"
	apix "github.com/xfrr/go-cqrsify/pkg/apix"
)

type etagSuite struct {
	suite.Suite
}

func (s *etagSuite) Test_StrongETagFromBytes_FormatAndValue() {
	data := []byte("hello")
	sum := sha256.Sum256(data)
	expected := `"sha256:` + hex.EncodeToString(sum[:]) + `"`

	s.Equal(expected, apix.StrongETagFromBytes(data))
}

func (s *etagSuite) Test_WeakETagFromBytes_FormatAndValue() {
	data := []byte("hello")
	sum := sha256.Sum256(data)
	expected := `W/"sha256:` + hex.EncodeToString(sum[:]) + `"`

	s.Equal(expected, apix.WeakETagFromBytes(data))
}

func TestETagSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(etagSuite))
}

type preEncodeWriteJSONSuite struct {
	suite.Suite
}

type sample struct {
	A int    `json:"a"`
	B string `json:"b"`
}

func (s *preEncodeWriteJSONSuite) mustMarshal(v any) []byte {
	b, err := json.Marshal(v)
	s.Require().NoError(err)
	return b
}

func (s *preEncodeWriteJSONSuite) Test_WritesBodyHeadersStatus_StrongETag() {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/x", nil)
	payload := sample{A: 1, B: "z"}

	apix.PreEncodeAndWriteJSON(rec, req, http.StatusOK, "application/json", payload, false)

	s.Equal(http.StatusOK, rec.Code)
	s.Equal("application/json", rec.Header().Get("Content-Type"))

	// Check ETag equals strong hash of encoded body
	body := rec.Body.Bytes()
	sum := sha256.Sum256(body)
	expectedETag := `"sha256:` + hex.EncodeToString(sum[:]) + `"`
	s.Equal(expectedETag, rec.Header().Get("ETag"))

	// Body should be the JSON encoding of payload
	var got sample
	s.Require().NoError(json.Unmarshal(body, &got))
	s.Equal(payload, got)
}

func (s *preEncodeWriteJSONSuite) Test_WritesBodyHeadersStatus_WeakETag() {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/x", nil)
	payload := sample{A: 7, B: "foo"}

	// Precompute what the body will be to validate weak etag
	body := s.mustMarshal(payload)
	sum := sha256.Sum256(body)
	expectedETag := `W/"sha256:` + hex.EncodeToString(sum[:]) + `"`

	apix.PreEncodeAndWriteJSON(rec, req, http.StatusAccepted, "application/json", payload, true)

	s.Equal(http.StatusAccepted, rec.Code)
	s.Equal("application/json", rec.Header().Get("Content-Type"))
	s.Equal(expectedETag, rec.Header().Get("ETag"))

	var got sample
	s.Require().NoError(json.Unmarshal(rec.Body.Bytes(), &got))
	s.Equal(payload, got)
}

func (s *preEncodeWriteJSONSuite) Test_ConditionalGET_IfNoneMatch_Match_Returns304_AndSetsETag() {
	// Prepare payload and compute expected ETag (strong)
	payload := sample{A: 2, B: "etag"}
	body := s.mustMarshal(payload)
	sum := sha256.Sum256(body)
	expectedETag := `"sha256:` + hex.EncodeToString(sum[:]) + `"`

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/x", nil)
	req.Header.Set("If-None-Match", expectedETag)

	apix.PreEncodeAndWriteJSON(rec, req, http.StatusOK, "application/json", payload, false)

	s.Equal(http.StatusNotModified, rec.Code)
	s.Equal(expectedETag, rec.Header().Get("ETag"))
	s.Empty(rec.Header().Get("Content-Type"), "no content-type on 304 bodyless response")
	s.Empty(rec.Body.Bytes())
}

func (s *preEncodeWriteJSONSuite) Test_ConditionalGET_IfNoneMatch_MultipleValues_SomeMatch_Returns304() {
	payload := sample{A: 3, B: "multi"}
	body := s.mustMarshal(payload)
	sum := sha256.Sum256(body)
	matching := `"sha256:` + hex.EncodeToString(sum[:]) + `"`

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/x", nil)
	// Include spaces and a non-matching token prior to the matching one
	req.Header.Set("If-None-Match", `"sha256:deadbeef", `+matching+`, "sha256:cafe"`)
	apix.PreEncodeAndWriteJSON(rec, req, http.StatusOK, "application/json", payload, false)

	s.Equal(http.StatusNotModified, rec.Code)
	s.Equal(matching, rec.Header().Get("ETag"))
	s.Empty(rec.Body.Bytes())
}

func (s *preEncodeWriteJSONSuite) Test_ConditionalGET_IfNoneMatch_Wildcard_Star_Returns304() {
	payload := sample{A: 9, B: "star"}
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/x", nil)
	req.Header.Set("If-None-Match", `*`)

	apix.PreEncodeAndWriteJSON(rec, req, http.StatusOK, "application/json", payload, false)

	s.Equal(http.StatusNotModified, rec.Code)
	s.NotEmpty(rec.Header().Get("ETag"))
	s.Empty(rec.Body.Bytes())
}

func (s *preEncodeWriteJSONSuite) Test_ConditionalGET_IgnoresWhenStatusNot200() {
	// If status != 200 OK, even a matching If-None-Match should NOT 304.
	payload := sample{A: 5, B: "no304"}
	body := s.mustMarshal(payload)
	sum := sha256.Sum256(body)
	matching := `"sha256:` + hex.EncodeToString(sum[:]) + `"`

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/x", nil)
	req.Header.Set("If-None-Match", matching)

	apix.PreEncodeAndWriteJSON(rec, req, http.StatusCreated, "application/json", payload, false)

	s.Equal(http.StatusCreated, rec.Code)
	s.Equal(matching, rec.Header().Get("ETag")) // still sets ETag
	s.NotEmpty(rec.Body.Bytes())                // body is written
}

func (s *preEncodeWriteJSONSuite) Test_EncodingError_WritesProblem500() {
	// json.Marshal on a channel returns an error -> triggers problem writer
	ch := make(chan int)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/x", nil)

	apix.PreEncodeAndWriteJSON(rec, req, http.StatusOK, "application/json", ch, false)

	s.Equal(http.StatusInternalServerError, rec.Code)
	s.Equal(apix.ContentTypeProblemJSON.String(), rec.Header().Get("Content-Type"))
	s.NotEmpty(rec.Body.Bytes())

	// Basic shape of RFC 7807 problem
	var body map[string]any
	s.Require().NoError(json.Unmarshal(rec.Body.Bytes(), &body))
	s.Equal("Encoding error", body["title"])
	s.InDelta(float64(http.StatusInternalServerError), body["status"], 0)
}

func TestPreEncodeAndWriteJSONSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(preEncodeWriteJSONSuite))
}
