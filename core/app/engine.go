package app

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.uber.org/zap"

	"github.com/Raphy42/weekend/core/logger"
	"github.com/Raphy42/weekend/core/scheduler"
	"github.com/Raphy42/weekend/core/scheduler/schedulable"
	"github.com/Raphy42/weekend/pkg/channel"
)

type Engine struct {
	manifests []schedulable.Manifest
	errors    chan error
}

func (e *Engine) Manifest() schedulable.Manifest {
	return schedulable.Of(
		schedulable.Name("wk.app.engine"),
		func(ctx context.Context) error {
			ctx, span := otel.Tracer("wk.app.engine").Start(ctx, "run")
			defer span.End()

			log := logger.FromContext(ctx)

			log.Debug("engine starting up")
			handles := make([]*scheduler.Handle, 0)
			for _, manifest := range e.manifests {
				log.Info("scheduling background job", zap.String("wk.manifest.name", manifest.Name))
				handle, err := scheduler.Schedule(ctx, manifest, nil)
				if err != nil {
					return err
				}
				handles = append(handles, handle)
			}
			result := scheduler.Coalesce(ctx, handles...)
			go func() {
				ctx, span := otel.Tracer("wk.app.engine").Start(ctx, "goroutine")
				defer span.End()

				for {
					select {
					case <-ctx.Done():
						log.Info("engine background job terminated")
						return
					case value := <-result:
						if value.Error != nil {
							_ = channel.Send(ctx, value.Error, e.errors)
						}
					}
				}
			}()
			return nil
		},
	)
}
