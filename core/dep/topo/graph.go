package topo

// toposort and behavior based on https://github.com/philopon/go-toposort
// MIT License

type Graph struct {
	nodes   []string
	inputs  map[string]int
	outputs map[string]map[string]int
}

func New() *Graph {
	return &Graph{
		nodes:   make([]string, 0),
		inputs:  make(map[string]int),
		outputs: make(map[string]map[string]int),
	}
}

func (g *Graph) Clone() *Graph {
	return &(*g)
}

func (g *Graph) AddNode(name string) bool {
	g.nodes = append(g.nodes, name)

	if _, ok := g.outputs[name]; ok {
		return false
	}
	g.outputs[name] = make(map[string]int)
	g.inputs[name] = 0
	return true
}

func (g *Graph) AddNodes(names ...string) bool {
	for _, name := range names {
		if ok := g.AddNode(name); !ok {
			return false
		}
	}
	return true
}

func (g *Graph) AddEdge(from, to string) bool {
	m, ok := g.outputs[from]
	if !ok {
		return false
	}

	m[to] = len(m) + 1
	g.inputs[to]++

	return true
}

func (g *Graph) unsafeRemoveEdge(from, to string) {
	delete(g.outputs[from], to)
	g.inputs[to]--
}

func (g *Graph) RemoveEdge(from, to string) bool {
	if _, ok := g.outputs[from]; !ok {
		return false
	}
	g.unsafeRemoveEdge(from, to)
	return true
}
