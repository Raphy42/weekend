package scheduler

import (
	"context"
	"reflect"

	"github.com/google/uuid"
	"github.com/palantir/stacktrace"

	"github.com/Raphy42/weekend/core/errors"
)

type Schedulable interface {
	func() |
		func() error |
		func(ctx context.Context) |
		func(ctx context.Context) error |
		func(ctx context.Context, args interface{}) |
		func(ctx context.Context, args interface{}) error |
		func(ctx context.Context, args interface{}) (interface{}, error)
}

type SchedulableFn func(ctx context.Context, args interface{}) (interface{}, error)
type Manifest struct {
	Name string
	ID   uuid.UUID
	Fn   SchedulableFn
}

// oooooooh this is dirty
// too bad golang generics don't allow specialisation :(
// but it's less dirty than ye big ol'trusty `interface{}`
func makeImpl[S Schedulable](schedulable S) SchedulableFn {
	v := reflect.ValueOf(schedulable).Interface()
	switch fn := v.(type) {
	case func():
		return func(ctx context.Context, args interface{}) (interface{}, error) {
			fn()
			return nil, nil
		}
	case func() error:
		return func(ctx context.Context, args interface{}) (interface{}, error) {
			return nil, fn()
		}
	case func(ctx context.Context):
		return func(ctx context.Context, args interface{}) (interface{}, error) {
			fn(ctx)
			return nil, nil
		}
	case func(ctx context.Context) error:
		return func(ctx context.Context, args interface{}) (interface{}, error) {
			return nil, fn(ctx)
		}
	case func(ctx context.Context, args interface{}) error:
		return func(ctx context.Context, args interface{}) (interface{}, error) {
			return nil, fn(ctx, args)
		}
	case func(ctx context.Context, args interface{}) (interface{}, error):
		return fn
	default:
		panic(
			stacktrace.NewMessageWithCode(errors.EUnreachable, "invalid Schedulable '%T'", fn),
		)
	}
}

func Make[S Schedulable](name string, schedulable S) Manifest {
	return Manifest{
		Name: name,
		Fn:   makeImpl(schedulable),
		ID:   uuid.New(),
	}
}
