package main

import (
	"context"

	"github.com/palantir/stacktrace"

	"github.com/Raphy42/weekend/core/app"
	"github.com/Raphy42/weekend/modules/platform"
)

func main() {
	modules := std.New(
		platform.Module(),
	)

	sdk, err := app.New("test",
		app.WithModules(modules...),
	)
	if err != nil {
		panic(err)
	}
	if err := sdk.Start(context.Background()); err != nil {
		panic(stacktrace.Propagate(err, "unable to start application"))
	}

	<-sdk.Wait()
}
