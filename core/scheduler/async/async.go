package async

import (
	"context"
	"reflect"

	"github.com/palantir/stacktrace"
	"github.com/rs/xid"

	"github.com/Raphy42/weekend/core/errors"
)

type Async interface {
	func() |
		func() error |
		func(ctx context.Context) |
		func(ctx context.Context) error |
		func(ctx context.Context) (any, error) |
		func(ctx context.Context, args any) |
		func(ctx context.Context, args any) error |
		func(ctx context.Context, args any) (any, error)
}

type Fn func(ctx context.Context, args any) (any, error)
type Manifest struct {
	Name string
	ID   xid.ID
	Fn   Fn
}

// oooooooh this is dirty
// too bad golang generics don't allow specialisation :(
// but it's less dirty than ye big ol'trusty `any`
func makeImpl[S Async](schedulable S) Fn {
	v := reflect.ValueOf(schedulable).Interface()
	switch fn := v.(type) {
	case func():
		return func(_ context.Context, args any) (any, error) {
			fn()
			return args, nil
		}
	case func() error:
		return func(_ context.Context, args any) (any, error) {
			return args, fn()
		}
	case func(ctx context.Context):
		return func(ctx context.Context, args any) (any, error) {
			fn(ctx)
			return args, nil
		}
	case func(ctx context.Context) error:
		return func(ctx context.Context, args any) (any, error) {
			return args, fn(ctx)
		}
	case func(ctx context.Context) (any, error):
		return func(ctx context.Context, _ any) (any, error) {
			return fn(ctx)
		}
	case func(ctx context.Context, args any) error:
		return func(ctx context.Context, args any) (any, error) {
			return args, fn(ctx, args)
		}
	case func(ctx context.Context, args any) (any, error):
		return fn
	default:
		panic(
			stacktrace.NewMessageWithCode(errors.EUnreachable, "invalid Async '%T'", fn),
		)
	}
}

func Of[S Async](name string, schedulable S) Manifest {
	return Manifest{
		Name: name,
		Fn:   makeImpl(schedulable),
		ID:   xid.New(),
	}
}
