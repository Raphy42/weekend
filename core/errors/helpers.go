package errors

import (
	"context"

	"github.com/palantir/stacktrace"

	"github.com/Raphy42/weekend/pkg/bitmask"
	"github.com/Raphy42/weekend/pkg/std/slice"
)

//All executes a list of `func() error` and stops at the first non-nil error encountered
func All(fns ...func() error) error {
	for idx, fn := range fns {
		if err := fn(); err != nil {
			return stacktrace.Propagate(err, "function %T, index %d, returned non-nil error", fn, idx)
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
			return stacktrace.Propagate(err, "function %T, index %d, returned non-nil error", fn, idx)
		}
	}
	return nil
}

func HasFlag(err error, flag uint16) bool {
	if v := stacktrace.GetCode(err); v != stacktrace.NoCode {
		return bitmask.Has(uint16(v), flag)
	}
	return false
}

func HasAnyFlag(err error, flags ...uint16) bool {
	return slice.Any(flags, func(flag uint16) bool {
		return HasFlag(err, flag)
	})
}

func IsTransient(err error) bool {
	return HasFlag(err, KTransient)
}

func HasCode(err error, code stacktrace.ErrorCode) bool {
	return stacktrace.GetCode(err) == code
}
