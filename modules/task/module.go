package task

import "github.com/Raphy42/weekend/core/dep"

var (
	ModuleName = dep.Name("wk", "task")
)

func Module() dep.Module {
	return dep.Declare(
		ModuleName,
		dep.Factories(),
	)
}
