package concurrent_set

import (
	"sync"

	"github.com/Raphy42/weekend/pkg/set"
)

type Set[K comparable, V any] struct {
	lock  sync.RWMutex
	inner map[K]V
}

func New[K comparable, V any]() Set[K, V] {
	return Set[K, V]{
		inner: make(map[K]V),
	}
}

func From[K comparable, V any](values map[K]V) Set[K, V] {
	return Set[K, V]{
		inner: values,
	}
}

func (s *Set[K, V]) Insert(key K, value V) {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.inner[key] = value
}

func (s *Set[K, V]) Get(key K) (V, bool) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	v, ok := s.inner[key]
	return v, ok
}

func (s *Set[K, V]) Values() []V {
	s.lock.RLock()
	defer s.lock.RUnlock()

	return set.Values(s.inner)
}

func (s *Set[K, V]) Keys() []K {
	s.lock.RLock()
	defer s.lock.RUnlock()

	return set.Keys(s.inner)
}

func (s *Set[K, V]) Delete(key K) {
	s.lock.Lock()
	defer s.lock.Unlock()

	delete(s.inner, key)
}

func (s *Set[K, V]) Edit(key K, fn func(value *V) *V) {
	s.lock.Lock()
	defer s.lock.Unlock()

	v, ok := s.inner[key]
	var newV *V
	if !ok {
		newV = fn(nil)
	} else {
		newV = fn(&*&v)
	}
	if newV != nil {
		s.inner[key] = *newV
	}
}

func (s *Set[K, V]) Iter() map[K]V {
	s.lock.RLock()
	defer s.lock.RUnlock()

	return *&s.inner
}
