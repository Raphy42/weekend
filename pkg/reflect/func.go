package reflect

import (
	"fmt"
	"reflect"

	"github.com/palantir/stacktrace"
)

type FuncT struct {
	Value
	Ins  []Type
	Outs []Type
}

func (f FuncT) ReturnsResult() bool {
	return len(f.Outs) == 2 && f.Outs[1].Implements(reflect.TypeOf(fmt.Errorf("")))
}

func Func(fn interface{}) (*FuncT, error) {
	t := reflect.ValueOf(fn)
	if t.Kind() == reflect.Pointer {
		return Func(t.Elem().Interface())
	}
	if t.Kind() != reflect.Func {
		return nil, stacktrace.NewError("unexpected %T in reflect.Func call", fn)
	}

	tt := t.Type()
	ins := make([]Type, tt.NumIn())
	outs := make([]Type, tt.NumOut())

	for i := 0; i < tt.NumIn(); i++ {
		ins[i] = tt.In(i)
	}

	for i := 0; i < tt.NumOut(); i++ {
		outs[i] = tt.Out(i)
	}

	return &FuncT{
		Value: t,
		Ins:   ins,
		Outs:  outs,
	}, nil
}
