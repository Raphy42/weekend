package kernel

import (
	"context"
)

type Scheduler interface {
	Schedule(ctx context.Context, schedulable Schedulable, args interface{}) (*Handle, error)
}
