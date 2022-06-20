package task

import (
	"context"
	"fmt"
	"time"

	"github.com/palantir/stacktrace"
	"github.com/rs/xid"
	"go.uber.org/zap"

	"github.com/Raphy42/weekend/core/logger"
	"github.com/Raphy42/weekend/core/message"
	"github.com/Raphy42/weekend/core/scheduler"
	"github.com/Raphy42/weekend/core/scheduler/async"
	"github.com/Raphy42/weekend/core/supervisor"
	"github.com/Raphy42/weekend/pkg/chrono"
	"github.com/Raphy42/weekend/pkg/concurrent_set"
	"github.com/Raphy42/weekend/pkg/set"
)

type Worker struct {
	id              xid.ID
	tasks           map[string]Task
	bus             message.Bus
	running         concurrent_set.Set[xid.ID, *scheduler.Future]
	announceMailbox message.Mailbox
}

func NewWorker(bus message.Bus, tasks ...Task) *Worker {
	return &Worker{
		id: xid.New(),
		tasks: set.From(tasks, func(item Task) (string, Task) {
			return item.Name, item
		}),
		bus:     bus,
		running: concurrent_set.New[xid.ID, *scheduler.Future](),
	}
}

func (w *Worker) Announce(ctx context.Context) error {
	if w.announceMailbox == nil {
		subject, err := w.bus.Subject(ctx, workerAnnounceSubject)
		if err != nil {
			return stacktrace.Propagate(err, "unable to open worker announce subject")
		}
		w.announceMailbox = subject
	}
	subject := w.announceMailbox

	jobNames := set.CollectSlice(w.tasks, func(k string, v Task) (string, bool) {
		return v.Name, true
	})
	if err := subject.Emit(ctx, NewWorkerAnnounceMessage(w.id, jobNames...)); err != nil {
		return stacktrace.Propagate(err, "could not announce worker to controller")
	}

	return nil
}

func (w *Worker) handleExecuteMessage(
	ctx context.Context,
	payload *ExecuteMessage,
	updateMailbox message.Mailbox,
) (*scheduler.Future, error) {
	task, ok := w.tasks[payload.Name]
	if !ok {
		return nil, stacktrace.NewError("no such manifest '%s' registered in this worker", payload.Name)
	}
	runner := supervisor.New(
		async.Name("wk", "task", payload.Name, "runner"),
		supervisor.NewSpec(task.AsyncManifest, &payload.Manifest,
			supervisor.WithSupervisionStrategy(supervisor.OneForOneSupervisionStrategy),
			supervisor.WithRestartStrategy(supervisor.TransientRestartStrategy),
			supervisor.WithShutdownStrategy(supervisor.ImmediateShutdownStrategy),
		),
	)
	runnerManifest := runner.Manifest()
	future, err := scheduler.Schedule(ctx, runnerManifest, nil)
	if err != nil {
		return nil, stacktrace.Propagate(err, "could not start task")
	}

	if err := updateMailbox.Emit(ctx, NewTaskExecutingMessage(task.AsyncManifest.ID, future.ID)); err != nil {
		return future, err
	}

	return future, nil
}

func (w *Worker) PollingExecutor(ctx context.Context) (*async.Manifest, error) {
	log := logger.FromContext(ctx).With(zap.Stringer("wk.worker.id", w.id))
	log.Info("worker starting")
	for _, manifest := range w.tasks {
		log.Debug("task registered",
			zap.String("wk.task.name", manifest.Name),
			zap.Stringer("wk.task.async.id", manifest.AsyncManifest.ID),
		)
	}

	updateMailbox, err := w.bus.Subject(ctx, workerUpdateSubject)
	if err != nil {
		return nil, stacktrace.Propagate(err, "unable to create worker mailbox")
	}

	executeRoutine := async.Of(
		async.Name("wk", "worker", w.id.String(), "execute"),
		func(ctx context.Context) error {
			mailbox, err := w.bus.Subject(ctx, fmt.Sprintf(pollExecuteSubjectFmt, w.id))
			if err != nil {
				return stacktrace.Propagate(err, "unable to create worker mailbox")
			}
			msgC, cancel := mailbox.ReadC(ctx)
			defer cancel()

			for msg := range msgC {
				switch msg.Kind {
				case MTaskExecute:
					payload := msg.Payload.(*ExecuteMessage)
					future, err := w.handleExecuteMessage(ctx, payload, mailbox)
					if err != nil {
						return stacktrace.Propagate(err, "unrecoverable error")
					}
					w.running.Insert(future.ID, future)
				}
			}
			return nil
		},
	)
	pollingRoutine := async.Of(
		async.Name("wk", "worker", w.id.String(), "poll"),
		func(ctx context.Context) error {
			// todo improve this bullshit
			ticker := chrono.NewTicker(time.Second * 2)
			errC := ticker.TickErr(ctx, func() error {
				for _, future := range w.running.Values() {
					result, done, err := future.TryPoll(ctx, time.Millisecond*10)
					if done {
						log = log.With(zap.Stringer("wk.future.id", future.ID))
						if err != nil {
							log.Error("task failed", zap.Error(err))
						} else {
							log.Info("task success")
						}
						if err := updateMailbox.Emit(ctx, NewTaskExecutedMessage(future.ID, result, err)); err != nil {
							return stacktrace.Propagate(err, "could not dispatch message")
						}
					}
				}
				return nil
			})
			return <-errC
		},
	)

	super := supervisor.New(
		async.Name("wk", "worker", w.id.String(), "executor"),
		supervisor.NewSpec(pollingRoutine, nil),
		supervisor.NewSpec(executeRoutine, nil),
	)
	manifest := super.Manifest()
	return &manifest, nil
}
