package apix_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/suite"
	apix "github.com/xfrr/go-cqrsify/pkg/apix"
)

type ProblemJSONSuite struct {
	suite.Suite
}

func (s *ProblemJSONSuite) decode(b []byte) map[string]any {
	var m map[string]any
	s.Require().NoError(json.Unmarshal(b, &m))
	return m
}

func (s *ProblemJSONSuite) Test_MarshalJSON_MergesExtensions_WithoutClobber() {
	p := apix.NewProblem(
		http.StatusBadRequest,
		"Bad Request",
		"explanation",
		apix.WithType("https://e/x"),
		apix.WithInstance("/req/123"),
		apix.WithExtensions(map[string]any{
			"extra":  "v",
			"status": 999,             // should NOT override base field
			"title":  "should_ignore", // should NOT override base field
		}),
	)

	raw, err := p.MarshalJSON()
	s.Require().NoError(err)

	got := s.decode(raw)

	// base fields present
	s.Equal("https://e/x", got["type"])
	s.Equal("Bad Request", got["title"])
	s.Equal(float64(http.StatusBadRequest), got["status"])
	s.Equal("explanation", got["detail"])
	s.Equal("/req/123", got["instance"])

	// extension merged
	s.Equal("v", got["extra"])

	// clobbering prevented
	s.NotEqual(float64(999), got["status"])
	s.NotEqual("should_ignore", got["title"])

	// no "Extensions" key should be serialized
	_, has := got["Extensions"]
	s.False(has)
}

func (s *ProblemJSONSuite) Test_MarshalJSON_NoExtensions() {
	p := apix.NewProblem(http.StatusNotFound, "Not Found", "")
	raw, err := p.MarshalJSON()
	s.Require().NoError(err)
	got := s.decode(raw)

	s.Equal("about:blank", got["type"])
	s.Equal("Not Found", got["title"])
	s.Equal(float64(http.StatusNotFound), got["status"])
	_, hasExt := got["Extensions"]
	s.False(hasExt)
}

func TestProblemJSONSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(ProblemJSONSuite))
}

type ProblemBuildersSuite struct {
	suite.Suite
}

func (s *ProblemBuildersSuite) Test_NewProblem_WithOptions() {
	p := apix.NewProblem(
		http.StatusForbidden,
		"Forbidden",
		"nope",
		apix.WithType("https://e/forbidden"),
		apix.WithInstance("/ops/1"),
		apix.WithExtensions(map[string]any{"code": "E123"}),
	)

	s.Equal("https://e/forbidden", p.Type)
	s.Equal("Forbidden", p.Title)
	s.Equal(http.StatusForbidden, p.Status)
	s.Equal("nope", p.Detail)
	s.Equal("/ops/1", p.Instance)
	s.Equal("E123", p.Extensions["code"])
}

func (s *ProblemBuildersSuite) Test_Shorthand_Constructors() {
	s.Run("Conflict", func() {
		p := apix.NewConflictProblem("x")
		s.Equal(http.StatusConflict, p.Status)
		s.Equal("Conflict", p.Title)
		s.Equal("x", p.Detail)
	})
	s.Run("BadRequest", func() {
		p := apix.NewBadRequestProblem("y")
		s.Equal(http.StatusBadRequest, p.Status)
		s.Equal("Bad Request", p.Title)
		s.Equal("y", p.Detail)
	})
	s.Run("NotFound", func() {
		p := apix.NewNotFoundProblem("z")
		s.Equal(http.StatusNotFound, p.Status)
		s.Equal("Not Found", p.Title)
		s.Equal("z", p.Detail)
	})
	s.Run("Unauthorized", func() {
		p := apix.NewUnauthorizedProblem("u")
		s.Equal(http.StatusUnauthorized, p.Status)
		s.Equal("Unauthorized", p.Title)
		s.Equal("u", p.Detail)
	})
	s.Run("Forbidden", func() {
		p := apix.NewForbiddenProblem("f")
		s.Equal(http.StatusForbidden, p.Status)
		s.Equal("Forbidden", p.Title)
		s.Equal("f", p.Detail)
	})
	s.Run("InternalServerError", func() {
		p := apix.NewInternalServerErrorProblem("e")
		s.Equal(http.StatusInternalServerError, p.Status)
		s.Equal("Internal Server Error", p.Title)
		s.Equal("e", p.Detail)
	})
	s.Run("UnsupportedMediaType", func() {
		p := apix.NewUnsupportedMediaTypeProblem("m")
		s.Equal(http.StatusUnsupportedMediaType, p.Status)
		s.Equal("Unsupported Media Type", p.Title)
		s.Equal("m", p.Detail)
	})
}

func TestProblemBuildersSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(ProblemBuildersSuite))
}

type WriteProblemSuite struct {
	suite.Suite
}

func (s *WriteProblemSuite) decodeRecorder(rec *httptest.ResponseRecorder) map[string]any {
	return s.decode(rec.Body.Bytes())
}

func (s *WriteProblemSuite) decode(b []byte) map[string]any {
	var m map[string]any
	s.Require().NoError(json.Unmarshal(b, &m))
	return m
}

func (s *WriteProblemSuite) Test_WriteProblem_SetsHeaders_Status_AndBody() {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/things/42", nil)

	prob := apix.NewBadRequestProblem(
		"bad input",
		apix.WithInstance(req.URL.Path),
		apix.WithExtensions(map[string]any{"field": "name"}),
	)

	apix.WriteProblem(rec, prob)

	// status (WriteJSON should set it based on WithStatus)
	s.Equal(http.StatusBadRequest, rec.Code)

	// headers
	s.Equal(apix.ContentTypeProblemJSON.String(), rec.Header().Get("Content-Type"))
	s.Equal("no-store", rec.Header().Get("Cache-Control"))

	// body
	body := s.decodeRecorder(rec)
	s.Equal("about:blank", body["type"])
	s.Equal("Bad Request", body["title"])
	s.Equal(float64(http.StatusBadRequest), body["status"])
	s.Equal("bad input", body["detail"])
	s.Equal("/things/42", body["instance"])
	s.Equal("name", body["field"])
}

func TestWriteProblemSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(WriteProblemSuite))
}
