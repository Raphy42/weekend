package app

import (
	"context"
	"os"
	"os/signal"
	"sync"

	"github.com/palantir/stacktrace"
	"go.uber.org/zap"

	"github.com/Raphy42/weekend/core/dep"
	"github.com/Raphy42/weekend/core/errors"
	"github.com/Raphy42/weekend/core/logger"
	"github.com/Raphy42/weekend/core/scheduler"
	"github.com/Raphy42/weekend/pkg/chrono"
)

type App struct {
	name      string
	lock      sync.RWMutex
	ctx       context.Context
	cancel    context.CancelFunc
	container *dep.Container
	scheduler *scheduler.Scheduler
}

func New(name string, opts ...BuilderOption) (*App, error) {
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

func (a *App) Start(rootCtx context.Context) error {
	log := logger.FromContext(rootCtx).With(zap.String("wk.app", a.name))
	log.Info("starting application")

	timer := chrono.NewChrono()
	timer.Start()

	a.lock.Lock()
	defer a.lock.Unlock()

	if err := errors.ValidateContext(rootCtx); err != nil {
		return err
	}

	ctx, cancel := signal.NotifyContext(rootCtx, os.Interrupt, os.Kill)
	log.Debug("signal handlers installed")
	a.ctx = ctx
	a.cancel = cancel

	if err := executeOnStart(); err != nil {
		return stacktrace.Propagate(err, "global application starting hooks could not be run successfully")
	}
	log.Debug("global start hooks ran successfully", zap.Duration("wk.init.duration", timer.Elapsed()))

	handle, err := a.scheduler.Schedule(ctx, a.container.Manifest(), nil)
	if err != nil {
		return stacktrace.Propagate(err, "unable to schedule container")
	}
	_, err = handle.Poll(ctx)
	if err != nil {
		return stacktrace.Propagate(err, "container bootstrap returned non-nil error")
	}
	log.Info("application initialised", zap.Duration("wk.init.duration", timer.Elapsed()))

	return nil
}

func (a *App) Wait() <-chan struct{} {
	defer func() {
		if err := executeOnStart(); err != nil {
			panic(stacktrace.Propagate(err, "global application stopping hooks could not be run successfully"))
		}
	}()

	if a.ctx == nil {
		panic(stacktrace.NewError("application has not been started prior to polling (missing `App.Start`)"))
	}
	return a.ctx.Done()
}
