package api

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"

	"github.com/Raphy42/weekend/core"
	"github.com/Raphy42/weekend/core/app"
	"github.com/Raphy42/weekend/core/logger"
	"github.com/Raphy42/weekend/core/scheduler/async"
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

func newServer(server *http.Server, builder *app.EngineBuilder) *Server {
	gin.SetMode(gin.ReleaseMode)

	engine := gin.New()
	engine.Use(otelgin.Middleware(core.Name()))
	server.Handler = engine

	s := &Server{
		engine: engine,
		server: server,
	}

	builder.Background(async.Of(
		async.Name("wk.api.serve"),
		s.listenAndServe,
	))
	builder.HealthCheck(s, time.Second*10, func(_ context.Context) error {
		return nil
	})

	return s
}

func (s *Server) Group(name string) *gin.RouterGroup {
	return s.engine.Group(name)
}
