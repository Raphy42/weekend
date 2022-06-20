package supervisor

import "github.com/Raphy42/weekend/core/scheduler/async"

type Strategy struct {
	Restart     RestartStrategy
	Shutdown    ShutdownStrategy
	Supervision SupervisionStrategy
}

func (s *Strategy) apply(opts ...StrategyOption) {
	for _, opt := range opts {
		opt(s)
	}
}

type StrategyOption func(strategy *Strategy)

func WithRestartStrategy(restart RestartStrategy) StrategyOption {
	return func(strategy *Strategy) {
		strategy.Restart = restart
	}
}

func WithShutdownStrategy(shutdown ShutdownStrategy) StrategyOption {
	return func(strategy *Strategy) {
		strategy.Shutdown = shutdown
	}
}

func WithSupervisionStrategy(supervision SupervisionStrategy) StrategyOption {
	return func(strategy *Strategy) {
		strategy.Supervision = supervision
	}
}

func NewStrategy() Strategy {
	return Strategy{
		Restart:     TransientRestartStrategy,
		Shutdown:    ImmediateShutdownStrategy,
		Supervision: OneForOneSupervisionStrategy,
	}
}

type Spec struct {
	Strategy Strategy
	Manifest async.Manifest
	Args     any
}

func NewSpec(manifest async.Manifest, args any, opts ...StrategyOption) Spec {
	strategy := NewStrategy()
	strategy.apply(opts...)

	return Spec{
		Strategy: strategy,
		Manifest: manifest,
		Args:     args,
	}
}
