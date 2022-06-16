package app

import (
	"sync"

	"github.com/palantir/stacktrace"

	"github.com/Raphy42/weekend/core/dep"
	"github.com/Raphy42/weekend/core/message"
	"github.com/Raphy42/weekend/core/scheduler"
)

type Builder struct {
	sync.RWMutex
	done    bool
	name    string
	bus     message.Bus
	modules []dep.Module
}

func (b *Builder) Apply(opts ...BuilderOption) error {
	b.Lock()
	defer b.Unlock()

	for _, opt := range opts {
		if err := opt(b); err != nil {
			return stacktrace.Propagate(err, "builder option returned non-nil error")
		}
	}
	return nil
}

func (b *Builder) Build() (*App, error) {
	b.Lock()
	defer b.Unlock()

	if b.bus == nil {
		b.bus = message.NewNoopBus()
	}

	container := dep.NewContainer(b.name)
	app := App{
		name:      b.name,
		container: container.Use(b.modules...),
		scheduler: scheduler.New(b.bus),
	}
	b.done = true
	return &app, nil
}

type BuilderOption func(builder *Builder) error

func newBuilder(name string) *Builder {
	return &Builder{
		name:    name,
		bus:     nil,
		modules: make([]dep.Module, 0),
	}
}
