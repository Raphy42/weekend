package backoff

import (
	"context"
	"sync"
	"time"

	"github.com/Raphy42/weekend/pkg/channel"
)

type Linear struct {
	sync.RWMutex
	initialInterval time.Duration
	counter         int64
	currentInterval time.Duration
}

func NewLinear(interval time.Duration) *Linear {
	return &Linear{
		initialInterval: interval,
		counter:         0,
		currentInterval: interval,
	}
}

func (l *Linear) Interval(ctx context.Context) <-chan time.Time {
	timer := make(chan time.Time, 1)
	go func() {
		for {
			l.RLock()
			currentInterval := l.currentInterval
			l.RUnlock()
			select {
			case <-ctx.Done():
				close(timer)
				return
			case t := <-time.After(currentInterval):
				l.Lock()
				_ = channel.Send(ctx, t, timer)
				l.counter += 1
				newDuration := time.Duration(l.counter) * l.initialInterval
				l.currentInterval = newDuration
				l.Unlock()
			}
		}
	}()
	return timer
}

func (l *Linear) Reset() {
	l.Lock()
	defer l.Unlock()
	l.counter = 0
}
