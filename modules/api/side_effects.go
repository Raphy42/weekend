package api

import (
	"context"

	"go.uber.org/zap"

	"github.com/Raphy42/weekend/core/logger"
	"github.com/Raphy42/weekend/pkg/reflect"
)

func debugRoutes(ctx context.Context, server *Server) error {
	log := logger.FromContext(ctx)
	for _, route := range server.engine.Routes() {
		log.Info(route.Path,
			zap.String("handler", reflect.Typename(route.Handler)),
			zap.String("method", route.Method),
		)
	}

	return nil
}
