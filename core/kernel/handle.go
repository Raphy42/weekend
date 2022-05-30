package kernel

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"

	"github.com/Raphy42/weekend/core/errors"
)

type HandleTimestamp struct {
	CreatedAt   time.Time
	ScheduledAt *time.Time
	StartedAt   *time.Time
	FinishedAt  *time.Time
}
type HandleTimestamps []HandleTimestamp

func (h *HandleTimestamps) newPeriod() {
	*h = append(*h, HandleTimestamp{
		CreatedAt: time.Now(),
	})
}

func (h *HandleTimestamps) last() *HandleTimestamp {
	v := (*h)[len(*h)-1]
	return &v
}

type HandleMetadata struct {
	// args
	Args Args

	// stack-related
	Scheduler string

	// manifest
	Idempotent bool

	// policies
	Tries       int
	MaxRetry    int
	Concurrency int

	// time
	Timestamps HandleTimestamps
}

type Handle struct {
	ID              uuid.UUID
	Name            string
	Context         *SchedulingContext
	Metadata        *HandleMetadata
	Options         Options
	ResultChan      <-chan interface{}
	ErrorChan       <-chan error
	CompletionHooks []func()
}

func NewHandle(options Options, args Args, ctx ...context.Context) (*Handle, chan interface{}, chan error) {
	errorChan := make(chan error)
	resultChan := make(chan interface{})
	handle := Handle{
		ID:      uuid.New(),
		Name:    options.Name(),
		Context: NewSchedulingContext(ctx...),
		Metadata: &HandleMetadata{
			Args:        args,
			Scheduler:   "<nobody>",
			Idempotent:  options.Idempotent,
			Concurrency: options.Concurrency,
			Tries:       0,
			MaxRetry:    options.Retry,
			Timestamps:  make(HandleTimestamps, 0),
		},
		ResultChan:      resultChan,
		ErrorChan:       errorChan,
		CompletionHooks: make([]func(), 0),
	}
	handle.Metadata.Timestamps.newPeriod()
	return &handle, resultChan, errorChan
}

func (h *Handle) validate() error {
	if err := IsTooManyRetries(h.Metadata.Tries, h.Metadata.MaxRetry); err != nil {
		return err
	}
	return nil
}

func (h *Handle) MarkAsScheduled(whom string) {
	now := time.Now()
	h.Metadata.Tries += 1
	h.Metadata.Scheduler = whom
	h.Metadata.Timestamps.last().ScheduledAt = &now
}

func (h *Handle) markAsStarted() {
	now := time.Now()
	h.Metadata.Timestamps.last().StartedAt = &now
}

func (h *Handle) markAsFinished() {
	now := time.Now()
	h.Metadata.Timestamps.last().FinishedAt = &now
	for _, hook := range h.CompletionHooks {
		hook()
	}
}

func (h *Handle) Poll(ctx context.Context) (interface{}, error) {
	for {
		select {
		case <-ctx.Done():
			h.markAsFinished()
			return nil, ctx.Err()
		case result := <-h.ResultChan:
			h.markAsFinished()
			return result, nil
		case err := <-h.ErrorChan:
			h.markAsFinished()
			return nil, err
		}
	}
}

type HandleGroup []*Handle

func (h HandleGroup) Poll(ctx context.Context) ([]interface{}, errors.Group) {
	var wg sync.WaitGroup

	results := make([]interface{}, 0)
	errs := make(errors.Group, 0)

	for idx, handle := range h {
		wg.Add(1)
		go func(routineCtx context.Context, h *Handle, index int) {
			results[index], errs[index] = h.Poll(routineCtx)
			wg.Done()
		}(ctx, handle, idx)
	}

	wg.Wait()
	return results, errs
}
