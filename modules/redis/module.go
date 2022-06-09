package redis

import "github.com/Raphy42/weekend/core/di"

const (
	ModuleName = "wk.redis"
)

func Module() di.Module {
	return di.Declare(
		ModuleName,
	)
}
