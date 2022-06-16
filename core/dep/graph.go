package dep

import (
	"context"

	"github.com/palantir/stacktrace"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"

	"github.com/Raphy42/weekend/core/dep/topo"
	"github.com/Raphy42/weekend/core/errors"
	"github.com/Raphy42/weekend/core/logger"
	"github.com/Raphy42/weekend/pkg/reflect"
)

type Graph struct {
	registry *Registry
	topo     *topo.Graph
}

func NewGraph(topo *topo.Graph, registry *Registry) *Graph {
	return &Graph{
		registry: registry,
		topo:     topo,
	}
}

func (g *Graph) executeDependency(ctx context.Context, dependency *Dependency) error {
	ctx, span := otel.Tracer("wk.core.dep").Start(ctx, dependency.Name())
	defer span.End()

	kindS := "factory"
	if dependency.Kind() == SideEffect {
		kindS = "side-effect"
	}
	span.SetAttributes(
		attribute.String("wk.dependency.kind", kindS),
		attribute.Stringer("wk.dependency.id", dependency.id),
		attribute.String("wk.dependency.name", dependency.Name()),
	)

	log := logger.FromContext(ctx).With(
		zap.String("wk.dependency.name", dependency.Name()),
		zap.String("wk.dependency.kind", kindS),
	)
	log.Debug("executing")

	value, err := dependency.Value()
	errors.Mustf(err, "%s '%s' has no value", kindS, dependency.Name())

	// fetch arguments from registry
	args := make([]reflect.Value, 0)
	funcT := value.(*reflect.FuncT)
	for idx, in := range funcT.Ins {
		instance, found := g.registry.FindByName(in.String())
		if !found {
			return stacktrace.NewError(
				"no dependency named '%s' in registry for factory '%s'", in.String(), dependency.Name(),
			)
		}
		if instance.Status() != InitialisedStatus {
			return stacktrace.NewError(
				"instance '%s', argument number %d has not been initialised prior to factory '%s' invocation",
				instance.Name(), idx, dependency.Name(),
			)
		}
		if instance.Name() == "context.Context" {
			instance.Solve(ctx, nil)
		}

		argValue, err := instance.Value()
		if err != nil {
			return stacktrace.Propagate(err,
				"dependency number %d of factory '%s' has invalid value", idx, dependency,
			)
		}
		args = append(args, reflect.ValueOf(argValue))
	}

	// two choices
	// (*T, error)		-> FuncT.CallResult
	// *T & the rest	-> FuncT.Call
	var result any
	if funcT.ReturnsResult() {
		result, err = funcT.CallResult(args...)
	} else if len(funcT.Outs) == 1 {
		result, err = funcT.Call(args...)
	}

	if len(funcT.Outs) == 0 && dependency.Kind() == Factory {
		log.Warn("factory has no outputs, it should be declared a side effect instead")
		return nil
	}

	if dependency.Kind() == SideEffect {
		return stacktrace.Propagate(err, "%s '%s' invocation failed", kindS, dependency.Name())
	}

	// get output dependency associated with factory
	out := funcT.Outs[0]
	outDependency, ok := g.registry.FindByName(out.String())
	if !ok {
		return stacktrace.NewError(
			"unable to find output dependency named '%s' in the registry", outDependency.Name(),
		)
	}

	outDependency.Solve(result, err)
	return stacktrace.Propagate(err, "%s '%s' invocation failed", kindS, dependency.Name())
}

// Solve solves the dependency DAG and instantiate it
// todo: check that toposort works
func (g *Graph) Solve(ctx context.Context) error {
	ctx, span := otel.Tracer("wk.core.dep").Start(ctx, "Graph.solve")
	defer span.End()

	log := logger.FromContext(ctx)

	sortedNodes, ok := g.topo.TopologicalSort()
	if !ok {
		return stacktrace.NewError("topological sort failed")
	}
	orderedFactories := make([]*Dependency, 0)
	for _, node := range sortedNodes {
		dependency, ok := g.registry.FindByName(node)
		if !ok {
			return stacktrace.NewError("dependency should be registered at this point")
		}
		if dependency.Kind() == Factory {
			orderedFactories = append(orderedFactories, dependency)
		}
	}

	for idx, factory := range orderedFactories {
		if err := g.executeDependency(ctx, factory); err != nil {
			span.RecordError(err)
			return stacktrace.Propagate(err, "factory '%s' execution (step %d of %d) failed", factory.Name(), idx+1, len(orderedFactories)+1)
		}
	}
	log.Info("all instances constructed", zap.Int("wk.factory.count", len(orderedFactories)))

	sideEffects := g.registry.Kind(SideEffect)
	for idx, sideEffect := range sideEffects {
		if err := g.executeDependency(ctx, sideEffect); err != nil {
			span.RecordError(err)
			return stacktrace.Propagate(err, "side effect '%s' execution (step %d of %d) failed", sideEffect.Name(), idx+1, len(sideEffects)+1)
		}
	}
	log.Info("all side effects executed", zap.Int("wk.side_effect.count", len(sideEffects)))

	return nil
}
