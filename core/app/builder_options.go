package app

import (
	"github.com/Raphy42/weekend/core/dep"
	"github.com/Raphy42/weekend/core/message"
)

func WithBus(bus message.Mailbox) BuilderOption {
	return func(builder *Builder) error {
		builder.bus = bus
		return nil
	}
}

func WithModules(modules ...dep.Module) BuilderOption {
	return func(builder *Builder) error {
		builder.modules = append(builder.modules, modules...)
		return nil
	}
}
