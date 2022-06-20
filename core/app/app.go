package app

import (
	"context"
	"os"
	"os/signal"
	"sync"

	"github.com/palantir/stacktrace"
	"go.opentelemetry.io/otel"
	"go.uber.org/zap"

	"github.com/Raphy42/weekend/core"
	"github.com/Raphy42/weekend/core/dep"
	"github.com/Raphy42/weekend/core/errors"
	"github.com/Raphy42/weekend/core/logger"
	"github.com/Raphy42/weekend/core/scheduler"
	"github.com/Raphy42/weekend/core/service"
	"github.com/Raphy42/weekend/pkg/channel"
	"github.com/Raphy42/weekend/pkg/chrono"
)

//todo reduce internal dependencies bloat from the application object

type App struct {
	name         string
	lock         sync.RWMutex
	ctx          context.Context
	cancel       context.CancelFunc
	container    *dep.Container
	scheduler    *scheduler.Scheduler
	engineHandle *scheduler.Future
	healthProbe  <-chan service.Health
	engine       *Engine
}

func New(name string, opts ...BuilderOption) (*App, error) {
	core.SetName(name)
	builder := newBuilder(name)
	if err := builder.Apply(opts...); err != nil {
		return nil, stacktrace.Propagate(err, "could not build application")
	}
	app, err := builder.Build()
	if err != nil {
		return nil, err
	}
	app.container.Use(app.module())
	return app, nil
}

func (a *App) asInjectable() *App {
	return a
}

func (a *App) module() dep.Module {
	return dep.Declare(
		"wk.core.app",
		dep.Factories(
			a.asInjectable,
		),
	)
}

func (a *App) Start(ctx context.Context) error {
	timer := chrono.NewChrono()
	timer.Start()

	log := logger.FromContext(ctx).With(zap.String("wk.app", a.name))
	log.Info("starting application")

	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt, os.Kill)
	log.Debug("signal handlers installed")

	ctx, span := otel.Tracer("wk.core.app").Start(ctx, "App.Start")
	defer span.End()

	a.ctx = ctx
	a.cancel = cancel

	if err := core.ExecuteOnStart(); err != nil {
		return stacktrace.Propagate(err, "global application stopping hooks could not be run successfully")
	}

	a.lock.Lock()
	defer a.lock.Unlock()

	span.AddEvent("wk.lifecycle.on_start")
	if err := core.ExecuteOnStart(); err != nil {
		return stacktrace.Propagate(err, "global application starting hooks could not be run successfully")
	}
	log.Debug("global start hooks ran successfully", zap.Duration("wk.init.duration", timer.Elapsed()))

	handle, err := a.scheduler.Schedule(ctx, a.container.Manifest(), nil)
	span.AddEvent("wk.container.scheduled")
	if err != nil {
		return stacktrace.Propagate(err, "unable to schedule container")
	}

	span.AddEvent("wk.handle.polling")
	_, err = handle.Poll(ctx)
	if err != nil {
		return stacktrace.Propagate(err, "container bootstrap returned non-nil error")
	}

	a.engineHandle, err = a.scheduler.Schedule(ctx, a.engine.Manifest(), nil)
	span.AddEvent("wk.engine.scheduled")
	if err != nil {
		return stacktrace.Propagate(err, "unable to schedule engine")
	}

	log.Info("application initialised", zap.Duration("wk.init.duration", timer.Elapsed()))

	return nil
}

func (a *App) Wait(ctx context.Context) error {
	ctx, span := otel.Tracer("wk.core.app").Start(ctx, "App.Wait")
	defer span.End()

	result := make(chan error)
	if a.engine == nil {
		errors.Must(channel.Send(ctx, stacktrace.NewError("application is missing `core.Module()`"), result))
	}

	if a.ctx == nil {
		panic(stacktrace.NewError("application has not been started prior to polling (missing `App.Start`)"))
	}

	go func() {
		select {
		case health := <-a.healthProbe:
			if health.Error != nil {
				errors.Must(channel.Send(ctx, health.Error, result))
			}
		case err := <-a.engineHandle.Error():
			errors.Must(channel.Send(ctx, err, result))
		case <-a.ctx.Done():
			errors.Must(channel.Send(ctx, nil, result))
		case err := <-a.engine.errors:
			errors.Must(channel.Send(ctx, err, result))
		}
	}()

	//todo actually execute this synchronously
	defer a.engineHandle.Cancel()
	// todo fix span being dropped before this can be invoked
	defer func() {
		span.RecordError(core.ExecuteOnStop())
	}()

	return <-result
}

func (a *App) SetEngine(engine *Engine) error {
	a.engine = engine
	return nil
}

func (a *App) SetRegistry(registry *service.Registry) {
	a.healthProbe = registry.HealthProbe()
}
