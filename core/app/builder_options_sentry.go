//go:build ops.sentry

package app

import (
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/palantir/stacktrace"

	"github.com/Raphy42/weekend/core"
)

func WithSentry(sentryDSN string) BuilderOption {
	return func(builder *Builder) error {
		core.RegisterOnStartHook(func() error {
			return stacktrace.Propagate(sentry.Init(sentry.ClientOptions{
				Dsn:              sentryDSN,
				AttachStacktrace: true,
				//todo adjust sample rate
				TracesSampleRate: 1.0,
			}), "unable to start sentry")
		})
		core.RegisterOnStopHook(func() error {
			if ok := sentry.Flush(time.Second * 2); !ok {
				return stacktrace.NewError("unable to flush sentry, some error might have been lost")
			}
			return nil
		})
		return nil
	}
}
