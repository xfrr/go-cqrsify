package messaginghttp

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const httpMethodQuery = "QUERY"

type QueryGINServer struct {
	cfg    ServerConfig
	h      *QueryHandler
	srv    *http.Server
	engine *gin.Engine
}

// NewQueryGINServer creates a new QueryGINServer with the given QueryHandler and options.
func NewQueryGINServer(handler *QueryHandler, engine *gin.Engine, opts ...ServerOption) *QueryGINServer {
	cfg := new(ServerConfig)
	for _, opt := range opts {
		opt(cfg)
	}

	return &QueryGINServer{h: handler, cfg: *cfg, engine: engine}
}

func (s *QueryGINServer) ListenAndServe(addr string) error {
	handler := gin.WrapH(s.h)
	s.engine.Handle(http.MethodPost, "/queries", handler)
	s.engine.Handle(httpMethodQuery, "/queries", handler)

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

func (s *QueryGINServer) Close() error {
	if s.srv != nil {
		return s.srv.Close()
	}
	return nil
}
