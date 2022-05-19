package strategy

import "weekend/core/errors"

type (
	FailureStrategy        int
	FailureStrategyHandler func(err error) bool
)

const (
	AnyErrorFailureStrategy FailureStrategy = iota
	PersistentErrorFailureStrategy
	CustomFailureStrategy
)

func AnyErrorFailureStrategyHandler() FailureStrategyHandler {
	return func(err error) bool {
		return err != nil
	}
}

func PersistentErrorFailureStrategyHandler() FailureStrategyHandler {
	return func(err error) bool {
		return errors.HasFlag(err, errors.KPersistent)
	}
}

func CustomFailureStrategyHandler(fn func(err error) bool) FailureStrategyHandler {
	return fn
}
