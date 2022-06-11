package dependency

import (
	"sync"

	"github.com/rs/xid"
)

type TypeRegistry struct {
	sync.RWMutex
	inner map[xid.ID]interface{}
	lut   map[interface{}]xid.ID
}

func NewTypeRegistry() *TypeRegistry {
	return &TypeRegistry{
		inner: make(map[xid.ID]interface{}),
		lut:   make(map[interface{}]xid.ID),
	}
}

func (t *TypeRegistry) Get(id xid.ID) (interface{}, bool) {
	t.RLock()
	defer t.RUnlock()

	v, ok := t.inner[id]
	return v, ok
}

func (t *TypeRegistry) Lookup(kind interface{}) (xid.ID, bool) {
	t.RLock()
	defer t.RUnlock()

	v, ok := t.lut[kind]
	return v, ok
}

func (t *TypeRegistry) Register(kind interface{}) xid.ID {
	if id, ok := t.Lookup(kind); ok {
		return id
	}

	t.Lock()
	defer t.Unlock()

	id := xid.New()
	t.inner[id] = kind
	t.lut[kind] = id
	return id
}
