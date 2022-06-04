package main

import (
	"context"

	"github.com/Raphy42/weekend/core/app"
	"github.com/Raphy42/weekend/modules/platform"
	"github.com/Raphy42/weekend/pkg/slice"
)

func main() {
	modules := slice.New(
		platform.Module(),
	)

	sdk, err := app.New("test",
		app.WithModules(modules...),
	)
	if err != nil {
		panic(err)
	}
	if err := sdk.Start(context.Background()); err != nil {
		panic(err)
	}

	<-sdk.Wait()
}
