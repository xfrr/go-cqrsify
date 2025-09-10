package retry

import (
	"net/http"
	"strconv"
	"strings"
	"time"
)

// HTTPError wraps an HTTP error with headers for hint extraction.
type HTTPError struct {
	StatusCode int
	Header     http.Header
	Err        error
	// ReceivedAt is when response was received (for date-based Retry-After)
	ReceivedAt time.Time
}

func (e *HTTPError) Error() string {
	return "http error: " + strconv.Itoa(e.StatusCode) + ": " + e.Err.Error()
}
func (e *HTTPError) Unwrap() error { return e.Err }

// Implement RetryAfterHint for convenience.
func (e *HTTPError) RetryAfter() (time.Duration, bool) {
	return parseRetryAfterHeader(e.Header.Get("Retry-After"), e.ReceivedAt)
}

// HTTPErrorFromResponse is a convenience to wrap an *http.Response you already have.
// Pass the underlying error you returned (e.g., io.ReadAll error or a sentinel).
func HTTPErrorFromResponse(resp *http.Response, underlying error) *HTTPError {
	h := http.Header{}
	for k, v := range resp.Header {
		// clone
		cp := append([]string(nil), v...)
		h[k] = cp
	}
	return &HTTPError{
		StatusCode: resp.StatusCode,
		Header:     h,
		Err:        underlying,
		ReceivedAt: time.Now(),
	}
}

// parseRetryAfterHeader supports:
//   - delta-seconds     -> "120"
//   - HTTP-date (RFC7231) -> "Wed, 21 Oct 2015 07:28:00 GMT"
func parseRetryAfterHeader(val string, now time.Time) (time.Duration, bool) {
	val = strings.TrimSpace(val)
	if val == "" {
		return 0, false
	}
	// Try delta-seconds
	if secs, err := strconv.Atoi(val); err == nil && secs >= 0 {
		return time.Duration(secs) * time.Second, true
	}
	// Try HTTP-date
	if t, err := http.ParseTime(val); err == nil {
		if now.IsZero() {
			now = time.Now()
		}
		if t.After(now) {
			return t.Sub(now), true
		}
		return 0, true // date is in the past -> 0s wait
	}
	return 0, false
}
