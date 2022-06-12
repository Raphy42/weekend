package dep

import (
	"fmt"
	"sync"

	"github.com/palantir/stacktrace"
	"github.com/rs/xid"

	"github.com/Raphy42/weekend/core/errors"
	"github.com/Raphy42/weekend/pkg/reflect"
)

type Kind int

const (
	//Factory represent a dependency which can instantiate new dependencies.
	//This is the entrypoint to the whole DI system and the solver will try to associate every declare Instance with
	//their own Factory.
	//Do not declare a Factory which uses scalar types like `int32` or `[]string{}`, always use opaque types or rethink
	//your implementation.
	//Factory should be implemented with a pure function, with no side-effects. If you need to mutate external state
	//use SideEffect instead.
	//Factory should only be used to declare heavy-lifter objects (services, drivers, resources) which have complex
	//lifecycles, or depends on each other.
	Factory Kind = iota
	//SideEffect is a dependency which should be invoked whenever the DI system has finished bootstrapping.
	SideEffect
	//Transitive represent intermediate dependency state.
	//For example a factory function outputting (*T, error) will be represented by a transitive dependency of value *T.
	//The transitive will then be associated with an Instance or a Token depending on its resolution at runtime.
	Transitive
	//Instance represent the actual value of the dependency, used and created by both factories and transitives
	Instance
	//Token is a placeholder dependency, used to inject custom synchronisation logic in the DAG dependency solver.
	Token
)

type Status int

const (
	NewStatus Status = iota
	InitialisedStatus
	CorruptedStatus
)

type Dependency struct {
	lock      sync.RWMutex
	id        xid.ID
	kind      Kind
	status    Status
	value     interface{}
	lastError error
}

func newDependency(kind Kind, value interface{}) *Dependency {
	status := NewStatus
	// factories should not need to be initialised
	if kind == Factory {
		status = InitialisedStatus
	}
	return &Dependency{
		id:     xid.New(),
		kind:   kind,
		value:  value,
		status: status,
	}
}

func NewFactory(value interface{}) (*Dependency, error) {
	funcT, err := reflect.Func(value)
	if err != nil {
		return nil, stacktrace.Propagate(err, "invalid factory type")
	}
	return newDependency(Factory, funcT), nil
}

type Transition struct {
	sync.RWMutex
	From  xid.ID
	Value interface{}
}

func (t *Transition) String() string {
	t.RLock()
	defer t.RUnlock()

	return fmt.Sprintf("Transition.%s.%s", reflect.Typename(t.Value), t.From)
}

func NewResult(from xid.ID) (*Dependency, error) {
	return newDependency(Transitive, &Transition{Value: nil, From: from}), nil
}

func NewSideEffect(value interface{}) (*Dependency, error) {
	funcT, err := reflect.Func(value)
	if err != nil {
		return nil, stacktrace.Propagate(err, "invalid side effect type")
	}
	return newDependency(SideEffect, funcT), nil
}

type InstanceContainer struct {
	sync.RWMutex
	Type  reflect.Type
	Value interface{}
}

func (i *InstanceContainer) String() string {
	i.RLock()
	defer i.RUnlock()

	return fmt.Sprintf("%s", i.Type)
}

func NewInstance(kind reflect.Type) (*Dependency, error) {
	return newDependency(Instance, &InstanceContainer{
		Type: kind,
	}), nil
}

func NewToken(value string) (*Dependency, error) {
	return newDependency(Token, &value), nil
}

func (d *Dependency) Solve(value interface{}, err error) {
	d.lock.Lock()
	defer d.lock.Unlock()

	d.lastError = err
	if d.kind == Transitive || d.kind == Instance {
		if err != nil {
			d.status = CorruptedStatus
		} else {
			if d.kind == Instance {
				d.value.(*InstanceContainer).Value = value
			}
			d.status = InitialisedStatus
		}
	}
}

func (d *Dependency) Status() Status {
	return d.status
}

func (d *Dependency) Value() (interface{}, error) {
	if d.lastError != nil && d.status != InitialisedStatus {
		return nil, stacktrace.Propagate(d.lastError, "dependency is in an invalid state")
	}
	if d.Kind() == Instance {
		return d.value.(*InstanceContainer).Value, nil
	}

	return d.value, nil
}

func (d *Dependency) ID() xid.ID {
	d.lock.RLock()
	defer d.lock.RUnlock()

	return d.id
}

func (d *Dependency) HasIO() bool {
	kind := d.Kind()
	return kind == SideEffect || kind == Factory
}

func (d *Dependency) Name() string {
	d.lock.RLock()
	defer d.lock.RUnlock()

	//todo refactor this mess with fmt.Stringer or generics
	switch d.kind {
	case Token:
		v := d.value.(*string)
		return *v
	case Instance:
		v := d.value.(*InstanceContainer)
		return v.String()
	case SideEffect, Factory:
		v := d.value.(*reflect.FuncT)
		return v.String()
	case Transitive:
		v := d.value.(*Transition)
		return v.String()
	}
	panic(stacktrace.NewErrorWithCode(errors.EUnreachable, "unknown dependency kind"))
}

func (d *Dependency) Kind() Kind {
	d.lock.RLock()
	defer d.lock.RUnlock()

	return d.kind
}
