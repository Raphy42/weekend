package scheduler

import (
	"context"
	"sync"

	"github.com/google/uuid"
	"github.com/palantir/stacktrace"
	"go.uber.org/zap"

	"github.com/Raphy42/weekend/core/logger"
)

type Scheduler struct {
	sync.RWMutex
	running map[uuid.UUID]*Handle
}

func New() *Scheduler {
	return &Scheduler{
		running: make(map[uuid.UUID]*Handle),
	}
}

func Schedule(ctx context.Context, manifest Manifest, args interface{}) (*Handle, error) {
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

func (s *Scheduler) Schedule(parent context.Context, manifest Manifest, args interface{}) (*Handle, error) {
	handle, resultChan, errChan := NewHandle(parent, manifest.ID)
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
	go func(ctx context.Context, resultC chan<- interface{}, errC chan<- error, f SchedulableFn, in interface{}) {
		defer func() {
			_ = log.Sync()
		}()

		result, err := f(ctx, in)
		if err != nil {
			select {
			case <-ctx.Done():
				log.Error("unexpected context termination, error was lost", zap.Error(err))
			case errC <- err:
				log.Debug("returned non-nil error", zap.Error(err))
				return
			}
		} else {
			select {
			case <-ctx.Done():
				log.Error("unexpected context termination, result was lost", zap.Any("result", result))
			case resultC <- result:
				log.Debug("finished with success")
				return
			}
		}
	}(handle.Context, resultChan, errChan, manifest.Fn, args)
	s.running[handle.ID] = handle

	return handle, nil
}
