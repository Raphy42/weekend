package nats

import "github.com/Raphy42/weekend/core/dep"

var (
	ModuleName = dep.Name("wk", "nats")
)

func Module() dep.Module {
	return dep.Declare(
		ModuleName,
		dep.Factories(
			natsFactory,
		),
	)
}
