//go:build ops.sentry

package errors

import (
	"time"

	"github.com/getsentry/sentry-go"
)

func InstallPanicHandler() {
	if err := recover(); err != nil {
		sentry.CurrentHub().Recover(err)
		sentry.Flush(time.Second * 5)
	}
}
