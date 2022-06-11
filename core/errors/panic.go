package errors

import (
	"fmt"

	"go.uber.org/zap"

	"github.com/Raphy42/weekend/core/logger"
)

func Mustf(err error, format string, args ...interface{}) {
	if err == nil {
		return
	}
	diag := Diagnostic(err)
	log := logger.New(logger.SkipCallFrame(1))
	log.Debug(fmt.Sprintf(format, args...), zap.NamedError("crash", err), zap.Stringer("errors.diagnostic", diag))
	panic(err)
}

func Must(err error) {
	if err == nil {
		return
	}
	diag := Diagnostic(err)
	log := logger.New(logger.SkipCallFrame(1))
	log.Debug("unrecoverable error", zap.NamedError("crash", err), zap.Stringer("errors.diagnostic", diag))
	panic(err)
}
