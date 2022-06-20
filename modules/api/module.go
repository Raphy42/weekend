package api

import "github.com/Raphy42/weekend/core/dep"

var (
	ModuleName = dep.Name("wk", "api")
)

func Module() dep.Module {
	return dep.Declare(
		ModuleName,
		dep.Factories(
			ginEngineFactory,
		),
		dep.SideEffects(
			startServer,
		),
	)
}
