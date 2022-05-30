package kernel

import (
	"fmt"
	"time"

	"github.com/Raphy42/weekend/core/kernel/strategy"
)

type Options struct {
	Delay      time.Duration
	Priority   int
	Idempotent bool
	Name       func() string

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

func DefaultOptions() *Options {
	options := Options{
		Delay:                      time.Millisecond * 0,
		Retry:                      DefaultMaximumRetries,
		RetryStrategy:              DefaultRetryStrategy,
		RetryStrategyHandler:       DefaultRetryStrategyHandler, //todo
		Priority:                   DefaultPriority,
		Concurrency:                DefaultConcurrency,
		ConcurrencyStrategy:        strategy.ConstantConcurrencyStrategy,
		ConcurrencyStrategyHandler: strategy.ConstantConcurrencyStrategyHandler(DefaultConcurrency),
		Idempotent:                 false,
	}
	localValue := &options
	namingFunction := func() string {
		return fmt.Sprintf("default.%d", localValue.Priority)
	}
	localValue.Name = namingFunction
	return localValue
}

func (o *Options) Apply(opts ...Option) {
	for _, opt := range opts {
		opt(o)
	}
}

type Option func(options *Options)
