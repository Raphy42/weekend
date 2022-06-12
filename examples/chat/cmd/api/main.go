package main

import (
	"context"
	"os"

	"github.com/Raphy42/weekend/core/app"
	"github.com/Raphy42/weekend/core/errors"
	"github.com/Raphy42/weekend/core/telemetry"
	"github.com/Raphy42/weekend/modules/core"
	"github.com/Raphy42/weekend/modules/redis"
)

// build with the following tags
// - ops.sentry
// - task.encoding.msgpack

const name = "chat.api"

func main() {

	defer errors.InstallPanicObserver()

	ctx, span := telemetry.Install(name, context.Background())
	defer span.End()

	sdk, err := app.New(name,
		app.WithSentry(os.Getenv("SENTRY_DSN")),
		app.WithModules(
			core.Module(
				core.WithContext(ctx),
				core.WithConfigFilenames("./examples/chat/common.yml"),
			),
			redis.Module(),
		),
	)
	errors.Mustf(err, "could not create application")

	errors.Mustf(sdk.Start(ctx), "could not start application")
	errors.Mustf(<-sdk.Wait(ctx), "application shut down with non-nil error")
}
