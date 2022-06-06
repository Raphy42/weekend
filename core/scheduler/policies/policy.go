package policies

type PolicyAction struct {
	Restart bool
}

type Policy struct {
	Idempotent        bool
	RetryPolicy       RetryPolicy
	Retry             int
	ConcurrencyPolicy ConcurrencyPolicy
	Concurrency       int
	BackoffPolicy     BackoffPolicy
}

func Default() Policy {
	return Policy{
		Idempotent:        false,
		RetryPolicy:       RetryNever,
		Retry:             0,
		ConcurrencyPolicy: ConcurrencyConstant,
		Concurrency:       16,
		BackoffPolicy:     BackoffExponential,
	}
}
