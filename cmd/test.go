package main

import (
	"context"

	"github.com/palantir/stacktrace"

	"github.com/Raphy42/weekend/core/app"
	"github.com/Raphy42/weekend/modules/platform"
	"github.com/Raphy42/weekend/pkg/slice"
)

func main() {
	modules := slice.New(
		platform.Module(),
	)

	badCtx, cancel := context.WithCancel(context.Background())
	cancel()

	sdk, err := app.New("test",
		app.WithModules(modules...),
	)
	if err != nil {
		panic(err)
	}
	if err := sdk.Start(badCtx); err != nil {
		panic(stacktrace.Propagate(err, "unable to start application"))
	}

	<-sdk.Wait()
}
