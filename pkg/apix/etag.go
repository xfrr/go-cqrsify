package apix

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"strings"
)

// StrongETagFromBytes returns a strong ETag: "sha256:<hex>"
func StrongETagFromBytes(b []byte) string {
	h := sha256.Sum256(b)
	return `"sha256:` + hex.EncodeToString(h[:]) + `"`
}

// WeakETagFromBytes returns a weak ETag: W/"sha256:<hex>"
func WeakETagFromBytes(b []byte) string {
	h := sha256.Sum256(b)
	return `W/"sha256:` + hex.EncodeToString(h[:]) + `"`
}

// PreEncodeAndWriteJSON encodes v once (for ETag) and writes it with conditional 304 support.
// If If-None-Match matches the computed ETag, this returns 304 without a body.
func PreEncodeAndWriteJSON(w http.ResponseWriter, r *http.Request, status int, contentType string, v any, weak bool) {
	// Pre-encode once for both ETag and body
	b, err := json.Marshal(v)
	if err != nil {
		// Fallback to problem (rare)
		WriteProblem(w, NewProblem(http.StatusInternalServerError, "Encoding error", err.Error()))
		return
	}
	var etag string
	if weak {
		etag = WeakETagFromBytes(b)
	} else {
		etag = StrongETagFromBytes(b)
	}

	// Conditional GET (If-None-Match)
	if inm := r.Header.Get("If-None-Match"); inm != "" && etagMatches(inm, etag) && status == http.StatusOK {
		w.Header().Set("ETag", etag)
		w.WriteHeader(http.StatusNotModified)
		return
	}

	w.Header().Set("ETag", etag)
	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(status)
	_, _ = w.Write(b)
}

// etagMatches performs a simple token compare; supports multiple values.
func etagMatches(header, etag string) bool {
	parts := strings.SplitSeq(header, ",")
	for p := range parts {
		if strings.TrimSpace(p) == etag || strings.TrimSpace(p) == "*" {
			return true
		}
	}
	return false
}
