//go:build !ops.sentry

package errors

import "github.com/Raphy42/weekend/core/logger"

func InstallPanicObserver() {
	logger.Flush()
	//noop
}
