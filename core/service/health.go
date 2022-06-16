package service

import (
	"sync"
	"time"

	"github.com/Raphy42/weekend/pkg/reflect"
)

type Health struct {
	Error     error
	LastCheck time.Time
}

type Registry struct {
	lock  sync.RWMutex
	last  chan Health
	inner map[string]*Health
}

func NewRegistry() *Registry {
	return &Registry{
		inner: make(map[string]*Health),
		last:  make(chan Health, 1),
	}
}

func (r *Registry) Set(service any, err error) {
	r.lock.Lock()
	defer r.lock.Unlock()

	name := reflect.Typename(service)

	now := time.Now()
	_, ok := r.inner[name]
	if !ok {
		r.inner[name] = &Health{}
	}
	r.inner[name].Error = err
	r.inner[name].LastCheck = now
	r.last <- *r.inner[name]
}

func (r *Registry) HealthProbe() <-chan Health {
	r.lock.RLock()
	defer r.lock.RUnlock()

	return r.last
}
