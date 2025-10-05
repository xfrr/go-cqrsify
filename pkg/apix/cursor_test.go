package apix_test

import (
	"encoding/base64"
	"encoding/json"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"
	apix "github.com/xfrr/go-cqrsify/pkg/apix"
)

// ------------------------------------------------------------
// CursorSigner tests
// ------------------------------------------------------------

type CursorSignerSuite struct {
	suite.Suite
	secret []byte
	signer apix.CursorSigner
}

func (s *CursorSignerSuite) SetupTest() {
	s.secret = []byte("super-secret-key")
	s.signer = apix.NewCursorSigner(s.secret)
}

func (s *CursorSignerSuite) Test_NewCursorSigner_CopiesSecret_NotAliased() {
	payload := map[string]any{"n": 1}

	// Mutate original secret AFTER creating the signer
	s.secret[0] = 'X'

	token, err := s.signer.Sign(payload)
	s.Require().NoError(err)

	// Should still verify with the signer (internal copy intact)
	out, err := s.signer.Verify(token)
	s.Require().NoError(err)
	var got map[string]any
	s.Require().NoError(json.Unmarshal(out, &got))
	s.InDelta(float64(1), got["n"], 0.0000001)

	// A new signer made with the mutated secret should FAIL to verify the old token
	signerWithMutated := apix.NewCursorSigner(s.secret)
	_, err = signerWithMutated.Verify(token)
	s.Error(err)
}

func (s *CursorSignerSuite) Test_SignAndVerify_RoundTrip_MapPayload() {
	payload := map[string]any{"id": "u1", "page": 42}

	token, err := s.signer.Sign(payload)
	s.Require().NoError(err)
	s.NotEmpty(token)
	s.Contains(token, ".", "token must have payload.signature format")

	raw, err := s.signer.Verify(token)
	s.Require().NoError(err)

	var got map[string]any
	s.Require().NoError(json.Unmarshal(raw, &got))
	s.Equal("u1", got["id"])
	s.InDelta(float64(42), got["page"], 0.0000001)
}

func (s *CursorSignerSuite) Test_Sign_Error_NotSerializable() {
	_, err := s.signer.Sign(make(chan int)) // json.Marshal should error
	s.Error(err)
}

func (s *CursorSignerSuite) Test_Verify_TamperedSignature_Fails() {
	payload := map[string]any{"x": "y"}
	token, err := s.signer.Sign(payload)
	s.Require().NoError(err)

	// Break signature by flipping one byte
	dot := strings.LastIndexByte(token, '.')
	s.Require().Positive(dot, "token must have a dot")

	pEnc := token[:dot]
	sEnc := token[dot+1:]

	sig, err := base64.RawURLEncoding.DecodeString(sEnc)
	s.Require().NoError(err)
	sig[0] ^= 0xFF // flip bits
	badSig := base64.RawURLEncoding.EncodeToString(sig)
	badToken := pEnc + "." + badSig

	_, err = s.signer.Verify(badToken)
	s.EqualError(err, "cursor signature mismatch")
}

func (s *CursorSignerSuite) Test_Verify_InvalidFormats() {
	// No dot
	_, err := s.signer.Verify("abc")
	s.Require().Error(err, "invalid cursor")

	// Invalid payload base64
	_, err = s.signer.Verify("%.%")
	s.Require().EqualError(err, "invalid cursor payload")

	// Valid payload / invalid signature base64
	p := base64.RawURLEncoding.EncodeToString([]byte(`{"a":1}`))
	_, err = s.signer.Verify(p + "." + "%")
	s.EqualError(err, "invalid cursor signature")
}

func TestCursorSignerSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(CursorSignerSuite))
}

// ------------------------------------------------------------
// BuildCursorLinks tests
// ------------------------------------------------------------

type BuildCursorLinksSuite struct {
	suite.Suite
}

func (s *BuildCursorLinksSuite) parseURL(u string) *url.URL {
	parsed, err := url.Parse(u)
	s.Require().NoError(err)
	return parsed
}

func (s *BuildCursorLinksSuite) Test_Self_NoQuery() {
	links := apix.BuildCursorLinks("/users", url.Values{}, "", "")
	s.Equal("/users", links.Self)
	s.Empty(links.Next)
	s.Empty(links.Prev)
}

func (s *BuildCursorLinksSuite) Test_Self_WithQuery_AndNextPrevCursors() {
	q := url.Values{}
	q.Set("filter[active]", "true")
	q.Set("sort", "-created_at")

	links := apix.BuildCursorLinks("/users", q, "NEXT123", "PREV456")

	// Self contains original query
	self := s.parseURL(links.Self)
	s.Equal("/users", self.Path)
	s.Equal("true", self.Query().Get("filter[active]"))
	s.Equal("-created_at", self.Query().Get("sort"))

	// Next has page[cursor]=NEXT123 plus original params
	next := s.parseURL(links.Next)
	s.Equal("/users", next.Path)
	s.Equal("NEXT123", next.Query().Get("page[cursor]"))
	s.Equal("true", next.Query().Get("filter[active]"))
	s.Equal("-created_at", next.Query().Get("sort"))

	// Prev has page[cursor]=PREV456 plus original params
	prev := s.parseURL(links.Prev)
	s.Equal("/users", prev.Path)
	s.Equal("PREV456", prev.Query().Get("page[cursor]"))
	s.Equal("true", prev.Query().Get("filter[active]"))
	s.Equal("-created_at", prev.Query().Get("sort"))
}

func (s *BuildCursorLinksSuite) Test_OnlyNext() {
	q := url.Values{"limit": []string{"10"}}
	links := apix.BuildCursorLinks("/items", q, "NEXT", "")

	s.NotEmpty(links.Next)
	s.Empty(links.Prev)

	u := s.parseURL(links.Next)
	s.Equal("/items", u.Path)
	s.Equal("10", u.Query().Get("limit"))
	s.Equal("NEXT", u.Query().Get("page[cursor]"))
}

func (s *BuildCursorLinksSuite) Test_OnlyPrev() {
	links := apix.BuildCursorLinks("/items", url.Values{"q": []string{"abc"}}, "", "P")
	s.Empty(links.Next)
	s.NotEmpty(links.Prev)

	u := s.parseURL(links.Prev)
	s.Equal("/items", u.Path)
	s.Equal("abc", u.Query().Get("q"))
	s.Equal("P", u.Query().Get("page[cursor]"))
}

func (s *BuildCursorLinksSuite) Test_QueryValuesCloned_NotMutatedByLaterChanges() {
	q := url.Values{}
	q.Set("a", "1")

	links := apix.BuildCursorLinks("/things", q, "N", "P")

	// Mutate the original after building
	q.Set("a", "2")

	// Links were built from a clone; they must not reflect later changes
	next := s.parseURL(links.Next)
	prev := s.parseURL(links.Prev)
	s.Equal("1", next.Query().Get("a"))
	s.Equal("1", prev.Query().Get("a"))
}

func (s *BuildCursorLinksSuite) Test_EncodesSpecialCharacters() {
	q := url.Values{}
	q.Set("name", "A B") // space should encode to '+'
	links := apix.BuildCursorLinks("/people", q, "CUR", "PRV")

	self := links.Self
	s.Contains(self, "name=A+B")

	s.Contains(links.Next, "page%5Bcursor%5D=CUR") // "page[cursor]" encoded
	s.Contains(links.Prev, "page%5Bcursor%5D=PRV")
}

func TestBuildCursorLinksSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(BuildCursorLinksSuite))
}
