package dependency

import "github.com/heimdalr/dag"

type Graph struct {
	registry  *TypeRegistry
	factories []*Factory
	dag       *dag.DAG
}
