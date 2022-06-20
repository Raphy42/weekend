package api

import (
	"context"

	"github.com/palantir/stacktrace"
	"go.uber.org/zap"

	"github.com/Raphy42/weekend/core/logger"
	"github.com/Raphy42/weekend/pkg/reflect"
)

func startServer(ctx context.Context, server *Server) error {
	log := logger.FromContext(ctx)

	if err := server.markAsReady(ctx); err != nil {
		return stacktrace.Propagate(err, "could not start server")
	}

	log.Info("listing routes")
	for _, route := range server.engine.Routes() {
		log.Info(route.Path,
			zap.String("handler", reflect.Typename(route.Handler)),
			zap.String("method", route.Method),
		)
	}

	return nil
}
