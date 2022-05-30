package backoff

import (
	"context"
	"time"

	"go.uber.org/atomic"
)

type Linear struct {
	initialInterval time.Duration
	counter         *atomic.Int64
	currentInterval *atomic.Duration
}

func NewLinear(interval time.Duration) *Linear {
	return &Linear{
		initialInterval: interval,
		counter:         atomic.NewInt64(0),
		currentInterval: atomic.NewDuration(interval),
	}
}

func (l *Linear) Interval(ctx context.Context) <-chan time.Time {
	timer := make(chan time.Time, 1)
	go func() {
		for {
			select {
			case <-ctx.Done():
				close(timer)
				return
			case t := <-time.After(l.currentInterval.Load()):
				//todo remove races conditions and maybe use RWLock instead of double atomic (needs semaphore for perfect sync)
				timer <- t
				value := l.counter.Inc()
				newDuration := time.Duration(value) * l.initialInterval
				l.currentInterval.Store(newDuration)
			}
		}
	}()
	return timer
}

func (l *Linear) Reset() {
	l.counter.Store(0)
}
