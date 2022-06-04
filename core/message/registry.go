package message

import (
	"sync"

	"github.com/palantir/stacktrace"
	"go.uber.org/atomic"

	"github.com/Raphy42/weekend/core/errors"
)

type RegistrySlot struct {
	next  *atomic.Uint32
	inner map[uint32]Handler
}

func newRegistrySlot() *RegistrySlot {
	return &RegistrySlot{
		next:  atomic.NewUint32(0),
		inner: make(map[uint32]Handler),
	}
}

func (r *RegistrySlot) Register(handler Handler) uint32 {
	slot := r.next.Inc()
	r.inner[slot] = handler
	return slot
}

func (r *RegistrySlot) Unregister(slot uint32) error {
	_, ok := r.inner[slot]
	if ok {
		delete(r.inner, slot)
		return nil
	}
	return stacktrace.NewErrorWithCode(
		errors.EUnreachable,
		"trying to unregister a handler with an invalid slot: %d", slot,
	)
}

type Registry struct {
	sync.RWMutex
	handlers map[string]*RegistrySlot
}

func (r *Registry) Register(kind string, handler Handler) uint32 {
	r.Lock()
	defer r.Unlock()

	_, ok := r.handlers[kind]
	if !ok {
		r.handlers[kind] = newRegistrySlot()
	}
	return r.handlers[kind].Register(handler)
}

func (r *Registry) Unregister(kind string, slot uint32) error {
	r.Lock()
	defer r.Unlock()

	_, ok := r.handlers[kind]
	if !ok {
		return stacktrace.NewErrorWithCode(
			errors.EUnreachable,
			"trying to unregister an unknown handler kind: '%s'", kind,
		)
	}

	if err := r.handlers[kind].Unregister(slot); err != nil {
		return stacktrace.Propagate(
			err,
			"handler has not been previously registered or slot is invalid: '%s'",
			kind,
		)
	}
	return nil
}
