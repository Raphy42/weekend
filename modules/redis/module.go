package redis

import (
	"github.com/Raphy42/weekend/core/dep"
)

var (
	ModuleName = dep.Name("wk", "redis")
)

func Module() dep.Module {
	return dep.Declare(
		ModuleName,
		dep.Factories(
			clientFactory,
		),
		dep.SideEffects(
			redisVersion,
			redisHealthCheck,
		),
	)
}
