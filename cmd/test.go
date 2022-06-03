package main

import (
	"context"
	"time"

	"github.com/Raphy42/weekend/core/errors"
	"github.com/Raphy42/weekend/core/kernel/scheduler"
)

func main() {
	complexFn := scheduler.Make("capitalize", func(ctx context.Context, args interface{}) (interface{}, error) {
		time.Sleep(time.Second * 2)
		return "foobar", nil
	})
	fn := scheduler.Make("pipeline", func(ctx context.Context, args interface{}) error {
		handle, err := scheduler.Schedule(ctx, complexFn, args)
		if err != nil {
			return err
		}
		_, err = handle.Poll(ctx)
		return err
	})

	sched := scheduler.New()
	handle, err := sched.Schedule(context.Background(), fn, nil)
	errors.Must(err)

	_, err = handle.Poll(context.Background())
	errors.Must(err)

}
