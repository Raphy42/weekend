package telemetry

import (
	"context"
	"sync"
	"time"

	"github.com/palantir/stacktrace"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/trace"

	"github.com/Raphy42/weekend/core/scheduler"
	"github.com/Raphy42/weekend/core/scheduler/schedulable"
)

type Telemetry struct {
	sync.RWMutex
	tracer *trace.TracerProvider
}

func NewTelemetry(tracer *trace.TracerProvider) *Telemetry {
	return &Telemetry{tracer: tracer}
}

func (t *Telemetry) start() schedulable.Manifest {
	return schedulable.Of(
		schedulable.Name("wk.telemetry.start"),
		func() {
			otel.SetTracerProvider(t.tracer)
		},
	)
}

func (t *Telemetry) stop() schedulable.Manifest {
	return schedulable.Of(
		schedulable.Name("wk.telemetry.stop"),
		func(ctx context.Context) error {
			t.Lock()
			defer t.Unlock()

			shutdownCtx, cancel := context.WithTimeout(ctx, time.Second*5)
			defer cancel()

			if err := t.tracer.Shutdown(shutdownCtx); err != nil {
				return stacktrace.Propagate(err, "could not terminate tracer")
			}
			return nil
		})
}

func (t *Telemetry) Manifest(applicationContext context.Context) schedulable.Manifest {
	start := t.start()
	stop := t.stop()

	pipeline := scheduler.MakePipeline(
		schedulable.Name("wk.telemetry"),
		scheduler.Hooks{
			OnStart: &start,
			OnStop:  &stop,
		},
		scheduler.Lock(applicationContext),
	)
	return pipeline.Manifest()
}
