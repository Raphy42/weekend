package dep

import (
	"context"
	"strings"

	"github.com/palantir/stacktrace"
	"go.opentelemetry.io/otel"
	"go.uber.org/zap"

	"github.com/Raphy42/weekend/core/dep/topo"
	"github.com/Raphy42/weekend/core/errors"
	"github.com/Raphy42/weekend/core/logger"
	"github.com/Raphy42/weekend/pkg/reflect"
	"github.com/Raphy42/weekend/pkg/slice"
)

type GraphBuilder struct {
	registry             *Registry
	factoryCandidates    []any
	sideEffectCandidates []any
}

func NewGraphBuilder() *GraphBuilder {
	return &GraphBuilder{
		registry:             NewRegistry(),
		factoryCandidates:    make([]any, 0),
		sideEffectCandidates: make([]any, 0),
	}
}

func (g *GraphBuilder) Build(ctx context.Context, modules ...Module) (*Graph, error) {
	ctx, span := otel.Tracer("wk.core.dep").Start(ctx, "GraphBuilder.Build")
	defer span.End()

	log := logger.FromContext(ctx)

	// group by stage
	for _, module := range modules {
		log.Info("loading module", zap.String("wk.module", module.Name))
		g.factoryCandidates = append(g.factoryCandidates, module.Factories...)
		g.sideEffectCandidates = append(g.sideEffectCandidates, module.SideEffects...)
	}

	// register factories
	for _, candidate := range g.factoryCandidates {
		factory, err := NewFactory(candidate)
		if err != nil {
			return nil, stacktrace.Propagate(err, "candidate was not a valid factory")
		}
		if err = g.registry.Register(factory); err != nil {
			return nil, stacktrace.Propagate(err, "could not register factory")
		}
	}

	// register side effects
	for _, candidate := range g.sideEffectCandidates {
		sideEffect, err := NewSideEffect(candidate)
		if err != nil {
			return nil, stacktrace.Propagate(err, "candidate was not a valid side effect")
		}
		if err = g.registry.Register(sideEffect); err != nil {
			return nil, stacktrace.Propagate(err, "could not register side effect")
		}
	}

	edges := make(map[string][]string)

	// get all input and output types from factories and side effects
	candidateMap := make(map[string]reflect.Type)
	for _, dependency := range g.registry.Kind(Factory, SideEffect) {
		value, err := dependency.Value()
		if err != nil {
			return nil, stacktrace.Propagate(err, "dependency has no valid value")
		}
		funcT := value.(*reflect.FuncT)
		for _, in := range funcT.Ins {
			candidateMap[in.String()] = in
			edges[in.String()] = slice.Append(edges[in.String()], dependency.Name())
		}
		for _, out := range funcT.Outs {
			candidateMap[out.String()] = out
			edges[dependency.Name()] = slice.Append(edges[dependency.Name()], out.String())
		}
	}
	log.Debug("instance candidates extracted", zap.Int("count", len(candidateMap)))

	// register all instance candidates as dependencies
	for _, candidate := range candidateMap {
		dependency, err := NewInstance(candidate)
		if err != nil {
			return nil, stacktrace.Propagate(err, "unable to create instance dependency")
		}
		if err = g.registry.Register(dependency); err != nil {
			return nil, stacktrace.Propagate(err, "unable to register instance")
		}
	}

	allDeps := g.registry.Kind(allDependencyKinds...)
	topoGraph := topo.New()

	// add all vertices to DAG
	for _, dependency := range allDeps {
		topoGraph.AddNode(dependency.Name())
	}

	// add all edges to DAG
	for from, tos := range edges {
		for _, to := range tos {
			topoGraph.AddEdge(from, to)
		}
	}

	// validate the graph by confirming that all root dependencies are factories
	var errs errors.Group
	for _, dependency := range topoGraph.Roots() {
		dep, ok := g.registry.FindByName(dependency)
		if !ok {
			panic("invalid dependency")
		}
		if !dep.HasIO() {
			name := dep.Name()
			if !strings.HasPrefix(name, "*") {
				errs.Append(stacktrace.NewError(
					"no registered factory can construct '%s', did you mean '*%s' ?", dep.Name(), dep.Name()),
				)
			} else {
				errs.Append(stacktrace.NewError(
					"no registered factory can construct '%s'", dep.Name(), dep.Name()),
				)
			}
		}
	}

	return NewGraph(topoGraph, g.registry), errs.Coalesce()
}
