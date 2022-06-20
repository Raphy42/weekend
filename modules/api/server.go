package api

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"

	"github.com/Raphy42/weekend/core"
	"github.com/Raphy42/weekend/core/logger"
	"github.com/Raphy42/weekend/pkg/channel"
)

type Server struct {
	engine *gin.Engine
	server *http.Server
	ready  chan struct{}
}

func (s *Server) listenAndServe(ctx context.Context) error {
	log := logger.FromContext(ctx)

	log.Debug("waiting for lock on server start to complete")
	<-s.ready
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
		ready:  make(chan struct{}, 1),
	}

	return s
}

func (s *Server) markAsReady(ctx context.Context) error {
	return channel.Send(ctx, struct{}{}, s.ready)
}

func (s *Server) Group(name string) *gin.RouterGroup {
	group := s.engine.Group(name)
	return group
}
