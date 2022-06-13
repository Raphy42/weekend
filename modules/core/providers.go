package core

import (
	"context"

	"github.com/palantir/stacktrace"

	"github.com/Raphy42/weekend/core/app"
	"github.com/Raphy42/weekend/core/config"
	"github.com/Raphy42/weekend/core/service"
)

func engineBuilderProvider() *app.EngineBuilder {
	return app.NewEngineBuilder()
}

func applicationContextProvider(ctx context.Context) func() context.Context {
	return func() context.Context {
		return ctx
	}
}

func healthProvider() *service.Registry {
	return service.NewRegistry()
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
