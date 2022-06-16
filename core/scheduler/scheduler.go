package scheduler

import (
	"context"
	"sync"

	"github.com/palantir/stacktrace"
	"github.com/rs/xid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"

	"github.com/Raphy42/weekend/core/errors"
	"github.com/Raphy42/weekend/core/logger"
	"github.com/Raphy42/weekend/core/message"
	"github.com/Raphy42/weekend/core/scheduler/async"
)

// todo handle the case where the bus is full

type Scheduler struct {
	sync.RWMutex
	id  xid.ID
	bus message.Bus
}

func New(bus message.Bus) *Scheduler {
	return &Scheduler{bus: bus, id: xid.New()}
}

func Schedule(ctx context.Context, manifest async.Manifest, args any) (*Handle, error) {
	ctx, span := otel.Tracer("wk.core.schedule").Start(ctx, "scheduleFromContext")
	span.SetAttributes(
		attribute.String("wk.manifest.name", manifest.Name),
		attribute.Stringer("wk.manifest.id", manifest.ID),
	)
	defer span.End()

	switch v := ctx.(type) {
	case *Context:
		if v == nil {
			return nil, stacktrace.NewError("invalid context")
		}
		return Schedule(*v, manifest, args)
	case Context:
		if v.Scheduler == nil {
			return nil, stacktrace.NewError("no scheduler was bound to this scheduling context")
		}
		return v.Scheduler.Schedule(ctx, manifest, args)
	default:
		schedulerI := ctx.Value(schedulerInjectionKey)
		if schedulerI == nil {
			return nil, stacktrace.NewError("unable to schedule, not a scheduling context, context is %T", v)
		}
		scheduler, ok := schedulerI.(*Scheduler)
		if !ok {
			return nil, stacktrace.NewError("invalid context value, expected *scheduler.Scheduler got %T", schedulerI)
		}
		return scheduler.Schedule(ctx, manifest, args)
	}
}

func (s *Scheduler) Schedule(parent context.Context, manifest async.Manifest, args any) (*Handle, error) {
	parent, span := otel.Tracer("wk.core.schedule").Start(parent, "schedule")
	span.SetAttributes(
		attribute.String("wk.manifest.name", manifest.Name),
		attribute.Stringer("wk.manifest.id", manifest.ID),
		attribute.Stringer("wk.scheduler.id", s.id),
	)
	defer span.End()

	handle, resultChan, errChan := NewHandle(parent, s.id, manifest)
	log := logger.FromContext(parent).With(
		zap.Stringer("wk.sched.id", handle.ID),
		zap.String("wk.sched.name", manifest.Name),
		zap.Stringer("wk.sched.parent.id", handle.Parent),
	)

	// children can now use `scheduler.Schedule` convenience method thanks to `scheduler.Context`
	handle.BindScheduler(s)
	s.Lock()
	defer s.Unlock()

	log.Debug("scheduling function")
	_ = s.bus.Emit(handle, NewScheduledMessage(manifest.Name, manifest.ID, handle.Parent))
	go func(ctx context.Context, resultC chan<- any, errC chan<- error, f async.Fn, in any, bus message.Bus) {
		ctx, goRoutineSpan := otel.Tracer("wk.core.schedule").Start(ctx, "goroutine")
		goRoutineSpan.SetAttributes(attribute.String("wk.manifest.name", manifest.Name))
		defer goRoutineSpan.End()

		// todo install telemetry
		defer errors.InstallPanicObserver()
		defer func() {
			_ = log.Sync()
		}()

		// welcome to the
		// actual function call
		result, err := f(ctx, in)

		if err != nil {
			goRoutineSpan.RecordError(err)
			_ = bus.Emit(ctx, NewFailureMessage(manifest.ID, handle.ID, err))
			select {
			case <-ctx.Done():
				log.Error("unexpected context termination, error was lost", zap.Error(err))
			case errC <- err:
				log.Debug("returned non-nil error", zap.Error(err))
				return
			}
		} else {
			_ = bus.Emit(ctx, NewSuccessMessage(manifest.ID, handle.ID))
			select {
			case <-ctx.Done():
				log.Error("unexpected context termination, result was lost", zap.Any("result", result))
			case resultC <- result:
				log.Debug("finished with success")
				return
			}
		}
	}(handle.Context, resultChan, errChan, manifest.Fn, args, s.bus)

	return handle, nil
}
