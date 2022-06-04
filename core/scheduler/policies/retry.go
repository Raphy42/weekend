package policies

type RetryPolicy int

const (
	RetryNever RetryPolicy = iota
	RetryAlways
	RetryConstant
)
