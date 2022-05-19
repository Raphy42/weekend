package strategy

type (
	RetryStrategy        int
	RetryStrategyHandler func(current, max int) bool
)

const (
	LinearBackoffRetryStrategy RetryStrategy = iota
	ExponentialBackoffRetryStrategy
	ManualRetryStrategy
	CustomRetryStrategy
)
