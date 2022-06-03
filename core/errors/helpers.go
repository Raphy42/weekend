package errors

import (
	"context"

	"github.com/palantir/stacktrace"

	"github.com/Raphy42/weekend/core/bitmask"
	"github.com/Raphy42/weekend/core/reflect"
)

//All executes a list of `func() error` and stops at the first non-nil error encountered
func All(fns ...func() error) error {
	for idx, fn := range fns {
		if err := fn(); err != nil {
			return stacktrace.Propagate(err, "function %s, index %d, returned non-nil error", reflect.Typename(fn), idx)
		}
	}
	return nil
}

//AllCtx executes a list of `func(context.Context) error` and stops if the `context.Context` is no longer valid,
// or at the first non-nil error encountered
func AllCtx(ctx context.Context, fns ...func(ctx context.Context) error) error {
	if err := ctx.Err(); err != nil {
		return stacktrace.PropagateWithCode(err, EInvalidContext, "invalid context")
	}
	for idx, fn := range fns {
		if err := fn(ctx); err != nil {
			return stacktrace.Propagate(err, "function %s, index %d, returned non-nil error", reflect.Typename(fn), idx)
		}
	}
	return nil
}

func HasFlag(err error, flag int16) bool {
	if v := stacktrace.GetCode(err); v != stacktrace.NoCode {
		return bitmask.Has(int16(v), flag)
	}
	return false
}

func T() error {
	return stacktrace.NewError("DUMMY ERROR MY ONLY PURPOSE IS TYPING DO NOT USE ME")
}
