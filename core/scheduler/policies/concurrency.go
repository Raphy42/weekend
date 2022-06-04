package policies

type ConcurrencyPolicy int

const (
	ConcurrencyAuto ConcurrencyPolicy = iota
	ConcurrencyConstant
)
