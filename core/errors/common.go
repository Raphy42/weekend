package errors

import (
	"context"

	"github.com/palantir/stacktrace"
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
	// ENil todo
	ENil = PersistentCode(DType, ANil)
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

func InvalidCast(expected, got any) error {
	return stacktrace.NewErrorWithCode(EInvalidCast, "could not cast any to %T, was %T", expected, got)
}
