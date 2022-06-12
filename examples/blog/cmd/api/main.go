package main

import (
	"context"
	"os"

	"github.com/Raphy42/weekend/core/app"
	"github.com/Raphy42/weekend/core/errors"
	"github.com/Raphy42/weekend/modules/core"
	"github.com/Raphy42/weekend/modules/telemetry"
)

// build with the following tags
// - ops.sentry
// - task.encoding.msgpack

func main() {
	defer errors.InstallPanicObserver()

	ctx := context.Background()

	sdk, err := app.New("api",
		app.WithSentry(os.Getenv("SENTRY_DSN")),
		app.WithModules(
			core.Module(
				core.WithContext(ctx),
				core.WithConfigFilenames("./examples/blog/common.yml"),
			),
			telemetry.Module(),
		),
	)
	errors.Mustf(err, "could not create application")

	errors.Mustf(sdk.Start(ctx), "could not start application")
	<-sdk.Wait()
}
