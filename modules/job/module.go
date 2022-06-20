package job

import "github.com/Raphy42/weekend/core/dep"

var (
	WorkerModuleName     = dep.Name("wk", "job", "worker")
	ControllerModuleName = dep.Name("wk", "job", "controller")
)

func WorkerModule(opts ...Option) dep.Module {
	options := defaultOptions()
	options.apply(opts...)

	return dep.Declare(
		WorkerModuleName,
		dep.Factories(
			newWorkerFactory(options),
		),
	)
}

func ControllerModule(opts ...Option) dep.Module {
	options := defaultOptions()
	options.apply(opts...)

	return dep.Declare(
		ControllerModuleName,
		dep.Factories(
			controllerFactory,
			apiEndpointFactory,
		),
	)
}
