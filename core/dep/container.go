package dep

import (
	"context"

	"github.com/palantir/stacktrace"
	"go.opentelemetry.io/otel"

	"github.com/Raphy42/weekend/core/scheduler/async"
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

func (c *Container) start(ctx context.Context) (any, error) {
	ctx, span := otel.Tracer("wk.core.dep").Start(ctx, "Container.start")

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

func (c *Container) Manifest() async.Manifest {
	return async.Of(
		async.Name("container", c.name),
		c.start,
	)
}

func (c *Container) UnsafeExecute(ctx context.Context, fn any) error {
	if c.graph == nil {
		return stacktrace.NewError("application has not been initialised")
	}
	return c.graph.UnsafeExecute(ctx, fn)
}
