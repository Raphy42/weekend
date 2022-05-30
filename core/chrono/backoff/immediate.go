package backoff

import (
	"context"
	"time"
)

type Immediate struct{}

func NewImmediate() *Immediate {
	return &Immediate{}
}

func (i Immediate) Time(ctx context.Context) <-chan time.Time {
	timer := make(chan time.Time, 1)
	go func() {
		for {
			select {
			case <-ctx.Done():
				close(timer)
				return
			case timer <- time.Now():
				// select is here used as a non-blocking dispatch
			}
		}
	}()
	return timer
}

func (i Immediate) Reset() {
	// noop
}
