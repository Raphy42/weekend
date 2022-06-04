package platform

import (
	"context"
	"runtime"

	"go.uber.org/zap"

	"github.com/Raphy42/weekend/core/di"
	"github.com/Raphy42/weekend/core/logger"
)

func platformInformation(ctx context.Context) {
	log := logger.FromContext(ctx)

	log.Info("platform information",
		zap.String("os.architecture", runtime.GOARCH),
		zap.String("os.kernel", runtime.GOOS),
		zap.String("go.version", runtime.Version()),
		zap.Int("go.concurrency", runtime.GOMAXPROCS(runtime.NumCPU())),
	)
}

func Module() di.Module {
	return di.Declare(
		"wk.platform",
		di.Invoke(platformInformation),
	)
}
