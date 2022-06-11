package grpc

import (
	"github.com/Raphy42/weekend/core/di"
)

func Module(namespace string) di.Module {
	return di.Declare(
		di.Name("wk", "grpc", namespace),
	)
}
