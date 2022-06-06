package backoff

import (
	"context"
	"time"
)

type Backoff interface {
	Duration(ctx context.Context) <-chan time.Duration
	Reset()
}
