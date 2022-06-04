package scheduler

import (
	"context"
	"fmt"

	"github.com/palantir/stacktrace"
	"github.com/rs/xid"

	"github.com/Raphy42/weekend/core/logger"
	"github.com/Raphy42/weekend/core/scheduler/schedulable"
)

type Supervisor struct {
	id       xid.ID
	children []schedulable.Manifest
}

func NewSupervisor(children ...schedulable.Manifest) *Supervisor {
	return &Supervisor{
		id:       xid.New(),
		children: children,
	}
}

func (s Supervisor) onStart() schedulable.Manifest {
	return schedulable.Make(
		fmt.Sprintf("supervisor.%s.start", s.id.String()),
		func(ctx context.Context) {
			log := logger.FromContext(ctx)

			log.Info("START OF SUPERVISION")
		},
	)
}

func (s Supervisor) onStop() schedulable.Manifest {
	return schedulable.Make(
		fmt.Sprintf("supervisor.%s.start", s.id.String()),
		func(ctx context.Context) {
			log := logger.FromContext(ctx)

			log.Info("END OF SUPERVISION")
		},
	)
}

func (s Supervisor) Manifest() schedulable.Manifest {
	onStart := s.onStart()
	onStop := s.onStop()
	pipeline := MakePipeline(
		fmt.Sprintf("supervisor.%s.pipeline", s.id.String()),
		Hooks{
			OnStart: &onStart,
			OnStop:  &onStop,
		},
		s.children...,
	)
	name := fmt.Sprintf("supervisor.%s.instance", s.id.String())
	return schedulable.Make(
		name,
		func(ctx context.Context, args interface{}) (interface{}, error) {
			bus, err := busFromContext(ctx)
			if err != nil {
				return nil, err
			}

			scheduler := New(bus)
			handle, err := scheduler.Schedule(ctx, pipeline.Manifest(), args)
			_ = bus.Emit(ctx, NewSupervisedMessage(name, handle.ID, handle.Parent, s.id))

			if err != nil {
				err = stacktrace.Propagate(err, "unable to start schedule supervisor pipeline")
				_ = bus.Emit(ctx, NewSupervisionFailureMessage(handle.ID, s.id, err))
				return nil, err
			}
			result, err := handle.Poll(ctx)
			if err != nil {
				err = stacktrace.Propagate(err, "pipeline execution returned non-nil error")
				_ = bus.Emit(ctx, NewSupervisionFailureMessage(handle.ID, s.id, err))
				return nil, err
			}

			_ = bus.Emit(ctx, NewSupervisionSuccessMessage(handle.ID, s.id))
			return result, nil
		},
	)
}
