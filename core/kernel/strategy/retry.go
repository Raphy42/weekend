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

func LinearBackoffRetryStrategyHandler() RetryStrategyHandler {
	return func(current, max int) bool {
		return false
	}
}

func ExponentialBackoffRetryStrategyHandler() RetryStrategyHandler {
	return func(current, max int) bool {
		return false
	}
}

func ManualRetryStrategyHandler() RetryStrategyHandler {
	return func(current, max int) bool {
		return false
	}
}

func CustomRetryStrategyHandler() RetryStrategyHandler {
	return func(current, max int) bool {
		return false
	}
}
