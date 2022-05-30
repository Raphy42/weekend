package main

import (
	"context"
	"math"
	"time"

	"github.com/Raphy42/weekend/core/kernel"
	"github.com/Raphy42/weekend/core/scheduler"
)

func main() {
	root := context.Background()
	ctx, cancel := context.WithTimeout(root, time.Second*15)
	defer cancel()

	system := scheduler.NewSystem(
		kernel.Task("test", func(ctx context.Context, args interface{}) (interface{}, error) {
			return true, nil
		}),
		kernel.Task("foobar", func(ctx context.Context, args interface{}) (interface{}, error) {
			t := time.NewTicker(5 * time.Second)
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-t.C:
				return nil, nil
			}
		}),
	)

	for i := 0; i < math.MaxInt8; i++ {
		_, _ = system.Next(ctx, kernel.NewEvent(scheduler.EventRunTask, &scheduler.RunTaskEventPayload{
			Name:    "test",
			Payload: nil,
		}))
		_, _ = system.Next(ctx, kernel.NewEvent(scheduler.EventRunTask, &scheduler.RunTaskEventPayload{
			Name:    "foobar",
			Payload: nil,
		}))
	}
	<-ctx.Done()
}
