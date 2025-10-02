package apix

import (
	"context"
	"net/http"
)

// HTTPRequestValidator allows validating http requests.
type HTTPRequestValidator interface {
	// Validate returns nil if the request is valid, or a Problem describing the issue.
	Validate(ctx context.Context, r *http.Request) *Problem
}
