package backoff

import (
	"context"
	"sync"
	"time"
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

func (l *Linear) Duration(ctx context.Context) <-chan time.Duration {
	timer := make(chan time.Duration)
	go func() {
		for {
			l.RLock()
			currentInterval := l.currentInterval
			l.RUnlock()
			select {
			case <-ctx.Done():
				close(timer)
				return
			case timer <- currentInterval:
				l.Lock()
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
