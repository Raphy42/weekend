package telemetry

import (
	"context"

	"github.com/palantir/stacktrace"
	"go.opentelemetry.io/otel/sdk/trace"

	"github.com/Raphy42/weekend/core/config"
	"github.com/Raphy42/weekend/core/di"
	"github.com/Raphy42/weekend/core/errors"
	"github.com/Raphy42/weekend/core/telemetry"
	"github.com/Raphy42/weekend/modules/core"
)

var (
	ModuleName = di.Name("wk", "telemetry")
)

func telemetryProvider(
	ctx context.Context,
	trace *trace.TracerProvider,
	builder *core.EngineBuilder,
) (*telemetry.Telemetry, error) {
	if trace == nil {
		return nil, stacktrace.NewErrorWithCode(
			errors.ENil,
			"OTEL tracer implementation provided is not valid (nil)",
		)
	}
	t := telemetry.NewTelemetry(trace)
	builder.Background(t.Manifest(ctx))
	return t, nil
}

func tracerProvider(ctx context.Context, config config.Config) (*trace.TracerProvider, error) {
	return telemetry.NewJaegerTracer(ctx, config)
}

func Module() di.Module {
	return di.Declare(
		ModuleName,
		di.Providers(
			telemetryProvider,
			tracerProvider,
		),
	)
}
