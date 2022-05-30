package backoff

import (
	"context"
	"time"
)

type Exponential struct {
	initialInterval time.Duration
}

func (e Exponential) Time(ctx context.Context) <-chan time.Time {
	//TODO implement me
	panic("implement me")
}

func (e Exponential) Reset() {
	//TODO implement me
	panic("implement me")
}
