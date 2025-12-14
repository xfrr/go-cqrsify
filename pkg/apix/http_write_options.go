package apix

import (
	"maps"
	"time"
)

type WriteOptions struct {
	// EscapeHTML tells the JSON encoder to escape HTML characters.
	EscapeHTML bool

	// ContentType to set; defaults to application/vnd.api+json.
	ContentType ContentType

	// Status code to use; defaults to 200 OK.
	Status int

	// provides directives to control caching behavior for both client
	// requests and server responses, primarily to optimize performance,
	// reduce server load, and minimize bandwidth usage.
	CacheControl string

	// LastModified indicates the date and time when the origin server
	// believes a resource was last modified, formatted in GMT.
	LastModified time.Time

	// Vary header informs caches which request headers influence
	// the content of the response, ensuring that the correct cached version
	// is served based on those headers
	Vary string

	// If-Modified-Since makes a request conditional, instructing the server
	// to return the requested resource with a 200 status code only if
	// it has been modified after the date specified in the header
	IfModifiedSince time.Time

	// UseWeakETag indicates if ETag should be weak (W/) or strong.
	UseWeakETag bool

	// ETag value to set.
	ETag string

	// Additional custom headers to set (if any).
	Headers map[string]string
}

func (o *WriteOptions) apply(opts ...WriteOption) {
	for _, opt := range opts {
		opt(o)
	}
	if o.ContentType == "" {
		o.ContentType = ContentTypeJSONAPI
	}
	if o.Status == 0 {
		o.Status = 200
	}
	if o.Headers != nil {
		for k, v := range o.Headers {
			if k != "" && v != "" {
				// Avoid overwriting standard headers
				switch k {
				case ContentTypeHeaderKey,
					CacheControlHeaderKey,
					LastModifiedHeaderKey,
					ETagHeaderKey,
					VaryHeaderKey,
					IfModifiedSinceHeaderKey:
				default:
					// Custom header
					o.Headers[k] = v
				}
			}
		}
	}
}

type WriteOption func(*WriteOptions)

// WithEscapeHTML enables HTML escaping in JSON output.
func WithEscapeHTML(escape bool) WriteOption {
	return func(o *WriteOptions) {
		o.EscapeHTML = escape
	}
}

// WithContentType sets the Content-Type header.
func WithContentType(ct ContentType) WriteOption {
	return func(o *WriteOptions) {
		o.ContentType = ct
	}
}

// WithStatusCode sets the HTTP status code.
func WithStatusCode(status int) WriteOption {
	return func(o *WriteOptions) {
		o.Status = status
	}
}

// WithCacheControl sets the Cache-Control header.
func WithCacheControl(cc string) WriteOption {
	return func(o *WriteOptions) {
		o.CacheControl = cc
	}
}

// WithLastModified sets the Last-Modified header.
func WithLastModified(t time.Time) WriteOption {
	return func(o *WriteOptions) {
		o.LastModified = t
	}
}

// WithVary sets the Vary header.
func WithVary(vary string) WriteOption {
	return func(o *WriteOptions) {
		o.Vary = vary
	}
}

// WithIfModifiedSince sets the If-Modified-Since header.
func WithIfModifiedSince(t time.Time) WriteOption {
	return func(o *WriteOptions) {
		o.IfModifiedSince = t
	}
}

// WithWeakETag sets whether the ETag should be weak (W/).
func WithWeakETag(weak bool) WriteOption {
	return func(o *WriteOptions) {
		o.UseWeakETag = weak
	}
}

// WithETag sets the ETag header value.
func WithETag(etag string) WriteOption {
	return func(o *WriteOptions) {
		o.ETag = etag
	}
}

// WithHeaders adds custom headers to the response.
func WithHeaders(headers map[string]string) WriteOption {
	return func(o *WriteOptions) {
		if o.Headers == nil {
			o.Headers = make(map[string]string)
		}
		maps.Copy(o.Headers, headers)
	}
}

// WithHeader adds a single custom header to the response.
func WithHeader(key, value string) WriteOption {
	return func(o *WriteOptions) {
		if o.Headers == nil {
			o.Headers = make(map[string]string)
		}
		o.Headers[key] = value
	}
}
