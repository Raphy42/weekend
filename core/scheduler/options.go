package scheduler

import (
	"time"

	"github.com/Raphy42/weekend/core/kernel"
	"github.com/Raphy42/weekend/core/kernel/strategy"
)

func WithDelay(delay time.Duration) kernel.Option {
	return func(options *kernel.Options) {
		options.Delay = delay
	}
}

func WithRetryCount(retry int) kernel.Option {
	return func(options *kernel.Options) {
		options.Retry = retry
		options.RetryStrategy = strategy.ExponentialBackoffRetryStrategy //todo
		options.RetryStrategyHandler = nil                               // todo
	}
}

func WithRetryStrategy(strat strategy.RetryStrategy) kernel.Option {
	return func(options *kernel.Options) {
		options.RetryStrategy = strat
		switch strat {
		case strategy.ManualRetryStrategy:
			panic("not implemented")
		case strategy.ExponentialBackoffRetryStrategy:
			panic("not implemented")
		case strategy.LinearBackoffRetryStrategy:
			panic("not implemented")
		default:
			panic("not implemented")
		}
	}
}

func WithCustomRetryStrategy(handler strategy.RetryStrategyHandler) kernel.Option {
	return func(options *kernel.Options) {
		options.RetryStrategy = strategy.CustomRetryStrategy
		options.RetryStrategyHandler = handler
	}
}

func WithFailureStrategy(failureStrategy strategy.FailureStrategy) kernel.Option {
	return func(options *kernel.Options) {
		options.FailureStrategy = failureStrategy
	}
}

func WithCustomFailureStrategy(handler strategy.FailureStrategyHandler) kernel.Option {
	return func(options *kernel.Options) {
		options.FailureStrategy = strategy.CustomFailureStrategy
		options.FailureStrategyHandler = handler
	}
}

func WithConcurrency(count int) kernel.Option {
	return func(options *kernel.Options) {
		options.Concurrency = count
		options.ConcurrencyStrategy = strategy.ConstantConcurrencyStrategy
		options.ConcurrencyStrategyHandler = strategy.ConstantConcurrencyStrategyHandler(count)
	}
}

func WithBinpack() kernel.Option {
	return func(options *kernel.Options) {
		options.Concurrency = -1
		options.ConcurrencyStrategy = strategy.BinpackConcurrencyStrategy
		options.ConcurrencyStrategyHandler = strategy.BinpackConcurrencyStrategyHandler()
	}
}
