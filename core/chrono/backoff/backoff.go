package backoff

import (
	"context"
	"time"
)

type Backoff interface {
	Time(ctx context.Context) <-chan time.Time
	Reset()
}
