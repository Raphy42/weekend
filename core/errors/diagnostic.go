package errors

import (
	"fmt"

	"github.com/palantir/stacktrace"
)

type DiagnosticManifest struct {
	Transient bool
	Domain    string
	Axiom     string
	Error     error
	external  bool
}

func (d DiagnosticManifest) String() string {
	if d.external {
		return "external"
	}
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
		external:  true,
	}
}

func Diagnostic(err error) DiagnosticManifest {
	c := stacktrace.GetCode(err)
	if c == stacktrace.NoCode {
		return noDiagnostic(err)
	}
	code := uint16(c)

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
