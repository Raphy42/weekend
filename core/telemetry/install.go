package telemetry

import (
	"context"

	"go.opentelemetry.io/otel/trace"

	"github.com/Raphy42/weekend/core"
	"github.com/Raphy42/weekend/core/errors"
)

func Install(name string) trace.Tracer {
	core.SetName(name)
	tracer, err := NewJaegerTracer(context.Background())
	errors.Mustf(err, "unable to create jaeger otel traer")

	_ = NewTelemetry(tracer)
	return tracer.Tracer(name)
}
