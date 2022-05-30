package kernel

import (
	"github.com/Raphy42/weekend/core/kernel/strategy"
)

const (
	DefaultMaximumRetries  = 5
	DefaultRetryStrategy   = strategy.ExponentialBackoffRetryStrategy
	DefaultFailureStrategy = strategy.AnyErrorFailureStrategy
	DefaultPriority        = 10
	DefaultConcurrency     = 16
)

var (
	DefaultRetryStrategyHandler = strategy.ExponentialBackoffRetryStrategyHandler()
)
