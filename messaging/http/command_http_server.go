package messaginghttp

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type CommandServer struct {
	cfg    ServerConfig
	h      *MessageHandler
	srv    *http.Server
	engine *gin.Engine
}

// NewGINCommandServer creates a new CommandHTTPServer with the given CommandBus and options.
func NewGINCommandServer(handler *CommandHandler, engine *gin.Engine, opts ...ServerOption) *CommandServer {
	cfg := new(ServerConfig)
	for _, opt := range opts {
		opt(cfg)
	}

	return &CommandServer{h: handler, cfg: *cfg, engine: engine}
}

func (s *CommandServer) ListenAndServe(addr string) error {
	s.engine.POST("/commands", func(ctx *gin.Context) {
		s.h.ServeHTTP(ctx.Writer, ctx.Request)
	})

	s.srv = &http.Server{
		Addr:              addr,
		Handler:           s.engine,
		ReadTimeout:       s.cfg.ReadTimeout,
		ReadHeaderTimeout: s.cfg.ReadHeaderTimeout,
		WriteTimeout:      s.cfg.WriteTimeout,
		IdleTimeout:       s.cfg.IdleTimeout,
	}

	return s.srv.ListenAndServe()
}

func (s *CommandServer) Close() error {
	if s.srv != nil {
		return s.srv.Close()
	}
	return nil
}
