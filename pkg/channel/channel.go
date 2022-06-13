package channel

import (
	"context"
	"sync"

	"github.com/palantir/stacktrace"

	"github.com/Raphy42/weekend/core/logger"
)

//Send reduces boilerplate for blocking unbuffered channels and ensures propagation of context termination
func Send[T any](ctx context.Context, value T, c chan<- T) error {
	select {
	case <-ctx.Done():
		return stacktrace.Propagate(ctx.Err(), "context canceled before value could be dispatched")
	case c <- value:
		return nil
	}
}

//Multicast dispatch the input channel message to the multiple others
func Multicast[T any](parent context.Context, in <-chan T, outs ...chan<- T) context.CancelFunc {
	ctx, cancel := context.WithCancel(parent)
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case msg := <-in:
				var wg sync.WaitGroup
				for _, out := range outs {
					wg.Add(1)
					go func(c chan<- T) {
						defer wg.Done()
						if err := Send(ctx, msg, c); err != nil {
							return
						}
					}(out)
				}
				wg.Wait()
			}
		}
	}()

	return func() {
		log := logger.FromContext(parent)

		log.Warn("reader context canceled")
		cancel()
	}
}
