package dep

import (
	"context"

	"github.com/palantir/stacktrace"
	"go.opentelemetry.io/otel"

	"github.com/Raphy42/weekend/core"
	"github.com/Raphy42/weekend/core/scheduler/schedulable"
)

type Container struct {
	name    string
	modules []Module
	graph   *Graph
}

func NewContainer(name string) *Container {
	return &Container{
		name:    name,
		modules: make([]Module, 0),
	}
}

func (c *Container) start(ctx context.Context) (interface{}, error) {
	ctx, span := otel.Tracer(core.Name()).Start(ctx, "Container.start")
	defer span.End()

	graph, err := NewGraphBuilder().
		Build(ctx, c.modules...)
	if err != nil {
		return nil, stacktrace.Propagate(err, "construction of dependency solver failed")
	}
	c.graph = graph
	return nil, c.graph.Solve(ctx)
}

func (c *Container) Use(modules ...Module) *Container {
	c.modules = append(c.modules, modules...)
	return c
}

func (c *Container) Manifest() schedulable.Manifest {
	return schedulable.Of(
		schedulable.Name("container", c.name),
		c.start,
	)
}
