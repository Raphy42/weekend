package telemetry

import (
	"context"
	"os"

	"github.com/palantir/stacktrace"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.10.0"

	"github.com/Raphy42/weekend/core/app"
	"github.com/Raphy42/weekend/core/config"
	"github.com/Raphy42/weekend/core/logger"
)

const (
	ConfServiceName       = ".telemetry.name"
	ConfCollectorEndpoint = ".telemetry.endpoint"
)

func NewJaegerTracer(ctx context.Context, config *config.Config) (*trace.TracerProvider, error) {
	log := logger.FromContext(ctx)

	collectorUrl, err := config.URL(ctx, ConfCollectorEndpoint)
	if err != nil {
		return nil, stacktrace.Propagate(err, "missing configuration URL entry: '%s'", ConfCollectorEndpoint)
	}

	appName, err := config.String(ctx, ConfServiceName, os.Args[0])
	if err != nil {
		return nil, stacktrace.Propagate(err, "missing configuration string entry: '%s'", app.ConfApplicationName)
	}

	tracer, err := jaeger.New(
		jaeger.WithCollectorEndpoint(
			jaeger.WithEndpoint(collectorUrl.String()),
		),
	)
	if err != nil {
		return nil, stacktrace.Propagate(err, "could not create jaeger exporter")
	}

	log.Info("jaeger tracer provider instantiated")
	return trace.NewTracerProvider(
		trace.WithBatcher(tracer),
		trace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(appName),
		)),
	), nil
}
