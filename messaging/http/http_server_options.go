package messaginghttp

import (
	"time"
)

type ServerConfig struct {
	ReadTimeout       time.Duration
	ReadHeaderTimeout time.Duration
	WriteTimeout      time.Duration
	IdleTimeout       time.Duration
}

type ServerOption func(*ServerConfig)

// WithReadTimeout sets the maximum duration for reading the entire request, including the body.
func WithReadTimeout(d time.Duration) ServerOption {
	return func(c *ServerConfig) {
		c.ReadTimeout = d
	}
}

// WithReadHeaderTimeout sets the amount of time allowed to read request headers.
func WithReadHeaderTimeout(d time.Duration) ServerOption {
	return func(c *ServerConfig) {
		c.ReadHeaderTimeout = d
	}
}

// WithWriteTimeout sets the maximum duration before timing out writes of the response.
func WithWriteTimeout(d time.Duration) ServerOption {
	return func(c *ServerConfig) {
		c.WriteTimeout = d
	}
}

// WithIdleTimeout sets the maximum amount of time to wait for the next request when keep-alives are enabled.
func WithIdleTimeout(d time.Duration) ServerOption {
	return func(c *ServerConfig) {
		c.IdleTimeout = d
	}
}
