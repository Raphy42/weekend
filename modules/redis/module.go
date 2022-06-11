package redis

import "github.com/Raphy42/weekend/core/di"

var (
	ModuleName = di.Name("wk", "redis")
)

func Module() di.Module {
	return di.Declare(
		ModuleName,
	)
}
