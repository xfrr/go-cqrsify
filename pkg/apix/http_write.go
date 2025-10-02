package apix

import (
	"encoding/json"
	"net/http"
)

const (
	// CacheControlHeaderKey is the standard HTTP Cache-Control header name.
	CacheControlHeaderKey = "Cache-Control"
	// LastModifiedHeaderKey is the standard HTTP Last-Modified header name.
	LastModifiedHeaderKey = "Last-Modified"
	// ETagHeaderKey is the standard HTTP ETag header name.
	ETagHeaderKey = "ETag"
	// VaryHeaderKey is the standard HTTP Vary header name.
	VaryHeaderKey = "Vary"
	// IfModifiedSinceHeaderKey is the standard HTTP If-Modified-Since header name.
	IfModifiedSinceHeaderKey = "If-Modified-Since"
	// ContentTypeHeaderKey is the standard HTTP Content-Type header name.
	ContentTypeHeaderKey = "Content-Type"
)

// WriteJSON writes a JSON response with the given status and content type.
// Additional options can be provided via WriteOption.
//
// If no ContentType is set, defaults to application/vnd.api+json.
// If no Status is set, defaults to 200 OK.
func WriteJSON(w http.ResponseWriter, v any, opts ...WriteOption) {
	options := new(WriteOptions)
	options.apply(opts...)

	b, err := json.Marshal(v)
	if err != nil {
		WriteProblem(w, NewProblem(http.StatusInternalServerError, "Encoding error", err.Error()))
		return
	}

	switch {
	case options.UseWeakETag:
		w.Header().Set(ETagHeaderKey, WeakETagFromBytes(b))
	case options.ETag != "":
		w.Header().Set(ETagHeaderKey, options.ETag)
	default:
		w.Header().Set(ETagHeaderKey, StrongETagFromBytes(b))
	}

	w.Header().Set(ContentTypeHeaderKey, options.ContentType.String())

	if options.CacheControl != "" {
		w.Header().Set(CacheControlHeaderKey, options.CacheControl)
	}
	if !options.LastModified.IsZero() {
		w.Header().Set(LastModifiedHeaderKey, options.LastModified.UTC().Format(http.TimeFormat))
	}
	if options.Vary != "" {
		w.Header().Set(VaryHeaderKey, options.Vary)
	}
	if !options.IfModifiedSince.IsZero() {
		w.Header().Set(IfModifiedSinceHeaderKey, options.IfModifiedSince.UTC().Format(http.TimeFormat))
	}
	if options.Headers != nil {
		for k, v := range options.Headers {
			if k != "" && v != "" {
				w.Header().Set(k, v)
			}
		}
	}

	w.WriteHeader(options.Status)
	_, _ = w.Write(b)
}

// IsJSONAPIContentNegotiable checks if the request Accept header is compatible.
// You can use it to return 406 Not Acceptable if needed.
func IsJSONAPIContentNegotiable(r *http.Request) bool {
	accept := r.Header.Get("Accept")
	return accept == "" ||
		accept == "*/*" ||
		containsToken(accept, ContentTypeJSONAPI.String())
}

// containsToken is a simple substring matcher for media types in Accept.
// TODO: implement a weighted media type parser.
func containsToken(s, token string) bool {
	return len(s) >= len(token) && (s == token || (len(s) > len(token) && (contains(s, token))))
}

func contains(s, sub string) bool {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
