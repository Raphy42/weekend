package errors

import (
	"context"

	"github.com/palantir/stacktrace"

	"github.com/Raphy42/weekend/core/reflect"
)

var (
	//EInvalidContext signals that a context is no longer valid, but should have been at the time of invocation
	EInvalidContext = PersistentCode(DSynchro, AInvariant)
	//ENotImplemented signals that this part of logic was not implemented
	ENotImplemented = PersistentCode(DLogic, AUnimplemented)
	//EInvalidCast todo
	EInvalidCast = PersistentCode(DType, AInvariant)
	//EUnreachable todo
	EUnreachable = PersistentCode(DLogic, AUnreachable)
)

func ValidateContext(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return stacktrace.PropagateWithCode(err, EInvalidContext, "expected a valid context")
	}
	return nil
}

func NotImplemented(reason string) error {
	return stacktrace.NewErrorWithCode(ENotImplemented, "reached unimplemented part of code: %s", reason)
}

func InvalidCast(expected, got interface{}) error {
	return stacktrace.NewErrorWithCode(EInvalidCast, "could not cast interface{} to %T, was %s", expected, reflect.Typename(got))
}
