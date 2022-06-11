package dependency

import (
	"context"

	"github.com/palantir/stacktrace"

	"github.com/Raphy42/weekend/core/di"
	"github.com/Raphy42/weekend/core/scheduler/schedulable"
)

type Container struct {
	name    string
	modules []di.Module
	graph   *Graph
}

func NewContainer(name string) *Container {
	return &Container{
		name:    name,
		modules: make([]di.Module, 0),
	}
}

func (c *Container) start(ctx context.Context) (interface{}, error) {
	builder := NewGraphBuilder()
	for _, module := range c.modules {
		builder.Factories(module.Providers...)
	}

	graph, err := builder.Build()
	if err != nil {
		return c, stacktrace.Propagate(err, "could not build dependency graph")
	}
	c.graph = graph
	return c, err
}

func (c *Container) Use(modules ...di.Module) *Container {
	c.modules = append(c.modules, modules...)
	return c
}

func (c *Container) Manifest() schedulable.Manifest {
	return schedulable.Of(
		schedulable.Name("container", c.name),
		c.start,
	)
}
