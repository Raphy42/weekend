package dep

import (
	"context"
	"sort"

	"github.com/heimdalr/dag"
	"github.com/palantir/stacktrace"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"

	"github.com/Raphy42/weekend/core"
	"github.com/Raphy42/weekend/core/errors"
	"github.com/Raphy42/weekend/core/logger"
	"github.com/Raphy42/weekend/pkg/reflect"
)

type Graph struct {
	registry *Registry
	inner    *dag.DAG
}

func NewGraph(dag *dag.DAG, registry *Registry) *Graph {
	return &Graph{inner: dag, registry: registry}
}

func (g *Graph) executeDependency(ctx context.Context, dependency *Dependency) error {
	ctx, span := otel.Tracer("Graph.executeDependency").Start(ctx, dependency.Name())
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
	var result interface{}
	if funcT.ReturnsResult() {
		result, err = funcT.CallResult(args...)
	} else if len(funcT.Outs) == 1 {
		result, err = funcT.Call(args...)
	}
	log.Debug("underlying function called via reflection")

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

func (g *Graph) childrenFactories(dependency *Dependency) ([]*Dependency, error) {
	deps := make([]*Dependency, 0)
	children, err := g.inner.GetOrderedDescendants(dependency.Name())
	if err != nil {
		return nil, err
	}
	for _, id := range children {
		dep, ok := g.registry.FindByName(id)
		if !ok {
			// todo error
			return nil, stacktrace.NewError("TODO")
		}
		if dep.Kind() == Factory {
			deps = append(deps, dep)
		}
	}
	return deps, nil
}

// Solve solves the dependency DAG and instantiate it
// todo: a working algorithm with no hacks
func (g *Graph) Solve(ctx context.Context) error {
	ctx, span := otel.Tracer(core.Name()).Start(ctx, "Graph.solve")
	defer span.End()

	log := logger.FromContext(ctx)
	registeredFactories := g.registry.Kind(Factory)

	sort.Slice(registeredFactories, func(i, j int) bool {
		a := registeredFactories[i]
		b := registeredFactories[j]

		aD, err := g.inner.GetOrderedAncestors(a.Name())
		span.RecordError(err)
		errors.Must(err)

		bD, err := g.inner.GetOrderedAncestors(b.Name())
		span.RecordError(err)
		errors.Must(err)
		return len(aD) < len(bD)
	})

	orderedFactories := registeredFactories
	for idx, factory := range orderedFactories {
		if err := g.executeDependency(ctx, factory); err != nil {
			span.RecordError(err)
			return stacktrace.Propagate(err, "factory '%s' execution (step %d of %d) failed", factory.Name(), idx+1, len(orderedFactories)+1)
		}
	}
	log.Info("all instances constructed")

	sideEffects := g.registry.Kind(SideEffect)
	for idx, sideEffect := range sideEffects {
		if err := g.executeDependency(ctx, sideEffect); err != nil {
			span.RecordError(err)
			return stacktrace.Propagate(err, "side effect '%s' execution (step %d of %d) failed", sideEffect.Name(), idx+1, len(sideEffects)+1)
		}
	}
	log.Info("all side effects executed")

	return nil
}
