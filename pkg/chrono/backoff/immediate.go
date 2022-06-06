package backoff

import (
	"context"
	"time"
)

type Immediate struct{}

func NewImmediate() *Immediate {
	return &Immediate{}
}

func (i Immediate) Duration(ctx context.Context) <-chan time.Duration {
	timer := make(chan time.Duration)
	go func() {
		for {
			select {
			case <-ctx.Done():
				close(timer)
				return
			case timer <- 0:
				// select is here used as a non-blocking dispatch
			}
		}
	}()
	return timer
}

func (i Immediate) Reset() {
	// noop
}
