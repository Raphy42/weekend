package errors

import (
	"testing"

	"github.com/palantir/stacktrace"
	"github.com/stretchr/testify/assert"
)

var (
	resourceFooNotFound        = uint16(PersistentCode(DResource, ANotFound))
	resourceNotFoundDiagnostic = DiagnosticManifest{
		Transient: false,
		Domain:    "resource",
		Axiom:     "not_found",
	}
)

func TestDiagnostic(t *testing.T) {
	a := assert.New(t)
	err := stacktrace.NewErrorWithCode(stacktrace.ErrorCode(resourceFooNotFound), "bar with id 'toto' could not be found")
	diag := Diagnostic(err)
	a.Equal(resourceNotFoundDiagnostic.Domain, diag.Domain)
	a.Equal(resourceNotFoundDiagnostic.Transient, diag.Transient)
	a.Equal(resourceNotFoundDiagnostic.Axiom, diag.Axiom)
}
