package core

import (
	"github.com/Raphy42/weekend/core/dep"
)

var (
	ModuleName = dep.Name("wk", "platform")
)

func Module(opts ...ModuleOption) dep.Module {
	options := defaultModuleOptions()
	options.apply(opts...)

	return dep.Declare(
		ModuleName,
		dep.Factories(
			healthProvider,
			engineBuilderProvider,
			applicationContextProvider(options.rootCtx),
			configFromFilenamesProvider(options.configFilenames...),
		),
		dep.SideEffects(
			platformInformation,
			applicationEngineInjector,
		),
	)
}
