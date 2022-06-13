package scheduler

import (
	"context"

	"github.com/palantir/stacktrace"
	"go.opentelemetry.io/otel"

	"github.com/Raphy42/weekend/core/scheduler/schedulable"
	"github.com/Raphy42/weekend/pkg/reflect"
)

type lock struct {
	ctx context.Context
}

func Lock(topCtx context.Context) schedulable.Manifest {
	return schedulable.Of(
		schedulable.Name("wk.context.lock", reflect.Typename(topCtx)),
		func(runningCtx context.Context) error {
			runningCtx, span := otel.Tracer("wk.scheduler").Start(runningCtx, "lock")
			defer span.End()

			select {
			case <-runningCtx.Done():
				return stacktrace.Propagate(runningCtx.Err(), "scheduling context was cancelled before lock context")
			case <-topCtx.Done():
				return nil
			}
		},
	)
}
