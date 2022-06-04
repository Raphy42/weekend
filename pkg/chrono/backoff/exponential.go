package backoff

import (
	"context"
	"time"

	"go.uber.org/atomic"

	"github.com/Raphy42/weekend/pkg/channel"
	"github.com/Raphy42/weekend/pkg/slice"
)

var (
	fibonacciSequence = []int64{0, 1, 1, 2, 3, 5, 8, 13, 21, 34, 55, 89, 144, 233, 377, 610, 987, 1597, 2584, 4181, 6765, 10946}
)

//todo really implement IPBA and not some fibonacci based backoff

//Exponential implements the exponential backoff as defined by [this wikipedia article](https://en.wikipedia.org/wiki/Exponential_backoff).
// More specifically a homebrew version based on the IPBA algorithm from [Intelligent Paging Backoff Algorithm for IEEE 802.11 MAC Protocol](https://www.researchgate.net/figure/IPBA-algorithm-description_fig6_258402123)
type Exponential struct {
	cw              time.Duration
	cwf             time.Duration
	x               time.Duration
	initialInterval time.Duration
	slots           []time.Duration
	shouldReset     *atomic.Bool
}

func NewExponentialBackoff(interval time.Duration) *Exponential {
	sequence := slice.Map(fibonacciSequence, func(idx int, in int64, eOut *[]time.Duration) {
		(*eOut)[idx] = time.Duration(in) * interval
	})
	return &Exponential{
		cw:              interval,
		cwf:             interval / 2,
		x:               0,
		initialInterval: interval,
		slots:           sequence,
		shouldReset:     atomic.NewBool(false),
	}
}

func (e *Exponential) nextInterval() time.Duration {
	for _, interval := range e.slots {
		if interval >= e.x && interval <= e.cw {
			if !e.shouldReset.Load() {
				e.cwf = e.cw
				e.cw *= 2
				e.x = e.cwf
			} else {
				e.cw = e.initialInterval
				e.cwf = e.initialInterval / 2
				e.x = 0
				e.shouldReset.Store(false)
			}
			return interval
		}
	}
	return e.initialInterval
}

func (e *Exponential) Time(ctx context.Context) <-chan time.Time {
	timer := make(chan time.Time, 1)
	go func() {
		for {
			select {
			case <-ctx.Done():
				close(timer)
				return
			case t := <-time.After(e.nextInterval()):
				_ = channel.Send(ctx, t, timer)
			}
		}
	}()
	return timer
}

func (e *Exponential) Reset() {
	e.shouldReset.Store(true)
}
