package dep

import (
	"sync"

	"github.com/rs/xid"

	"github.com/Raphy42/weekend/pkg/std/slice"
)

type Registry struct {
	lock         sync.RWMutex
	dependencies map[xid.ID]*Dependency
}

func NewRegistry() *Registry {
	reg := Registry{
		dependencies: make(map[xid.ID]*Dependency),
	}

	return &reg
}

func (r *Registry) Get(id xid.ID) (*Dependency, bool) {
	r.lock.Lock()
	defer r.lock.RUnlock()

	v, ok := r.dependencies[id]
	return v, ok
}

func (r *Registry) FindByName(name string) (*Dependency, bool) {
	r.lock.RLock()
	defer r.lock.RUnlock()

	for _, dep := range r.dependencies {
		if dep.Name() == name {
			return dep, true
		}
	}
	return nil, false
}

func (r *Registry) Lookup(dependency *Dependency) (xid.ID, bool) {
	r.lock.RLock()
	defer r.lock.RUnlock()

	for id := range r.dependencies {
		if dependency.ID() == id {
			return id, true
		}
	}
	return xid.NilID(), false
}

func (r *Registry) Register(dependency *Dependency) error {
	_, found := r.Lookup(dependency)
	if found {
		return nil
	}

	r.lock.Lock()
	defer r.lock.Unlock()

	r.dependencies[dependency.ID()] = dependency
	return nil
}

var (
	allDependencyKinds = slice.New(Factory, Instance, SideEffect)
)

func (r *Registry) Kind(kinds ...Kind) []*Dependency {
	r.lock.RLock()
	defer r.lock.RUnlock()

	deps := make([]*Dependency, 0)
	for _, dep := range r.dependencies {
		if slice.Contains(kinds, dep.Kind()) {
			deps = append(deps, dep)
		}
	}
	return deps
}
