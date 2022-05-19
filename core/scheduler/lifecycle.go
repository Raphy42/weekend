package scheduler

import (
	"context"
	"sync"
	"time"
)

type lifecycleTmpContext struct {
	context.Context
	cancel context.CancelFunc
}

func newLifecycleContext(parentCtx context.Context, duration time.Duration) lifecycleTmpContext {
	ctx, cancel := context.WithTimeout(parentCtx, duration)
	return lifecycleTmpContext{
		Context: ctx,
		cancel:  cancel,
	}
}

type Lifecycle int

const (
	//LBackground will run declared lifecycles for the whole duration of the application.
	// Use this lifecycle if you want to declare background jobs, as the other hooks will use a temporary
	// `context.Context` with a small timeout, thus killing any running lifecycle at the end of their respective stages.
	LBackground Lifecycle = iota
	//LOnStart will run declared lifecycles at startup invocation stage
	LOnStart
	//LOnStop will run declared lifecycles at shutdown invocation stage
	LOnStop
	//LOnCrash will run declared lifecycles at crash invocation stage
	LOnCrash
)

//LifecycleHook is an helper type for lifecycle hooks
type LifecycleHook func(ctx context.Context) error

type lifecycleStore map[Lifecycle][]LifecycleHook

func newLifecycleStore() lifecycleStore {
	return make(lifecycleStore)
}

func (l lifecycleStore) AppendHook(kind Lifecycle, hook LifecycleHook) {
	_, ok := l[kind]
	if !ok {
		l[kind] = make([]LifecycleHook, 0)
	}
	l[kind] = append(l[kind], hook)
}

//Lifecycles is the user entrypoint to the underlying lifecycle hook system.
// See the `Lifecycle` enum for further explanation of the whole system.
type Lifecycles struct {
	sync.RWMutex
	registered lifecycleStore
}

//Lifecycle allows user to register callbacks/hooks to be executed whenever their module have been initialised.
// Calling this method anywhere else other than inside a provider is undefined behavior.
func (l *Lifecycles) Lifecycle(kind Lifecycle, hooks ...LifecycleHook) *Lifecycles {
	for _, hook := range hooks {
		l.registered.AppendHook(kind, hook)
	}
	return l
}

func (l *Lifecycles) hooks(kind Lifecycle) []LifecycleHook {
	hooks := make([]LifecycleHook, 0)
	for _, lifecycle := range l.registered[kind] {
		hooks = append(hooks, lifecycle)
	}
	return hooks
}

func (l *Lifecycles) startHooks() []LifecycleHook {
	return l.hooks(LOnStart)
}

func (l *Lifecycles) stopHooks() []LifecycleHook {
	return l.hooks(LOnStop)
}

func (l *Lifecycles) crashHooks() []LifecycleHook {
	return l.hooks(LOnCrash)
}

func (l *Lifecycles) backgroundHooks() []LifecycleHook {
	return l.hooks(LBackground)
}
