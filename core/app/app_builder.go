package app

import (
	"sync"

	"github.com/palantir/stacktrace"

	"github.com/Raphy42/weekend/core/dependency"
	"github.com/Raphy42/weekend/core/di"
	"github.com/Raphy42/weekend/core/message"
	"github.com/Raphy42/weekend/core/scheduler"
)

type Builder struct {
	name    string
	bus     message.Bus
	modules []di.Module
}

func (b *Builder) Apply(opts ...BuilderOption) error {
	for _, opt := range opts {
		if err := opt(b); err != nil {
			return stacktrace.Propagate(err, "builder option returned non-nil error")
		}
	}
	return nil
}

func (b *Builder) Build() (*App, error) {
	if b.bus == nil {
		b.bus = message.NewInMemoryBus()
	}

	container := dependency.NewContainer(b.name)
	return &App{
		lock:      sync.RWMutex{},
		name:      b.name,
		container: container.Use(b.modules...),
		scheduler: scheduler.New(b.bus),
	}, nil
}

type BuilderOption func(builder *Builder) error

func newBuilder(name string) *Builder {
	return &Builder{
		name:    name,
		bus:     nil,
		modules: make([]di.Module, 0),
	}
}
