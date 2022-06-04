package scheduler

import (
	"context"

	"github.com/palantir/stacktrace"
	"go.uber.org/zap"

	"github.com/Raphy42/weekend/core/logger"
	"github.com/Raphy42/weekend/core/scheduler/schedulable"
)

type Hooks struct {
	OnStart *schedulable.Manifest
	OnStop  *schedulable.Manifest
}

func DefaultHooks() Hooks {
	return Hooks{}
}

type Pipeline struct {
	Name      string
	Manifests []schedulable.Manifest
	Hooks     Hooks
}

func MakePipeline(name string, hooks Hooks, manifests ...schedulable.Manifest) *Pipeline {
	return &Pipeline{
		Name:      name,
		Manifests: manifests,
		Hooks:     hooks,
	}
}

func (p Pipeline) Manifest() schedulable.Manifest {
	return schedulable.Make(p.Name, func(ctx context.Context, input interface{}) (interface{}, error) {
		log := logger.FromContext(ctx).With(zap.String("wk.pipeline", p.Name))

		args := input

		if p.Hooks.OnStart != nil {
			handle, err := Schedule(ctx, *p.Hooks.OnStart, args)
			if err != nil {
				return nil, stacktrace.Propagate(err, "could not schedule OnStart hook")
			}
			log.Info("scheduled OnStart hook")
			args, err = handle.Poll(ctx)
			if err != nil {
				return nil, stacktrace.Propagate(err, "OnStart hook failed")
			}
		}

		for idx, manifest := range p.Manifests {
			handle, err := Schedule(ctx, manifest, args)
			if err != nil {
				return &idx, stacktrace.Propagate(
					err,
					"pipeline execution failed: '%s' %s", manifest.Name, manifest.ID,
				)
			}
			args, err = handle.Poll(ctx)
			if err != nil {
				return &idx, stacktrace.Propagate(
					err,
					"pipeline step returned non nil error: '%s' %s", manifest.Name, manifest.ID,
				)
			}
		}

		id := len(p.Manifests)
		if p.Hooks.OnStop != nil {
			handle, err := Schedule(ctx, *p.Hooks.OnStop, args)
			if err != nil {
				return &id, stacktrace.Propagate(err, "could not schedule OnStop hook")
			}
			log.Info("scheduled OnStop hook")
			args, err = handle.Poll(ctx)
			if err != nil {
				return &id, stacktrace.Propagate(err, "OnStop hook failed")
			}
		}

		return args, nil
	})
}
