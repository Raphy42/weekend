package scheduler

import (
	"context"

	"github.com/palantir/stacktrace"
	"go.uber.org/atomic"
	"go.uber.org/zap"

	"github.com/Raphy42/weekend/core/errors"
	"github.com/Raphy42/weekend/core/kernel"
	"github.com/Raphy42/weekend/core/logger"
	"github.com/Raphy42/weekend/core/reflect"
)

type Runner struct {
	opts  []kernel.Option
	slots map[string]*atomic.Int32
}

func NewRunner(opts ...kernel.Option) *Runner {
	return &Runner{
		opts:  opts,
		slots: make(map[string]*atomic.Int32),
	}
}

func (s *Runner) hydrateHandle(handle *kernel.Handle) {
	handle.Metadata.Scheduler = reflect.Typename(s)
	handle.CompletionHooks = append(handle.CompletionHooks, func() {
		slot, ok := s.slots[handle.Name]
		if !ok {
			panic(
				stacktrace.NewErrorWithCode(
					errors.EUnreachable,
					"atomic counter was not found, should be initialised by the time this is called",
				))
		}
		slot.Dec()
	})
}

func (r *Runner) scheduleImpl(parentCtx context.Context, schedulable kernel.Schedulable, options kernel.Options, args interface{}) (*kernel.Handle, error) {
	log := logger.FromContext(parentCtx)

	if err := errors.ValidateContext(parentCtx); err != nil {
		return nil, stacktrace.Propagate(err, "can not schedule anything without a valid context")
	}

	handle, resultChan, errorChan := kernel.NewHandle(options, kernel.NewLiteralArgs(args), parentCtx)
	options.Apply(r.opts...)
	r.hydrateHandle(handle)
	log = log.With(zap.String("kernel.handle.id", handle.ID.String()))
	incrementCounter := func() {
		_, ok := r.slots[handle.Name]
		if !ok {
			// this is totally wonky AF, might need some RWLock to prevent races
			// todo fix race condition
			r.slots[handle.Name] = atomic.NewInt32(1)
		} else {
			r.slots[handle.Name].Inc()
		}
	}

	go func(s kernel.Schedulable, h *kernel.Handle, results chan<- interface{}, errors chan<- error) {
		h.MarkAsScheduled(reflect.Typename(s))
		log.Debug("scheduling started")
		incrementCounter()
		result, err := s.Execute(h.Context, h.Metadata.Args)
		if err != nil {
			log.Error("execution failure")
			errors <- err
		} else {
			log.Info("execution success")
			results <- result
		}
	}(schedulable, handle, resultChan, errorChan)

	return handle, nil
}

func (s *Runner) Schedule(ctx context.Context, schedulable kernel.Schedulable, args interface{}) (*kernel.Handle, error) {
	processOpts := schedulable.Options()
	processOpts.Apply(s.opts...)
	return s.scheduleImpl(ctx, schedulable, processOpts, args)
}
