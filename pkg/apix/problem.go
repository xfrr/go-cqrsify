package apix

import (
	"encoding/json"
	"net/http"
)

// Problem represents RFC 7807 "Problem Details for HTTP APIs".
// See: https://datatracker.ietf.org/doc/html/rfc7807
type Problem struct {
	Type       string         `json:"type"`             // absolute URI identifying the problem type
	Title      string         `json:"title"`            // short, human-readable summary
	Status     int            `json:"status"`           // HTTP status
	Detail     string         `json:"detail,omitempty"` // human-readable explanation
	Instance   string         `json:"instance,omitempty"`
	Extensions map[string]any `json:"-"`
}

// MarshalJSON merges base fields with extensions (RFC allows custom members).
func (p Problem) MarshalJSON() ([]byte, error) {
	type base Problem
	m := map[string]any{}
	b := base(p)

	// Encode base first
	raw, err := json.Marshal(b)
	if err != nil {
		return nil, err
	}

	// Unmarshal to map for merge
	if err = json.Unmarshal(raw, &m); err != nil {
		return nil, err
	}

	// Merge extensions (if any), without clobbering standard keys
	for k, v := range p.Extensions {
		if _, exists := m[k]; !exists {
			m[k] = v
		}
	}
	return json.Marshal(m)
}

// NewProblem builds a minimal Problem with optional extensions.
func NewProblem(status int, title, detail string, opts ...func(*Problem)) Problem {
	p := Problem{
		Type:   "about:blank", // RFC default when you don't have a typed URI
		Title:  title,
		Status: status,
		Detail: detail,
	}
	for _, o := range opts {
		o(&p)
	}
	return p
}

// WithType sets the problem type URI.
func WithType(uri string) func(*Problem) {
	return func(p *Problem) { p.Type = uri }
}

// WithInstance sets the specific occurrence URI/path.
func WithInstance(inst string) func(*Problem) {
	return func(p *Problem) { p.Instance = inst }
}

// WithExtensions attaches custom members (e.g., field errors).
func WithExtensions(ext map[string]any) func(*Problem) {
	return func(p *Problem) { p.Extensions = ext }
}

// WriteProblem writes a Problem to the response with the correct content-type.
func WriteProblem(w http.ResponseWriter, p Problem) {
	WriteJSON(w, p,
		WithContentType(ContentTypeProblemJSON),
		WithStatus(p.Status),
		WithCacheControl("no-store"),
	)
}

// NewConflictProblem is a 409 Conflict helper.
func NewConflictProblem(detail string, opts ...func(*Problem)) Problem {
	return NewProblem(http.StatusConflict, "Conflict", detail, opts...)
}

// NewBadRequestProblem is a 400 Bad Request helper.
func NewBadRequestProblem(detail string, opts ...func(*Problem)) Problem {
	return NewProblem(http.StatusBadRequest, "Bad Request", detail, opts...)
}

// NewNotFoundProblem is a 404 Not Found helper.
func NewNotFoundProblem(detail string, opts ...func(*Problem)) Problem {
	return NewProblem(http.StatusNotFound, "Not Found", detail, opts...)
}

// NewUnauthorizedProblem is a 401 Unauthorized helper.
func NewUnauthorizedProblem(detail string, opts ...func(*Problem)) Problem {
	return NewProblem(http.StatusUnauthorized, "Unauthorized", detail, opts...)
}

// NewForbiddenProblem is a 403 Forbidden helper.
func NewForbiddenProblem(detail string, opts ...func(*Problem)) Problem {
	return NewProblem(http.StatusForbidden, "Forbidden", detail, opts...)
}

// NewInternalServerErrorProblem is a 500 Internal Server Error helper.
func NewInternalServerErrorProblem(detail string, opts ...func(*Problem)) Problem {
	return NewProblem(http.StatusInternalServerError, "Internal Server Error", detail, opts...)
}

// NewUnsupportedMediaTypeProblem is a 415 Unsupported Media Type helper.
func NewUnsupportedMediaTypeProblem(detail string, opts ...func(*Problem)) Problem {
	return NewProblem(http.StatusUnsupportedMediaType, "Unsupported Media Type", detail, opts...)
}
