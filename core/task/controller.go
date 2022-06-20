package task

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/palantir/stacktrace"
	"go.uber.org/atomic"
	"go.uber.org/zap"

	"github.com/Raphy42/weekend/core/logger"
	"github.com/Raphy42/weekend/core/message"
	"github.com/Raphy42/weekend/core/repository/kv"
	"github.com/Raphy42/weekend/core/scheduler/async"
	"github.com/Raphy42/weekend/pkg/channel"
	"github.com/Raphy42/weekend/pkg/concurrent_set"
	"github.com/Raphy42/weekend/pkg/set"
	"github.com/Raphy42/weekend/pkg/slice"
)

type WorkerLoad struct {
	ID    string
	Tasks map[string]*atomic.Int32
}

type Controller struct {
	bus           message.Bus
	kv            kv.KV
	workers       concurrent_set.Set[string, message.Mailbox]
	lastUpdate    concurrent_set.Set[string, time.Time]
	dispatch      concurrent_set.Set[string, []string]
	load          concurrent_set.Set[string, WorkerLoad]
	dispatchQueue chan Manifest
}

func NewController(bus message.Bus, kv kv.KV) *Controller {
	return &Controller{
		bus:           bus,
		workers:       concurrent_set.New[string, message.Mailbox](),
		dispatch:      concurrent_set.New[string, []string](),
		lastUpdate:    concurrent_set.New[string, time.Time](),
		load:          concurrent_set.New[string, WorkerLoad](),
		kv:            kv,
		dispatchQueue: make(chan Manifest, 256),
	}
}

func (c *Controller) Dispatch(ctx context.Context, manifest Manifest) error {
	return channel.Send(ctx, manifest, c.dispatchQueue)
}

type WorkerInfo struct {
	ID         string
	Load       map[string]int32
	LastUpdate time.Time
}

func (c *Controller) Workers() []WorkerInfo {
	return slice.Map(
		set.Keys(c.workers.Iter()),
		func(idx int, in string) WorkerInfo {
			load, _ := c.load.Get(in)
			lastUpdate, _ := c.lastUpdate.Get(in)
			return WorkerInfo{
				ID: in,
				Load: set.Map(load.Tasks, func(in *atomic.Int32) int32 {
					return in.Load()
				}),
				LastUpdate: lastUpdate,
			}
		},
	)
}

func (c *Controller) dispatchImpl(ctx context.Context, manifest Manifest, persist ...bool) error {
	// todo store manifest in KV
	log := logger.FromContext(ctx).With(
		zap.String("wk.task.name", manifest.Name),
	)

	candidates := make([]string, 0)
	for id, canRun := range c.dispatch.Iter() {
		if slice.Contains(canRun, manifest.Name) {
			candidates = append(candidates, id)
		}
	}
	log.Debug("potential worker candidates found", zap.Int("count", len(candidates)))

	lowestLoad := int32(math.MaxInt32)
	var atomicLoad *atomic.Int32
	bestCandidate := -1
	for idx, candidate := range candidates {
		loads, ok := c.load.Get(candidate)
		if !ok {
			continue
		}
		load, ok := loads.Tasks[manifest.Name]
		if !ok {
			return stacktrace.NewError("worker '%s' has no associated load entry for '%s'", candidate, manifest.Name)
		}
		if l := load.Load(); l < lowestLoad {
			lowestLoad = l
			bestCandidate = idx
			atomicLoad = load
		}
	}

	if bestCandidate == -1 {
		log.Warn("no candidate matches requirement for immediate task dispatch")
		return nil
	}
	candidateID := candidates[bestCandidate]
	candidateMailbox, ok := c.workers.Get(candidateID)
	log = log.With(zap.String("wk.worker.candidate.id", candidateID))

	if !ok || candidateMailbox == nil {
		log.Warn("best candidate mailbox has gone stale, redispatching")
		return c.dispatchImpl(ctx, manifest, false)
	}

	log.Info("dispatching task to worker", zap.Int32("wk.worker.load", lowestLoad))
	if err := candidateMailbox.Emit(ctx, NewTaskExecuteMessage(manifest.Name, manifest.Payload, manifest.PayloadType)); err != nil {
		return err
	}
	atomicLoad.Inc()

	return nil
}

func (c *Controller) handleWorkerAnnounce(ctx context.Context, msg message.Message, payload *WorkerAnnounceMessage) error {
	log := logger.FromContext(ctx)

	id := payload.WorkerID.String()
	_, ok := c.dispatch.Get(id)
	if !ok {
		log.Info("new worker joined",
			zap.String("wk.worker.id", id),
			zap.Strings("wk.worker.tasks", payload.Tasks),
		)
		c.dispatch.Insert(id, payload.Tasks)
		mailbox, err := c.bus.Subject(ctx, fmt.Sprintf(pollExecuteSubjectFmt, id))
		if err != nil {
			return err
		}
		c.load.Insert(id, WorkerLoad{
			ID: id,
			Tasks: set.From(payload.Tasks, func(item string) (string, *atomic.Int32) {
				return item, atomic.NewInt32(0)
			}),
		})
		c.workers.Insert(id, mailbox)
	}
	c.lastUpdate.Insert(id, msg.CreatedAt())
	c.dispatch.Insert(id, payload.Tasks)
	return nil
}

func (c *Controller) producerImpl(ctx context.Context) error {
	log := logger.FromContext(ctx)

	announces, err := c.bus.Subject(ctx, workerAnnounceSubject)
	if err != nil {
		return stacktrace.Propagate(err, "unable to open announcer topic")
	}
	messages, cancel := announces.ReadC(ctx)
	defer cancel()

	ticker := time.NewTicker(time.Second * 5)

	for {
		select {
		case manifest := <-c.dispatchQueue:
			if err := c.dispatchImpl(ctx, manifest, true); err != nil {
				return stacktrace.Propagate(err, "internal dispatch error")
			}
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			for _, id := range c.lastUpdate.Keys() {
				instant, ok := c.lastUpdate.Get(id)
				if !ok {
					log.Warn("found missing worker update entry while executing book-keeping",
						zap.String("wk.worker.id", id),
					)
				}
				if time.Now().After(instant.Add(time.Second * 5)) {
					log.Warn("worker marked as offline", zap.String("wk.worker.id", id))
					c.dispatch.Delete(id)
					c.workers.Delete(id)
					c.lastUpdate.Delete(id)
				}
			}
		case msg := <-messages:
			switch msg.Kind {
			case MWorkerAnnounce:
				payload := msg.Payload.(*WorkerAnnounceMessage)
				err = stacktrace.Propagate(c.handleWorkerAnnounce(ctx, msg, payload), "could not process worker announce")
			}
			if err != nil {
				return err
			}
		}
	}
}

func (c *Controller) Producer() async.Manifest {
	return async.Of(
		async.Name("wk", "task", "producer"),
		c.producerImpl,
	)
}
