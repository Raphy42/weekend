package api

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"

	"github.com/Raphy42/weekend/core"
	"github.com/Raphy42/weekend/core/logger"
)

type Server struct {
	engine *gin.Engine
	server *http.Server
}

func (s *Server) listenAndServe(ctx context.Context) error {
	log := logger.FromContext(ctx)

	log.Info("server is listening",
		zap.String("addr", s.server.Addr),
	)
	return s.server.ListenAndServe()
}

func newServer(server *http.Server) *Server {
	gin.SetMode(gin.ReleaseMode)

	engine := gin.New()
	engine.Use(otelgin.Middleware(core.Name()))
	server.Handler = engine

	s := &Server{
		engine: engine,
		server: server,
	}

	return s
}

func (s *Server) Group(name string) *gin.RouterGroup {
	group := s.engine.Group(name)
	return group
}
