package core

import (
	"context"
	"runtime"

	"github.com/palantir/stacktrace"
	"go.uber.org/zap"

	"github.com/Raphy42/weekend/core/app"
	"github.com/Raphy42/weekend/core/logger"
	"github.com/Raphy42/weekend/core/service"
)

func platformInformation(ctx context.Context) error {
	log := logger.FromContext(ctx)

	log.Info("platform information",
		zap.String("os.architecture", runtime.GOARCH),
		zap.String("os.kernel", runtime.GOOS),
		zap.String("go.version", runtime.Version()),
		zap.Int("go.concurrency", runtime.GOMAXPROCS(runtime.NumCPU())),
	)

	return ctx.Err()
}

func applicationServiceHealthInjector(app *app.App, registry *service.Registry) {
	app.SetRegistry(registry)
}

func applicationEngineInjector(ctx context.Context, app *app.App, builder *app.EngineBuilder) error {
	log := logger.FromContext(ctx)

	engine, err := builder.Build()
	if err != nil {
		return stacktrace.Propagate(err, "could not build application engine")
	}
	log.Debug("attaching engine")
	return app.SetEngine(engine)
}
