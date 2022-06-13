package service

import (
	"sync"
	"time"

	"github.com/Raphy42/weekend/pkg/reflect"
	"github.com/Raphy42/weekend/pkg/std/set"
)

type Health struct {
	Error     error
	LastCheck time.Time
}

type Registry struct {
	lock  sync.RWMutex
	inner map[string]Health
}

func NewRegistry() *Registry {
	return &Registry{
		inner: make(map[string]Health),
	}
}

func (r *Registry) Set(service any, err error) {
	r.lock.Lock()
	defer r.lock.Unlock()

	name := reflect.Typename(service)

	now := time.Now()
	health, ok := r.inner[name]
	if !ok {
		r.inner[name] = Health{
			Error:     err,
			LastCheck: now,
		}
		return
	}
	health.Error = err
	health.LastCheck = now
}

func (r *Registry) Health() map[string]Health {
	r.lock.RLock()
	defer r.lock.RUnlock()

	return set.Clone(r.inner)
}
