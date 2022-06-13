package errors

import (
	"fmt"

	"go.uber.org/zap"

	"github.com/Raphy42/weekend/core/logger"
)

func Mustf(err error, format string, args ...any) {
	if err == nil {
		return
	}

	defer func() {
		panic(err)
	}()

	diag := Diagnostic(err)
	log := logger.New(logger.SkipCallFrame(1))
	log.Debug(fmt.Sprintf(format, args...), zap.Error(err), zap.Stringer("errors.diagnostic", diag))
}

func Must(err error) {
	if err == nil {
		return
	}

	defer func() {
		panic(err)
	}()

	diag := Diagnostic(err)
	log := logger.New(logger.SkipCallFrame(1))
	log.Debug("unrecoverable error", zap.Error(err), zap.Stringer("errors.diagnostic", diag))
	panic(err)
}
