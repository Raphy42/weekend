package scheduler

import "github.com/Raphy42/weekend/core/scheduler/strategy"

const (
	DefaultMaximumRetries  = 5
	DefaultRetryStrategy   = strategy.ExponentialBackoffRetryStrategy
	DefaultFailureStrategy = strategy.AnyErrorFailureStrategy
	DefaultPriority        = 10
	DefaultConcurrency     = 16
)
