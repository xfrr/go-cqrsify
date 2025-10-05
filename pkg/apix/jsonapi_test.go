package apix_test

import (
	"bytes"
	"net/http"
	"testing"

	apix "github.com/xfrr/go-cqrsify/pkg/apix"

	"github.com/stretchr/testify/suite"
)

type userAttrs struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

// ContentType & Stringer
type ContentTypeSuite struct {
	suite.Suite
}

func (s *ContentTypeSuite) Test_String_ReturnsRawValue() {
	ct := apix.ContentTypeJSONAPI
	s.Equal("application/vnd.api+json", ct.String())
}

func TestContentTypeSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(ContentTypeSuite))
}

// SingleDocument & ManyDocument
type DocumentSuite struct {
	suite.Suite
	attrsA userAttrs
	attrsB userAttrs
}

func (s *DocumentSuite) SetupTest() {
	s.attrsA = userAttrs{Name: "Alice", Age: 30}
	s.attrsB = userAttrs{Name: "Bob", Age: 25}
}

func (s *DocumentSuite) Test_NewSingle_BuildsMinimalDocument() {
	doc := apix.NewSingle("users", "u1", s.attrsA)

	s.Equal("users", doc.Data.Type)
	s.Equal("u1", doc.Data.ID)
	s.Equal(s.attrsA, doc.Data.Attributes)
	s.Nil(doc.Links)
	s.Nil(doc.Meta)
	s.Nil(doc.Included)
}

func (s *DocumentSuite) Test_NewMany_BuildsMinimalDocumentAndPreservesOrder() {
	items := []struct {
		ID    string
		Attrs userAttrs
	}{
		{ID: "u1", Attrs: s.attrsA},
		{ID: "u2", Attrs: s.attrsB},
	}
	doc := apix.NewMany("users", items)

	s.Require().Len(doc.Data, 2)
	s.Equal("users", doc.Data[0].Type)
	s.Equal("u1", doc.Data[0].ID)
	s.Equal(s.attrsA, doc.Data[0].Attributes)
	s.Equal("users", doc.Data[1].Type)
	s.Equal("u2", doc.Data[1].ID)
	s.Equal(s.attrsB, doc.Data[1].Attributes)
}

func (s *DocumentSuite) Test_Single_MarshalThenUnmarshal_RoundTrip() {
	in := apix.NewSingle("users", "u1", s.attrsA)

	data, err := in.MarshalJSON()
	s.Require().NoError(err)

	out, err := apix.UnmarshalSingleDocument[userAttrs](data)
	s.Require().NoError(err)
	s.Equal(in, out)
}

func (s *DocumentSuite) Test_Many_MarshalThenUnmarshal_RoundTrip() {
	in := apix.NewMany("users", []struct {
		ID    string
		Attrs userAttrs
	}{
		{"u1", s.attrsA},
		{"u2", s.attrsB},
	})

	data, err := in.MarshalJSON()
	s.Require().NoError(err)

	out, err := apix.UnmarshalManyDocument[userAttrs](data)
	s.Require().NoError(err)
	s.Equal(in, out)
}

func (s *DocumentSuite) Test_UnmarshalSingleDocument_InvalidJSON_ReturnsError() {
	_, err := apix.UnmarshalSingleDocument[userAttrs]([]byte(`{`))
	s.Error(err)
}

func (s *DocumentSuite) Test_UnmarshalManyDocument_InvalidJSON_ReturnsError() {
	_, err := apix.UnmarshalManyDocument[userAttrs]([]byte(`{`))
	s.Error(err)
}

func TestDocumentSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(DocumentSuite))
}

// Content Negotiation helpers
type NegotiationSuite struct {
	suite.Suite
}

func (s *NegotiationSuite) newReq() *http.Request {
	req, _ := http.NewRequest(http.MethodGet, "http://example.com", nil)
	return req
}

func (s *NegotiationSuite) Test_IsJSONAPIContentNegotiable_Table() {
	tests := []struct {
		name    string
		accept  string
		allowed bool
	}{
		{"EmptyAccept_Allows", "", true},
		{"Wildcard_Allows", "*/*", true},
		{"ExactJSONAPI_Allows", "application/vnd.api+json", true},
		{"WeightedJSONAPI_Allows", "application/json, application/vnd.api+json; q=0.9", true},
		{"ProblemJSON_Rejects", "application/problem+json", false},
		{"PlainJSON_Rejects", "application/json", false},
		{"ListWithoutToken_Rejects", "text/html,application/xml;q=0.8", false},
		{"SupersetTokenSubstring_AllowsBecauseSimpleContains", "application/vnd.api+json; charset=utf-8", true},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			req := s.newReq()
			req.Header.Set("Accept", tt.accept)
			got := apix.IsJSONAPIContentNegotiable(req)
			s.Equal(tt.allowed, got)
		})
	}
}

func (s *NegotiationSuite) Test_IsJSONAPIContentType_Table() {
	tests := []struct {
		name        string
		contentType string
		want        bool
	}{
		{"Empty_False", "", false},
		{"ExactJSONAPI_True", "application/vnd.api+json", true},
		{"WithParams_True", "application/vnd.api+json; charset=utf-8", true},
		{"Different_False", "application/json", false},
		{"InvalidValue_False", ";;", false},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			req := s.newReq()
			if tt.contentType != "" {
				req.Header.Set("Content-Type", tt.contentType)
			}
			got := apix.IsJSONAPIContentType(req)
			s.Equal(tt.want, got)
		})
	}
}

func TestNegotiationSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(NegotiationSuite))
}

// Internal containsToken / contains helpers (indirect behavioral test)
type ContainsBehaviorSuite struct {
	suite.Suite
}

func (s *ContainsBehaviorSuite) Test_SubstringMatchingBehavior() {
	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Accept", "text/plain, application/vnd.api+json;version=1; q=0.7")
	s.True(apix.IsJSONAPIContentNegotiable(req))

	req2, _ := http.NewRequest(http.MethodGet, "/", nil)
	req2.Header.Set("Accept", "text/plain,application/vnd.api+jsonx")
	s.True(apix.IsJSONAPIContentNegotiable(req2))
}

func (s *ContainsBehaviorSuite) Test_LongHeaderWithTokenInside() {
	var buf bytes.Buffer
	buf.WriteString("text/html; q=0.5, ")
	for range 10 {
		buf.WriteString("application/xml; q=0.1, ")
	}
	buf.WriteString("application/vnd.api+json; q=1.0")

	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Accept", buf.String())

	s.True(apix.IsJSONAPIContentNegotiable(req))
}

func TestContainsBehaviorSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(ContainsBehaviorSuite))
}
