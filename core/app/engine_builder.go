package app

import (
	"context"
	"time"

	"github.com/palantir/stacktrace"
	"go.uber.org/zap"

	"github.com/Raphy42/weekend/core/logger"
	"github.com/Raphy42/weekend/core/scheduler/async"
	"github.com/Raphy42/weekend/core/supervisor"
	"github.com/Raphy42/weekend/pkg/chrono"
	"github.com/Raphy42/weekend/pkg/reflect"
	"github.com/Raphy42/weekend/pkg/std/slice"
)

type EngineBuilder struct {
	background []async.Manifest
}

func NewEngineBuilder() *EngineBuilder {
	return &EngineBuilder{
		background: make([]async.Manifest, 0),
	}
}

func (e *EngineBuilder) Background(manifests ...async.Manifest) *EngineBuilder {
	e.background = append(e.background, manifests...)
	return e
}

func (e *EngineBuilder) HealthCheck(service any, interval time.Duration, fn func(ctx context.Context) error) *EngineBuilder {
	typename := reflect.Typename(service)
	e.Background(async.Of(
		async.Name("wk", typename, "health_check"),
		func(ctx context.Context) error {
			log := logger.FromContext(ctx).With(zap.String("service", typename))
			log.Info("healthcheck started")

			return <-chrono.NewTicker(interval).
				TickErr(ctx, func() error {
					err := fn(ctx)
					if err != nil {
						log.Error("health check failed", zap.Error(err))
					}
					return stacktrace.Propagate(
						err,
						"health_check failed for '%T'", service,
					)
				})
		},
	))
	return e
}

func (e *EngineBuilder) Build() (*Engine, error) {
	return &Engine{
		supervisor: supervisor.New("app_engine", slice.Map(
			e.background,
			func(_ int, in async.Manifest) supervisor.Spec {
				return supervisor.NewSpec(
					in,
					nil,
					// restart if error is not transient
					supervisor.WithRestartStrategy(supervisor.TransientRestartStrategy),
					// shutdown immediately
					supervisor.WithShutdownStrategy(supervisor.ImmediateShutdownStrategy),
					// restart only failed child
					supervisor.WithSupervisionStrategy(supervisor.OneForOneSupervisionStrategy),
				)
			},
		)...),
		manifests: e.background,
		errors:    nil,
	}, nil
}
