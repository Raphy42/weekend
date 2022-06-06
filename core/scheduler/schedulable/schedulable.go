package schedulable

import (
	"context"
	"reflect"

	"github.com/palantir/stacktrace"
	"github.com/rs/xid"

	"github.com/Raphy42/weekend/core/errors"
	"github.com/Raphy42/weekend/core/scheduler/policies"
)

type Schedulable interface {
	func() |
	func() error |
	func(ctx context.Context) |
	func(ctx context.Context) error |
	func(ctx context.Context) (interface{}, error) |
	func(ctx context.Context, args interface{}) |
	func(ctx context.Context, args interface{}) error |
	func(ctx context.Context, args interface{}) (interface{}, error)
}

type Fn func(ctx context.Context, args interface{}) (interface{}, error)
type Manifest struct {
	Name   string
	ID     xid.ID
	Fn     Fn
	Policy policies.Policy
}

// oooooooh this is dirty
// too bad golang generics don't allow specialisation :(
// but it's less dirty than ye big ol'trusty `interface{}`
func makeImpl[S Schedulable](schedulable S) Fn {
	v := reflect.ValueOf(schedulable).Interface()
	switch fn := v.(type) {
	case func():
		return func(_ context.Context, args interface{}) (interface{}, error) {
			fn()
			return args, nil
		}
	case func() error:
		return func(_ context.Context, args interface{}) (interface{}, error) {
			return args, fn()
		}
	case func(ctx context.Context):
		return func(ctx context.Context, args interface{}) (interface{}, error) {
			fn(ctx)
			return args, nil
		}
	case func(ctx context.Context) error:
		return func(ctx context.Context, args interface{}) (interface{}, error) {
			return args, fn(ctx)
		}
	case func(ctx context.Context) (interface{}, error):
		return func(ctx context.Context, _ interface{}) (interface{}, error) {
			return fn(ctx)
		}
	case func(ctx context.Context, args interface{}) error:
		return func(ctx context.Context, args interface{}) (interface{}, error) {
			return args, fn(ctx, args)
		}
	case func(ctx context.Context, args interface{}) (interface{}, error):
		return fn
	default:
		panic(
			stacktrace.NewMessageWithCode(errors.EUnreachable, "invalid Schedulable '%T'", fn),
		)
	}
}

func Make[S Schedulable](name string, schedulable S, pols ...policies.Policy) Manifest {
	policy := policies.Default()
	if len(pols) != 0 {
		policy = pols[0]
	}

	return Manifest{
		Name:   name,
		Fn:     makeImpl(schedulable),
		ID:     xid.New(),
		Policy: policy,
	}
}
