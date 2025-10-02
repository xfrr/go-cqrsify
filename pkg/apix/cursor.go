package apix

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/url"
	"strings"
)

// CursorSigner creates/verifies opaque signed cursors.
type CursorSigner struct {
	secret []byte
}

func NewCursorSigner(secret []byte) CursorSigner {
	return CursorSigner{secret: append([]byte(nil), secret...)}
}

// Sign encodes payload as base64url(JSON || "." || hex(HMAC256))
func (s CursorSigner) Sign(payload any) (string, error) {
	b, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}
	mac := hmac.New(sha256.New, s.secret)
	mac.Write(b)
	sig := mac.Sum(nil)

	encPayload := base64.RawURLEncoding.EncodeToString(b)
	encSig := base64.RawURLEncoding.EncodeToString(sig)
	return encPayload + "." + encSig, nil
}

// Verify decodes and validates signature; returns raw JSON bytes.
func (s CursorSigner) Verify(token string) ([]byte, error) {
	dot := strings.LastIndexByte(token, '.')
	if dot <= 0 {
		return nil, errors.New("invalid cursor")
	}
	pEnc := token[:dot]
	sEnc := token[dot+1:]

	p, err := base64.RawURLEncoding.DecodeString(pEnc)
	if err != nil {
		return nil, errors.New("invalid cursor payload")
	}
	sig, err := base64.RawURLEncoding.DecodeString(sEnc)
	if err != nil {
		return nil, errors.New("invalid cursor signature")
	}
	mac := hmac.New(sha256.New, s.secret)
	mac.Write(p)
	if !hmac.Equal(sig, mac.Sum(nil)) {
		return nil, errors.New("cursor signature mismatch")
	}
	return p, nil
}

// BuildCursorLinks builds JSON:API pagination links with opaque signed cursors.
// basePath should be a clean path like "/users".
func BuildCursorLinks(basePath string, selfQuery url.Values, nextCursor, prevCursor string) *Links {
	l := &Links{
		Self: basePath,
	}
	if len(selfQuery) > 0 {
		l.Self = basePath + "?" + selfQuery.Encode()
	}
	// Next
	if nextCursor != "" {
		q := cloneValues(selfQuery)
		q.Set("page[cursor]", nextCursor)
		l.Next = basePath + "?" + q.Encode()
	}
	// Prev
	if prevCursor != "" {
		q := cloneValues(selfQuery)
		q.Set("page[cursor]", prevCursor)
		l.Prev = basePath + "?" + q.Encode()
	}
	return l
}

func cloneValues(v url.Values) url.Values {
	out := make(url.Values, len(v))
	for k, vals := range v {
		cp := make([]string, len(vals))
		copy(cp, vals)
		out[k] = cp
	}
	return out
}
