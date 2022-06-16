package supervisor

import (
	"context"
	"sync"
	"time"

	"github.com/palantir/stacktrace"
	"github.com/rs/xid"
	"go.opentelemetry.io/otel"
	"go.uber.org/atomic"
	"go.uber.org/zap"

	"github.com/Raphy42/weekend/core/errors"
	"github.com/Raphy42/weekend/core/logger"
	"github.com/Raphy42/weekend/core/message"
	"github.com/Raphy42/weekend/core/scheduler"
	"github.com/Raphy42/weekend/core/scheduler/async"
	"github.com/Raphy42/weekend/pkg/std/set"
)

type Supervisor struct {
	lock      sync.RWMutex
	name      string
	bus       message.Bus
	scheduler *scheduler.Scheduler
	specLut   map[xid.ID]*Spec
	children  map[xid.ID]*scheduler.Handle
	restarts  map[xid.ID]*atomic.Int32
}

func New(name string, children ...Spec) *Supervisor {
	bus := message.NewInMemoryBus()
	sched := scheduler.New(bus)
	specLut := set.From(children, func(spec Spec) (xid.ID, *Spec) {
		return spec.Manifest.ID, &spec
	})

	return &Supervisor{
		name:      name,
		bus:       bus,
		scheduler: sched,
		specLut:   specLut,
		children:  make(map[xid.ID]*scheduler.Handle),
		restarts:  make(map[xid.ID]*atomic.Int32),
	}
}

func (s *Supervisor) registerHandle(handle *scheduler.Handle) {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.children[handle.ID] = handle
}

func (s *Supervisor) startChildren(ctx context.Context) error {
	ctx, span := otel.Tracer("wk.core.supervisor").Start(ctx, "startChildren")
	defer span.End()

	log := logger.FromContext(ctx)

	for _, spec := range s.specLut {
		log.Debug("starting supervision of child",
			zap.String("wk.supervisor.name", s.name),
			zap.String("wk.manifest.name", spec.Manifest.Name),
		)
		if err := s.bus.Emit(ctx, scheduler.NewScheduleMessage(spec.Manifest.ID, spec.Args)); err != nil {
			return stacktrace.Propagate(err, "unable to schedule child: %s", spec.Manifest.Name)
		}
	}
	return nil
}

func (s *Supervisor) terminateChildren(ctx context.Context) {
	s.lock.Lock()
	defer s.lock.Unlock()

	ctx, span := otel.Tracer("wk.core.supervisor").Start(ctx, "terminateChildren")
	defer span.End()

	log := logger.FromContext(ctx)
	log.Debug("terminating children", zap.Int("count", len(s.children)))

	for _, handle := range s.children {
		handle.Cancel()
	}
}

func (s *Supervisor) restart(ctx context.Context, id, handleID xid.ID, cause error) error {
	ctx, span := otel.Tracer("wk.core.supervisor").Start(ctx, "restart")
	defer span.End()

	log := logger.FromContext(ctx).With(zap.Stringer("wk.manifest.id", id))

	log.Info("trying to restart process")
	spec, ok := s.specLut[id]
	if !ok {
		return stacktrace.Propagate(cause, "retry: manifest '%s' has no associated spec", id)
	}
	// mark for restart if error is transient
	if errors.IsTransient(cause) {
		if err := s.bus.Emit(ctx, scheduler.NewScheduleMessage(spec.Manifest.ID, spec.Args)); err != nil {
			return stacktrace.Propagate(err, "retry: could not broadcast bus restart")
		}
	}

	counter := s.restarts[id].Inc()
	if counter > 3 {
		log.Error("maximum retry count exceeded")
		return stacktrace.Propagate(cause, "retry: maximum count exceeded")
	}

	strategy := spec.Strategy
	switch strategy.Supervision {
	case OneForOneSupervisionStrategy:
		log.Debug("re-scheduling")
		if err := s.bus.Emit(ctx, scheduler.NewScheduleMessage(spec.Manifest.ID, spec.Args)); err != nil {
			return stacktrace.Propagate(err, "retry: could not emit reschedule")
		}
	case OneForAllSupervisionStrategy:
		s.terminateChildren(ctx)
		log.Debug("re-scheduling all children")
		for _, spec := range s.specLut {
			log.Debug("re-scheduling child")
			if err := s.bus.Emit(ctx, scheduler.NewScheduleMessage(spec.Manifest.ID, spec.Args)); err != nil {
				return stacktrace.Propagate(err, "retry: could not emit reschedule")
			}
		}
	default:
		panic(stacktrace.NewErrorWithCode(errors.EUnreachable, "invalid supervision strategy"))
	}

	return nil
}

func (s *Supervisor) supervise(ctx context.Context) error {
	ctx, span := otel.Tracer("wk.core.supervisor").Start(ctx, "supervise")
	defer span.End()

	log := logger.FromContext(ctx).With(zap.String("wk.supervisor.name", s.name))
	log.Debug("supervisor started")

	messages, cancel := s.bus.Read(ctx)
	defer cancel()

	for {
		select {
		case <-ctx.Done():
			log.Info("supervisor terminating")
			shutdownCtx, cancelTmp := context.WithTimeout(context.Background(), time.Second*5)
			s.terminateChildren(shutdownCtx)
			cancelTmp()
			return nil
		case msg := <-messages:
			payload := msg.Payload
			switch msg.Kind {
			case scheduler.MSchedule:
				payload := payload.(*scheduler.ScheduleMessagePayload)
				if err := s.handleSchedule(ctx, payload); err != nil {
					return err
				}
			case scheduler.MFailure:
				payload := payload.(*scheduler.FailureMessagePayload)
				if err := s.handleFailure(ctx, payload); err != nil {
					return err
				}
			case scheduler.MProgress:
			case scheduler.MScheduled:
			case scheduler.MSuccess:
			}
		}
	}
}

func (s *Supervisor) stopChildren(ctx context.Context) error {
	ctx, span := otel.Tracer("wk.core.supervisor").Start(ctx, "stopChildren")
	defer span.End()

	log := logger.FromContext(ctx)

	for _, handle := range s.children {
		handle.Cancel()
	}
	log.Debug("terminating all supervised children")
	return nil
}

func (s *Supervisor) Manifest() async.Manifest {
	startChildren := async.Of(
		async.Name("wk", s.name, "supervisor", "start_children"),
		s.startChildren,
	)
	stopChildren := async.Of(
		async.Name("wk", s.name, "supervisor", "stop_children"),
		s.stopChildren,
	)
	superviseImpl := async.Of(
		async.Name("wk", s.name, "supervisor", "do"),
		s.supervise,
	)

	return scheduler.MakePipeline(
		async.Name("wk", s.name, "supervisor"),
		scheduler.Hooks{
			OnStart: &startChildren,
			OnStop:  &stopChildren,
		},
		superviseImpl,
	).
		Manifest()
}
