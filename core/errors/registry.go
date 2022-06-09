package errors

import "sync"

type Registry struct {
	sync.RWMutex
	domains map[uint16]string
	axioms  map[uint16]string
}

var globalRegistry = &Registry{
	domains: make(map[uint16]string),
	axioms:  make(map[uint16]string),
}

func RegisterDomain(flag uint16, name string) {
	globalRegistry.RegisterDomain(flag, name)
}

func RegisterAxiom(flag uint16, name string) {
	globalRegistry.RegisterAxiom(flag, name)
}

func (r *Registry) RegisterDomain(flag uint16, name string) *Registry {
	r.Lock()
	defer r.Unlock()

	r.domains[flag] = name
	return r
}

func (r *Registry) RegisterAxiom(flag uint16, name string) *Registry {
	r.Lock()
	defer r.Unlock()

	r.axioms[flag] = name
	return r
}

func (r *Registry) Domain(code uint16) string {
	r.RLock()
	defer r.RUnlock()

	domain, ok := r.domains[code]
	if !ok {
		return "<unknown>"
	}
	return domain
}

func (r *Registry) Axiom(code uint16) string {
	r.Lock()
	defer r.Unlock()

	axiom, ok := r.axioms[code]
	if !ok {
		return "<unknown>"
	}
	return axiom
}
