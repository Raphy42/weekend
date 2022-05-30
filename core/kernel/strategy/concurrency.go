package strategy

type (
	ConcurrencyStrategy        int
	ConcurrencyStrategyHandler func() int
)

const (
	ConstantConcurrencyStrategy ConcurrencyStrategy = iota
	BinpackConcurrencyStrategy
)

func ConstantConcurrencyStrategyHandler(value int) ConcurrencyStrategyHandler {
	return func() int {
		return value
	}
}

func BinpackConcurrencyStrategyHandler() ConcurrencyStrategyHandler {
	return func() int {
		panic("unimplemented")
	}
}
