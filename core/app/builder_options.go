package app

import (
	"github.com/Raphy42/weekend/core/di"
	"github.com/Raphy42/weekend/core/message"
)

func WithBus(bus message.Bus) BuilderOption {
	return func(builder *Builder) error {
		builder.bus = bus
		return nil
	}
}

func WithModules(modules ...di.Module) BuilderOption {
	return func(builder *Builder) error {
		builder.modules = append(builder.modules, modules...)
		return nil
	}
}
