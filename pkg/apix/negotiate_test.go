package apix_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/suite"
	apix "github.com/xfrr/go-cqrsify/pkg/apix"
)

type AcceptsSuite struct {
	suite.Suite
	req *http.Request
}

func (s *AcceptsSuite) SetupTest() {
	r, _ := http.NewRequest(http.MethodGet, "http://example.test/resource", nil)
	s.req = r
}

func (s *AcceptsSuite) Test_EmptyAccept_AllowsAny() {
	s.req.Header.Del("Accept")
	s.True(apix.Accepts(s.req, "application/json"))
	s.True(apix.Accepts(s.req, "text/html"))
}

func (s *AcceptsSuite) Test_Wildcard_Allows() {
	s.req.Header.Set("Accept", "*/*")
	s.True(apix.Accepts(s.req, "application/json"))
}

func (s *AcceptsSuite) Test_WildcardSubtype_Allows() {
	s.req.Header.Set("Accept", "application/*")
	s.True(apix.Accepts(s.req, "application/json"))
	s.True(apix.Accepts(s.req, "application/vnd.api+json"))
	s.False(apix.Accepts(s.req, "text/html"))
}

func (s *AcceptsSuite) Test_WildcardType_Allows() {
	s.req.Header.Set("Accept", "*/json")
	s.True(apix.Accepts(s.req, "application/json"))
	s.True(apix.Accepts(s.req, "text/json"))
	s.False(apix.Accepts(s.req, "text/html"))
}

func (s *AcceptsSuite) Test_ExactMatch_Allows() {
	s.req.Header.Set("Accept", "application/vnd.api+json")
	s.True(apix.Accepts(s.req, apix.ContentTypeJSONAPI.String()))
	s.False(apix.Accepts(s.req, "application/json"))
}

func (s *AcceptsSuite) Test_OrderAndQValues_Respected() {
	// q=0 entries must be ignored; earlier acceptable entries should allow.
	s.req.Header.Set("Accept", "application/xml;q=0, application/json;q=1.0")
	s.True(apix.Accepts(s.req, "application/json"))
	s.False(apix.Accepts(s.req, "application/xml"))
}

func (s *AcceptsSuite) Test_MultipleEntries_FirstUsableWins() {
	s.req.Header.Set("Accept", "text/html, application/json")
	s.True(apix.Accepts(s.req, "application/json"))
	s.False(apix.Accepts(s.req, "text/plain"))
}

func (s *AcceptsSuite) Test_SplitType_WithoutSlash_TreatedAsTypeOnly() {
	// When server asks for "application" (type only), splitType returns ("application","*")
	// so it should match application/* ranges as acceptable.
	s.req.Header.Set("Accept", "application/*")
	s.True(apix.Accepts(s.req, "application"))
}

func TestAcceptsSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(AcceptsSuite))
}

type RequireAcceptSuite struct {
	suite.Suite
}

func (s *RequireAcceptSuite) Test_ReturnsTrue_WhenAccepts() {
	req, _ := http.NewRequest(http.MethodGet, "/x", nil)
	req.Header.Set("Accept", "application/json")
	rec := httptest.NewRecorder()

	ok := apix.RequireAccept(rec, req, "application/json")
	s.True(ok)
}

func (s *RequireAcceptSuite) Test_ReturnsFalse_WritesProblem406_WhenNotAcceptable() {
	req, _ := http.NewRequest(http.MethodGet, "/x", nil)
	req.Header.Set("Accept", "text/html")
	rec := httptest.NewRecorder()

	ok := apix.RequireAccept(rec, req, "application/json")

	s.False(ok)
	// WriteProblem is expected to set 406 and write a body.
	s.Equal(http.StatusNotAcceptable, rec.Code)
	s.NotEmpty(rec.Body.String())
}

func TestRequireAcceptSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(RequireAcceptSuite))
}

type RequireContentTypeSuite struct {
	suite.Suite
}

func (s *RequireContentTypeSuite) Test_ExactMatch_ReturnsTrue() {
	req, _ := http.NewRequest(http.MethodPost, "/y", nil)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	ok := apix.RequireContentType(rec, req, "application/json")
	s.True(ok)
}

func (s *RequireContentTypeSuite) Test_WithParameters_Fails_MustBeExact() {
	req, _ := http.NewRequest(http.MethodPost, "/y", nil)
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	rec := httptest.NewRecorder()

	ok := apix.RequireContentType(rec, req, "application/json")
	s.False(ok)
	s.Equal(http.StatusUnsupportedMediaType, rec.Code)
	s.NotEmpty(rec.Body.String())
}

func (s *RequireContentTypeSuite) Test_DifferentType_Fails_415() {
	req, _ := http.NewRequest(http.MethodPost, "/y", nil)
	req.Header.Set("Content-Type", "text/plain")
	rec := httptest.NewRecorder()

	ok := apix.RequireContentType(rec, req, "application/json")
	s.False(ok)
	s.Equal(http.StatusUnsupportedMediaType, rec.Code)
}

func (s *RequireContentTypeSuite) Test_Shortcuts_JSONAPI() {
	req, _ := http.NewRequest(http.MethodPost, "/y", nil)
	req.Header.Set("Content-Type", apix.ContentTypeJSONAPI.String())
	rec := httptest.NewRecorder()

	s.True(apix.RequireJSONAPIContentType(rec, req))
}

func (s *RequireContentTypeSuite) Test_Shortcuts_MergePatch() {
	req, _ := http.NewRequest(http.MethodPost, "/y", nil)
	req.Header.Set("Content-Type", apix.ContentTypeMergePatch.String())
	rec := httptest.NewRecorder()

	s.True(apix.RequireMergePatchContentType(rec, req))
}

func (s *RequireContentTypeSuite) Test_Shortcuts_ProblemJSON() {
	req, _ := http.NewRequest(http.MethodPost, "/y", nil)
	req.Header.Set("Content-Type", apix.ContentTypeProblemJSON.String())
	rec := httptest.NewRecorder()

	s.True(apix.RequireProblemContentType(rec, req))
}

func (s *RequireContentTypeSuite) Test_Shortcuts_JSON() {
	req, _ := http.NewRequest(http.MethodPost, "/y", nil)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	s.True(apix.RequireJSONContentType(rec, req))
}

func TestRequireContentTypeSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(RequireContentTypeSuite))
}
