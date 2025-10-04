package messaginghttp

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type CommandGINServer struct {
	cfg    ServerConfig
	h      *MessageHandler
	srv    *http.Server
	engine *gin.Engine
}

// NewCommandGINServer creates a new CommandHTTPServer with the given CommandBus and options.
func NewCommandGINServer(handler *CommandHandler, engine *gin.Engine, opts ...ServerOption) *CommandGINServer {
	cfg := new(ServerConfig)
	for _, opt := range opts {
		opt(cfg)
	}

	return &CommandGINServer{h: handler, cfg: *cfg, engine: engine}
}

func (s *CommandGINServer) ListenAndServe(addr string) error {
	s.engine.POST("/commands", gin.WrapH(s.h))

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

func (s *CommandGINServer) Close() error {
	if s.srv != nil {
		return s.srv.Close()
	}
	return nil
}
