package scheduler

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Handle struct {
	ID      uuid.UUID
	Context *SchedulingContext

	// stack-related
	Scheduler string

	// manifest
	Idempotent bool

	// policies
	Retries        int
	MaxRetry       int
	Concurrency    int
	MaxConcurrency int

	// time
	CreatedAt   time.Time
	ScheduledAt *time.Time
	StartedAt   *time.Time
	FinishedAt  *time.Time
}

type Scheduler interface {
	Schedule(ctx context.Context, fn func(ctx context.Context) error, opts ...Option) error
}
