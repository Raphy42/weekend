package telemetry

import (
	"sync"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/trace"

	"github.com/Raphy42/weekend/core/logger"
)

type Telemetry struct {
	sync.RWMutex
	tracer *trace.TracerProvider
}

func NewTelemetry(tracer *trace.TracerProvider) *Telemetry {
	otel.SetTracerProvider(tracer)
	logger.New().Info("tracer installed")
	return &Telemetry{tracer: tracer}
}
