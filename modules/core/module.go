package core

import (
	"context"
	"runtime"

	"github.com/palantir/stacktrace"
	"go.uber.org/zap"

	"github.com/Raphy42/weekend/core/config"
	"github.com/Raphy42/weekend/core/dep"
	"github.com/Raphy42/weekend/core/logger"
)

var (
	ModuleName = dep.Name("wk", "platform")
)

func engineBuilderProvider() *EngineBuilder {
	return newEngineBuilder()
}

func applicationContextProvider(ctx context.Context) func() context.Context {
	return func() context.Context {
		return ctx
	}
}

func platformInformation(ctx context.Context) {
	log := logger.FromContext(ctx)

	log.Info("platform information",
		zap.String("os.architecture", runtime.GOARCH),
		zap.String("os.kernel", runtime.GOOS),
		zap.String("go.version", runtime.Version()),
		zap.Int("go.concurrency", runtime.GOMAXPROCS(runtime.NumCPU())),
	)
}

func configFromFilenamesProvider(filenames ...string) func(ctx context.Context) (*config.Config, error) {
	return func(ctx context.Context) (*config.Config, error) {
		cfg, err := configFromFilenames(ctx, filenames...)
		if err != nil {
			return nil, stacktrace.Propagate(err, "could not build application configuration")
		}
		return &config.Config{Configurable: cfg}, nil
	}
}

func Module(opts ...ModuleOption) dep.Module {
	options := defaultModuleOptions()
	options.apply(opts...)

	return dep.Declare(
		ModuleName,
		dep.Factories(
			engineBuilderProvider,
			applicationContextProvider(options.rootCtx),
			configFromFilenamesProvider(options.configFilenames...),
		),
		dep.SideEffects(platformInformation),
	)
}
