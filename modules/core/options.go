package core

import "context"

type ModuleOptions struct {
	configFilenames []string
	rootCtx         context.Context
}
type ModuleOption func(options *ModuleOptions)

func defaultModuleOptions() ModuleOptions {
	return ModuleOptions{
		configFilenames: make([]string, 0),
		rootCtx:         context.Background(),
	}
}

func (m *ModuleOptions) apply(opts ...ModuleOption) {
	for _, opt := range opts {
		opt(m)
	}
}

func WithContext(ctx context.Context) ModuleOption {
	return func(options *ModuleOptions) {
		options.rootCtx = ctx
	}
}

func WithConfigFilenames(filenames ...string) ModuleOption {
	return func(options *ModuleOptions) {
		options.configFilenames = append(options.configFilenames, filenames...)
	}
}
