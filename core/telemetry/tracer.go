package telemetry

import (
	"context"
	"runtime"

	"github.com/palantir/stacktrace"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.10.0"

	"github.com/Raphy42/weekend/core"
	"github.com/Raphy42/weekend/core/logger"
)

func NewJaegerTracer(ctx context.Context) (*trace.TracerProvider, error) {
	log := logger.FromContext(ctx)

	tracer, err := jaeger.New(
		jaeger.WithCollectorEndpoint(),
	)
	if err != nil {
		return nil, stacktrace.Propagate(err, "could not create jaeger exporter")
	}

	log.Info("jaeger tracer provider instantiated")
	return trace.NewTracerProvider(
		trace.WithBatcher(tracer),
		trace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(core.Name()),
			semconv.ServiceVersionKey.String(runtime.Version()),
		)),
		trace.WithSampler(trace.AlwaysSample()),
	), nil
}
