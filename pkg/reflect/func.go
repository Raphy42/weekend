package reflect

import (
	"reflect"

	"github.com/palantir/stacktrace"
)

type FuncT struct {
	Value
	Ins  []Type
	Outs []Type
}

func (f FuncT) ReturnsResult() bool {
	return len(f.Outs) == 2 && f.Outs[1].Implements(ErrorType)
}

func (f FuncT) String() string {
	return Typename(f.Value.Interface())
}

func (f FuncT) Call(args ...Value) (any, error) {
	result := f.Value.Call(args)
	if len(result) > 1 {
		return result, stacktrace.NewError("invalid return arity, wanted <=1 got %d", len(result))
	}

	if len(result) == 0 {
		return nil, nil
	}
	if result[0].CanConvert(ErrorType) && !result[0].IsNil() {
		return nil, result[0].Convert(ErrorType).Interface().(error)
	}
	return result[0].Interface(), nil
}

func (f FuncT) CallResult(args ...Value) (any, error) {
	result := f.Value.Call(args)
	if len(result) != 2 {
		return result, stacktrace.NewError("invalid return arity, wanted 2 got %d", len(result))
	}
	if !result[1].CanConvert(ErrorType) {
		return result, stacktrace.NewError("invalid second return type, wanted error got %s", result[1])
	}

	err := result[1]
	if err.IsNil() {
		return result[0].Interface(), nil
	}

	return result[0].Interface(), result[1].Convert(ErrorType).Interface().(error)
}

func Func(fn any) (*FuncT, error) {
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
