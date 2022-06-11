package reflect

import (
	"reflect"

	"github.com/palantir/stacktrace"

	"github.com/Raphy42/weekend/pkg/runtime"
)

// convenience reexports
type (
	Type       = reflect.Type
	Value      = reflect.Value
	Reflection struct {
		Type
		Value
	}
)

var (
	TypeOf    = reflect.TypeOf
	ValueOf   = reflect.ValueOf
	UnsafeNew = reflect.New
	PtrTo     = reflect.PtrTo
)

func Zero(t reflect.Type) (value reflect.Value, err error) {
	if t.Kind() == reflect.Pointer {
		return Zero(t.Elem())
	}
	err = runtime.Recover(func() {
		value = reflect.Zero(t)
	})
	if err != nil {
		err = stacktrace.Propagate(err, "could not create zero value of type '%s'", t.String())
	}
	return
}

func IsFunc(t Type) bool {
	return t.Kind() == reflect.Func
}
