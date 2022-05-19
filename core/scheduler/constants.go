package scheduler

import "weekend/core/scheduler/strategy"

const (
	DefaultMaximumRetries  = 5
	DefaultRetryStrategy   = strategy.ExponentialBackoffRetryStrategy
	DefaultFailureStrategy = strategy.AnyErrorFailureStrategy
	DefaultPriority        = 10
	DefaultConcurrency     = 16
)
