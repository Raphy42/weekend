package core

import (
	"sync"

	"github.com/palantir/stacktrace"
)

type lifecycles struct {
	sync.RWMutex
	onStart []func() error
	onStop  []func() error
}

func newLifecycles() *lifecycles {
	return &lifecycles{
		onStart: make([]func() error, 0),
		onStop:  make([]func() error, 0),
	}
}

var globalLifecycles = newLifecycles()

func RegisterOnStartHook(fn func() error) {
	globalLifecycles.Lock()
	defer globalLifecycles.Unlock()

	globalLifecycles.onStart = append(globalLifecycles.onStart, fn)
}

func RegisterOnStopHook(fn func() error) {
	globalLifecycles.Lock()
	defer globalLifecycles.Unlock()

	globalLifecycles.onStop = append(globalLifecycles.onStop, fn)
}

func ExecuteOnStart() error {
	globalLifecycles.RLock()
	defer globalLifecycles.RUnlock()

	for _, fn := range globalLifecycles.onStart {
		if err := fn(); err != nil {
			return stacktrace.Propagate(err, "onStart global handler failed")
		}
	}
	return nil
}

func ExecuteOnStop() error {
	globalLifecycles.RLock()
	defer globalLifecycles.RUnlock()

	for _, fn := range globalLifecycles.onStop {
		if err := fn(); err != nil {
			return stacktrace.Propagate(err, "onStop global handler failed")
		}
	}
	return nil
}
