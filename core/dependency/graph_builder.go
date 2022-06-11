package dependency

import (
	"github.com/heimdalr/dag"
	"github.com/palantir/stacktrace"
	"github.com/rs/xid"

	"github.com/Raphy42/weekend/core/errors"
	"github.com/Raphy42/weekend/pkg/reflect"
	"github.com/Raphy42/weekend/pkg/std/slice"
)

type GraphBuilder struct {
	factories []interface{}
	registry  *TypeRegistry
}

func NewGraphBuilder() *GraphBuilder {
	return &GraphBuilder{
		factories: make([]interface{}, 0),
		registry:  NewTypeRegistry(),
	}
}

func (g *GraphBuilder) Factories(factories ...interface{}) *GraphBuilder {
	g.factories = append(g.factories, factories...)
	return g
}

func (g *GraphBuilder) Build() (*Graph, error) {
	graph := dag.NewDAG()

	factories, err := slice.MapErr(g.factories, func(_ int, in interface{}) (*Factory, error) {
		funcT, err := reflect.Func(in)
		if err != nil {
			return nil, stacktrace.Propagate(err, "only functions are accepted as factories: invalid '%T'", in)
		}

		factoryID := g.registry.Register(funcT.Value)
		_ = graph.AddVertexByID(factoryID.String(), funcT.String())
		return &Factory{
			id: factoryID,
			inputs: slice.Map(funcT.Ins, func(i int, in reflect.Type) xid.ID {
				id := g.registry.Register(in)
				_ = graph.AddVertexByID(id.String(), in.String())
				_ = graph.AddEdge(id.String(), factoryID.String())
				return id
			}),
			outputs: slice.Map(funcT.Outs, func(i int, out reflect.Type) xid.ID {
				id := g.registry.Register(out)
				_ = graph.AddVertexByID(id.String(), out.String())
				_ = graph.AddEdge(factoryID.String(), id.String())
				return id
			}),
		}, nil
	})
	if err != nil {
		return nil, stacktrace.Propagate(err, "could not initialise factories")
	}

	for _, factory := range factories {
		for _, in := range factory.inputs {
			for _, out := range factory.outputs {
				_ = graph.AddEdge(in.String(), out.String())
			}
		}
	}

	var errs errors.Group
	for name := range graph.GetRoots() {
		id, _ := xid.FromString(name)
		vi, ok := g.registry.Get(id)
		if !ok {
			return nil, stacktrace.NewError("unregistered typeID: %s", name)
		}

		switch v := vi.(type) {
		case reflect.Type:
			if !reflect.IsFunc(v) {
				errs.Append(stacktrace.NewError(
					"no provider was found for '%s': id was '%s'", v.String(), name,
				))
			}
		}

	}
	if err = errs.Coalesce(); err != nil {
		return nil, err
	}

	return &Graph{
		factories: factories,
		registry:  g.registry,
		dag:       graph,
	}, nil
}
