package scheduler

import (
	"context"
	"sync"

	"github.com/palantir/stacktrace"
	"github.com/rs/xid"
	"go.uber.org/zap"

	"github.com/Raphy42/weekend/core/errors"
	"github.com/Raphy42/weekend/core/logger"
	"github.com/Raphy42/weekend/core/message"
	"github.com/Raphy42/weekend/core/scheduler/schedulable"
)

type Scheduler struct {
	sync.RWMutex
	id  xid.ID
	bus message.Bus
}

func New(bus message.Bus) *Scheduler {
	return &Scheduler{bus: bus, id: xid.New()}
}

func Schedule(ctx context.Context, manifest schedulable.Manifest, args interface{}) (*Handle, error) {
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
		return nil, stacktrace.NewError("unable to schedule, not a scheduling context")
	}
}

func (s *Scheduler) Schedule(parent context.Context, manifest schedulable.Manifest, args interface{}) (*Handle, error) {
	handle, resultChan, errChan := NewHandle(parent, s.id)
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
	_ = s.bus.Emit(handle, NewScheduledMessage(manifest.Name, handle.ID, handle.Parent))
	go func(ctx context.Context, resultC chan<- interface{}, errC chan<- error, f schedulable.Fn, in interface{}, bus message.Bus) {
		// todo install telemetry
		defer errors.InstallPanicObserver()
		defer func() {
			_ = log.Sync()
		}()

		// welcome to the
		// actual function call
		result, err := f(ctx, in)

		if err != nil {
			_ = bus.Emit(ctx, NewFailureMessage(handle.ID, err))
			select {
			case <-ctx.Done():
				log.Error("unexpected context termination, error was lost", zap.Error(err))
			case errC <- err:
				log.Debug("returned non-nil error", zap.Error(err))
				return
			}
		} else {
			_ = bus.Emit(ctx, NewSuccessMessage(handle.ID))
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
