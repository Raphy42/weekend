package scheduler

import (
	"time"

	"github.com/Raphy42/weekend/core/scheduler/strategy"
)

type Options struct {
	Delay      time.Duration
	Priority   int
	Idempotent bool

	// strategies
	Retry                      int
	RetryStrategy              strategy.RetryStrategy
	RetryStrategyHandler       strategy.RetryStrategyHandler
	Concurrency                int
	ConcurrencyStrategy        strategy.ConcurrencyStrategy
	ConcurrencyStrategyHandler strategy.ConcurrencyStrategyHandler
	FailureStrategy            strategy.FailureStrategy
	FailureStrategyHandler     strategy.FailureStrategyHandler
}

func newOptions() *Options {
	return &Options{
		Delay:                      time.Millisecond * 0,
		Retry:                      DefaultMaximumRetries,
		RetryStrategy:              DefaultRetryStrategy,
		RetryStrategyHandler:       nil, //todo
		Priority:                   DefaultPriority,
		Concurrency:                DefaultConcurrency,
		ConcurrencyStrategy:        strategy.ConstantConcurrencyStrategy,
		ConcurrencyStrategyHandler: strategy.ConstantConcurrencyStrategyHandler(DefaultConcurrency),
		Idempotent:                 false,
	}
}

func (o *Options) Apply(opts ...Option) {
	for _, opt := range opts {
		opt(o)
	}
}

type Option func(options *Options)

func WithDelay(delay time.Duration) Option {
	return func(options *Options) {
		options.Delay = delay
	}
}

func WithRetryCount(retry int) Option {
	return func(options *Options) {
		options.Retry = retry
		options.RetryStrategy = strategy.ExponentialBackoffRetryStrategy //todo
		options.RetryStrategyHandler = nil                               // todo
	}
}

func WithRetryStrategy(strat strategy.RetryStrategy) Option {
	return func(options *Options) {
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

func WithCustomRetryStrategy(handler strategy.RetryStrategyHandler) Option {
	return func(options *Options) {
		options.RetryStrategy = strategy.CustomRetryStrategy
		options.RetryStrategyHandler = handler
	}
}

func WithCustomFailureStrategy(handler strategy.FailureStrategyHandler) Option {
	return func(options *Options) {
		options.FailureStrategy = strategy.CustomFailureStrategy
		options.FailureStrategyHandler = handler
	}
}

func WithConcurrency(count int) Option {
	return func(options *Options) {
		options.Concurrency = count
		options.ConcurrencyStrategy = strategy.ConstantConcurrencyStrategy
		options.ConcurrencyStrategyHandler = strategy.ConstantConcurrencyStrategyHandler(count)
	}
}

func WithBinpack() Option {
	return func(options *Options) {
		options.Concurrency = -1
		options.ConcurrencyStrategy = strategy.BinpackConcurrencyStrategy
		options.ConcurrencyStrategyHandler = strategy.BinpackConcurrencyStrategyHandler()
	}
}
