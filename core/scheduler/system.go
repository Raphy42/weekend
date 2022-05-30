package scheduler

import (
	"context"

	"github.com/google/uuid"
	"github.com/palantir/stacktrace"
	"go.uber.org/atomic"
	"go.uber.org/zap"

	"github.com/Raphy42/weekend/core/errors"
	"github.com/Raphy42/weekend/core/kernel"
	"github.com/Raphy42/weekend/core/logger"
)

const (
	EventRunTask       = "task.run"
	EventCancelTask    = "task.cancel"
	EventTaskScheduled = "task#scheduled"
	EventTaskStarted   = "task#started"
	EventTaskSuccess   = "task#success"
	EventTaskFailure   = "task#failure"
)

type HandleMap map[uuid.UUID]kernel.HandleGroup

type Slots map[string]*atomic.Int32

type System struct {
	options      map[string]kernel.Options
	schedulables map[string][]kernel.Schedulable
	handles      map[string]HandleMap
	runner       kernel.Scheduler
	slots        Slots
}

func NewSystem(scheds ...kernel.Schedulable) *System {
	log := logger.New()

	options := make(map[string]kernel.Options)
	handles := make(map[string]HandleMap)
	schedulables := make(map[string][]kernel.Schedulable)
	runner := NewRunner()
	slots := make(Slots)
	for _, schedulable := range scheds {
		opts := schedulable.Options()
		name := opts.Name()
		options[name] = opts
		handles[name] = make(HandleMap)
		_, ok := schedulables[name]
		if ok {
			log.Warn("adding another schedulable to an already occupied slot", zap.String("kernel.schedulable.name", name))
		} else {
			schedulables[name] = make([]kernel.Schedulable, 0)
		}
		slots[name] = atomic.NewInt32(0)

		log.Debug("new schedulable registered", zap.String("kernel.schedulable.name", name))
		schedulables[name] = append(schedulables[name], schedulable)
	}

	return &System{
		options:      options,
		handles:      handles,
		runner:       runner,
		schedulables: schedulables,
		slots:        slots,
	}
}

func (s *System) Next(ctx context.Context, event *kernel.Event) (*kernel.Event, error) {
	log := logger.FromContext(ctx).With(zap.String("kernel.event.id", event.ID.String()))
	switch event.Kind {
	case EventRunTask:
		log.Debug("running task")
		event, err := s.runTask(ctx, event)
		if err != nil {
			return nil, stacktrace.Propagate(err, "could not run task")
		}
		return event, nil
	case EventCancelTask:
		log.Debug("cancelling task")
		event, err := s.cancelTask(ctx, event)
		if err != nil {
			return nil, stacktrace.Propagate(err, "could not cancel task")
		}
		return event, nil
	default:
		return nil, errors.NotImplemented("not sure about where to go about unknown events")
	}
}

func (s *System) cancelTask(ctx context.Context, event *kernel.Event) (*kernel.Event, error) {
	panic("not implemented")
}

type RunTaskEventPayload struct {
	Name    string
	Payload interface{}
}

type TaskScheduledEventPayload struct {
	ID uuid.UUID
}

func (s *System) runTask(ctx context.Context, event *kernel.Event) (*kernel.Event, error) {
	log := logger.FromContext(ctx)

	if event == nil {
		log.Warn("a nil event was received, ignoring")
		return nil, ctx.Err()
	}
	payload, ok := event.Payload.(*RunTaskEventPayload)
	if !ok {
		return nil, errors.InvalidCast(&RunTaskEventPayload{}, event.Payload)
	}
	schedulables, exists := s.schedulables[payload.Name]
	if !exists {
		return nil, errors.NotImplemented("needs to create error type")
	}

	handleGroup := make(kernel.HandleGroup, 0)
	for _, schedulable := range schedulables {
		handle, err := s.runner.Schedule(ctx, schedulable, payload.Payload)
		if err != nil {
			return nil, errors.NotImplemented("need to handle scheduling error policies here")
		}
		handleGroup = append(handleGroup, handle)
	}
	_, exists = s.handles[payload.Name]
	if !exists {
		s.handles[payload.Name] = make(HandleMap)
	}
	s.handles[payload.Name][event.ID] = handleGroup
	return kernel.NewEvent(EventTaskScheduled, &TaskScheduledEventPayload{ID: event.ID}), nil
}
