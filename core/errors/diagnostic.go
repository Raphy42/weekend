package errors

import (
	"fmt"
	"math"

	"github.com/palantir/stacktrace"
	"go.uber.org/zap"

	"github.com/Raphy42/weekend/core/logger"
)

type DiagnosticManifest struct {
	Transient bool
	Domain    string
	Axiom     string
	Error     error
}

func (d DiagnosticManifest) String() string {
	kind := "persistent"
	if d.Transient {
		kind = "transient"
	}
	return fmt.Sprintf("%s.%s.%s", kind, d.Domain, d.Axiom)
}

func noDiagnostic(err error) DiagnosticManifest {
	return DiagnosticManifest{
		Transient: false,
		Domain:    "unknown",
		Axiom:     "unknown",
		Error:     err,
	}
}

func Diagnostic(err error) DiagnosticManifest {
	code := int16(stacktrace.GetCode(err))
	if code == math.MaxInt16 {
		return noDiagnostic(err)
	}

	transient := IsTransientCode(code)
	domainCode := code & 0x0f00
	axiomCode := code & 0x00f0

	domain := globalRegistry.Domain(domainCode)
	axiom := globalRegistry.Axiom(axiomCode)

	return DiagnosticManifest{
		Transient: transient,
		Domain:    domain,
		Axiom:     axiom,
		Error:     err,
	}
}

func Mustf(err error, format string, args ...interface{}) {
	if err == nil {
		return
	}
	diag := Diagnostic(err)
	log := logger.New(logger.SkipCallFrame(1))
	log.Fatal(fmt.Sprintf(format, args...), zap.Error(err), zap.Stringer("errors.diagnostic", diag))
}

func Must(err error) {
	if err == nil {
		return
	}
	diag := Diagnostic(err)
	log := logger.New(logger.SkipCallFrame(1))
	log.Fatal("caught unrecoverable error", zap.Error(err), zap.Stringer("errors.diagnostic", diag))
}
