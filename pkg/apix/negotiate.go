package apix

import (
	"net/http"
	"strconv"
	"strings"
)

// mediaRange represents a parsed Accept entry with optional q weight.
type mediaRange struct {
	Type string
	Sub  string
	Q    float64
	// we keep it minimal; parameters beyond q= are ignored for performance
}

// Matches "type/sub" against Accept ranges (with wildcards), honoring client order and q.
func Accepts(r *http.Request, mediaType string) bool {
	m := parseAccept(r.Header.Get("Accept"))
	if len(m) == 0 {
		return true
	}
	wantType, wantSub := splitType(mediaType)
	for _, mr := range m {
		if mr.Q <= 0 {
			continue
		}
		if (mr.Type == "*" || mr.Type == wantType) && (mr.Sub == "*" || mr.Sub == wantSub) {
			return true
		}
	}
	return false
}

func splitType(mt string) (string, string) {
	for i := range mt {
		if mt[i] == '/' {
			return strings.ToLower(mt[:i]), strings.ToLower(mt[i+1:])
		}
	}
	return strings.ToLower(mt), "*"
}

// RequireAccept writes 406 Problem if Accept doesn't allow the mediaType.
func RequireAccept(w http.ResponseWriter, r *http.Request, mediaType string) bool {
	if Accepts(r, mediaType) {
		return true
	}
	WriteProblem(w, NewProblem(http.StatusNotAcceptable, "Not acceptable",
		"Accept must allow "+mediaType,
		WithType("https://example.com/problems/not-acceptable"),
		WithInstance(r.URL.Path)))
	return false
}

// RequireContentType enforces an exact Content-Type match (no parameters) and writes 415 if not.
func RequireContentType(w http.ResponseWriter, r *http.Request, contentType string) bool {
	ct := r.Header.Get("Content-Type")
	if ct == contentType {
		return true
	}
	WriteProblem(w, NewProblem(http.StatusUnsupportedMediaType, "Unsupported media type",
		"Content-Type must be "+contentType,
		WithType("https://example.com/problems/unsupported-media-type"),
		WithInstance(r.URL.Path)))
	return false
}

// RequireJSONAPIContentType is a shortcut for RequireContentType with JSON:API.
func RequireJSONAPIContentType(w http.ResponseWriter, r *http.Request) bool {
	return RequireContentType(w, r, ContentTypeJSONAPI.String())
}

// RequireMergePatchContentType is a shortcut for RequireContentType with Merge Patch.
func RequireMergePatchContentType(w http.ResponseWriter, r *http.Request) bool {
	return RequireContentType(w, r, ContentTypeMergePatch.String())
}

// RequireProblemContentType is a shortcut for RequireContentType with application/problem+json.
func RequireProblemContentType(w http.ResponseWriter, r *http.Request) bool {
	return RequireContentType(w, r, ContentTypeProblemJSON.String())
}

// RequireJSONContentType is a shortcut for RequireContentType with application/json.
func RequireJSONContentType(w http.ResponseWriter, r *http.Request) bool {
	return RequireContentType(w, r, ContentTypeJSON.String())
}

// parseAccept implements a tiny (but strict-enough) Accept parser with q-values.
// Examples:
//
//	Accept: application/vnd.api+json; q=1.0, application/json; q=0.8, */*;q=0.1
func parseAccept(s string) []mediaRange {
	if s == "" {
		return []mediaRange{{Type: "*", Sub: "*", Q: 1.0}}
	}
	parts := strings.Split(s, ",")
	out := make([]mediaRange, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		if mr, ok := parseMediaRange(p); ok {
			out = append(out, *mr)
		}
	}
	// We keep order (sender preference). Server will pick first matching we support.
	return out
}

func parseMediaRange(p string) (*mediaRange, bool) {
	typ := p
	q := 1.0
	if semi := strings.IndexByte(p, ';'); semi >= 0 {
		typ = strings.TrimSpace(p[:semi])
		paramStr := p[semi+1:]
		ps := strings.SplitSeq(paramStr, ";")
		for kv := range ps {
			kv = strings.TrimSpace(kv)
			if kv == "" {
				continue
			}
			if strings.HasPrefix(strings.ToLower(kv), "q=") {
				val := strings.TrimSpace(kv[2:])
				// very small fast-path parser; defaults on error
				if v, err := parseQ(val); err == nil {
					q = v
				}
			}
		}
	}
	if slash := strings.IndexByte(typ, '/'); slash >= 0 {
		return &mediaRange{
			Type: strings.ToLower(strings.TrimSpace(typ[:slash])),
			Sub:  strings.ToLower(strings.TrimSpace(typ[slash+1:])),
			Q:    q,
		}, true
	}
	return nil, false
}

func parseQ(s string) (float64, error) {
	// q-values are 0..1 with up to 3 decimals.
	return strconv.ParseFloat(s, 64)
}
