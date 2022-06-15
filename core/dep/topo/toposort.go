package topo

func (g *Graph) TopologicalSort() ([]string, bool) {
	// work on a clone instead of the actual underlying value
	self := g.Clone()

	orderedNodes := make([]string, 0, len(self.nodes))
	s := make([]string, 0, len(self.nodes))

	for _, n := range self.nodes {
		if self.inputs[n] == 0 {
			s = append(s, n)
		}
	}

	for len(s) > 0 {
		var n string
		n, s = s[0], s[1:]
		orderedNodes = append(orderedNodes, n)

		ms := make([]string, len(self.outputs[n]))
		for m, i := range self.outputs[n] {
			ms[i-1] = m
		}

		for _, m := range ms {
			self.unsafeRemoveEdge(n, m)

			if self.inputs[m] == 0 {
				s = append(s, m)
			}
		}
	}

	nodeCount := 0
	for _, v := range self.inputs {
		nodeCount += v
	}

	if nodeCount > 0 {
		return orderedNodes, false
	}

	return orderedNodes, true
}
