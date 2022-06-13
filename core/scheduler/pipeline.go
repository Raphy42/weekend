package scheduler

import (
	"context"

	"github.com/palantir/stacktrace"
	"go.opentelemetry.io/otel"
	"go.uber.org/zap"

	"github.com/Raphy42/weekend/core/logger"
	"github.com/Raphy42/weekend/core/scheduler/async"
)

type Hooks struct {
	OnStart *async.Manifest
	OnStop  *async.Manifest
}

type Pipeline struct {
	Name      string
	Manifests []async.Manifest
	Hooks     Hooks
}

func MakePipeline(name string, hooks Hooks, manifests ...async.Manifest) *Pipeline {
	return &Pipeline{
		Name:      name,
		Manifests: manifests,
		Hooks:     hooks,
	}
}

func (p Pipeline) Manifest() async.Manifest {
	return async.Of(p.Name, func(ctx context.Context, input any) (any, error) {
		ctx, span := otel.Tracer("wk.pipeline").Start(ctx, p.Name)
		defer span.End()

		log := logger.FromContext(ctx).With(zap.String("wk.pipeline", p.Name))

		args := input

		log.Debug("executing OnStart hooks")
		if p.Hooks.OnStart != nil {
			handle, err := Schedule(ctx, *p.Hooks.OnStart, args)
			span.RecordError(err)
			if err != nil {
				return nil, stacktrace.Propagate(err, "could not schedule OnStart hook")
			}

			log.Info("scheduled OnStart hook")
			args, err = handle.Poll(ctx)
			span.RecordError(err)
			if err != nil {
				return nil, stacktrace.Propagate(err, "OnStart hook failed")
			}
		}

		for idx, manifest := range p.Manifests {
			handle, err := Schedule(ctx, manifest, args)
			span.RecordError(err)
			if err != nil {
				return &idx, stacktrace.Propagate(
					err,
					"pipeline execution failed: '%s' %s", manifest.Name, manifest.ID,
				)
			}

			args, err = handle.Poll(ctx)
			span.RecordError(err)
			if err != nil {
				return &idx, stacktrace.Propagate(
					err,
					"pipeline step returned non nil error: '%s' %s", manifest.Name, manifest.ID,
				)
			}
		}

		id := len(p.Manifests)
		log.Debug("executing OnStop hooks")
		if p.Hooks.OnStop != nil {
			handle, err := Schedule(ctx, *p.Hooks.OnStop, args)
			span.RecordError(err)
			if err != nil {
				return &id, stacktrace.Propagate(err, "could not schedule OnStop hook")
			}

			log.Info("scheduled OnStop hook")
			args, err = handle.Poll(ctx)
			span.RecordError(err)
			if err != nil {
				return &id, stacktrace.Propagate(err, "OnStop hook failed")
			}
		}

		return args, nil
	})
}
