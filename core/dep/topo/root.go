package topo

import (
	"github.com/Raphy42/weekend/pkg/set"
)

func (g *Graph) Roots() []string {
	return set.CollectSlice(g.inputs, func(k string, v int) (string, bool) {
		if v == 0 {
			return k, true
		}
		return "", false
	})
}
