package app

import (
	"context"
	"time"

	"github.com/palantir/stacktrace"
	"go.uber.org/atomic"
	"go.uber.org/zap"

	"github.com/Raphy42/weekend/core/logger"
	"github.com/Raphy42/weekend/core/scheduler/async"
	"github.com/Raphy42/weekend/core/supervisor"
	"github.com/Raphy42/weekend/pkg/chrono"
	"github.com/Raphy42/weekend/pkg/reflect"
	"github.com/Raphy42/weekend/pkg/slice"
)

type EngineBuilder struct {
	done       *atomic.Bool
	background []async.Manifest
}

func NewEngineBuilder() *EngineBuilder {
	return &EngineBuilder{
		done:       atomic.NewBool(false),
		background: make([]async.Manifest, 0),
	}
}

func (e *EngineBuilder) Background(manifests ...async.Manifest) *EngineBuilder {
	if e.done.Load() {
		panic("trying to register a background routine in builder, but engine is already running")
	}
	e.background = append(e.background, manifests...)
	return e
}

func healthCheckImpl(
	typename string,
	interval time.Duration,
	service any,
	fn func(ctx context.Context) error) func(ctx context.Context) error {
	return func(ctx context.Context) error {
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
	}
}

func (e *EngineBuilder) HealthCheck(service any, interval time.Duration, fn func(ctx context.Context) error) *EngineBuilder {
	if e.done.Load() {
		panic("trying to register a health-check routine in builder, but engine is already running")
	}

	log := logger.New()

	typename := reflect.Typename(service)
	log.Debug("registered new health check", zap.String("wk.service.name", typename))
	manifest := async.Of(
		async.Name("wk", typename, "health_check"),
		healthCheckImpl(typename, interval, service, fn),
	)
	e.Background(manifest)
	return e
}

func (e *EngineBuilder) Build() (*Engine, error) {
	e.done.Store(true)

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
