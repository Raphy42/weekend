package grpc

import (
	"fmt"

	"github.com/Raphy42/weekend/core/di"
)

func Module(namespace string) di.Module {
	return di.Declare(
		fmt.Sprintf("wk.platform.%s", namespace),
	)
}
