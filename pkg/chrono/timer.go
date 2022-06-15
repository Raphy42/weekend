package chrono

import (
	"context"
	"time"
)

type Ticker struct {
	interval time.Duration
}

func NewTicker(interval time.Duration) *Ticker {
	return &Ticker{interval: interval}
}

func (t *Ticker) Tick(ctx context.Context, fn func()) {
	ticker := time.NewTicker(t.interval)
	defer ticker.Stop()

loop:
	for {
		select {
		case <-ticker.C:
			fn()
		case <-ctx.Done():
			break loop
		}
	}
}
func (t *Ticker) TickErr(ctx context.Context, fn func() error) <-chan error {
	errs := make(chan error)
	t.Tick(ctx, func() {
		if err := fn(); err != nil {
			errs <- err
		}
	})
	return errs
}
