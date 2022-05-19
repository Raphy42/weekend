package errors

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	persistentTestError int16 = 0x7000
	transientTestError  int16 = 0x0000

	serviceFooUnreachable = int16(Code(KTransient, DService, AUnreachable))
	resourceBarNotFound   = int16(Code(KPersistent, DResource, ANotFound))
)

func TestIsPersistent(t *testing.T) {
	a := assert.New(t)
	a.False(IsPersistentCode(transientTestError))
	a.True(IsPersistentCode(persistentTestError))
	a.True(IsPersistentCode(resourceBarNotFound))
}

func TestIsTransient(t *testing.T) {
	a := assert.New(t)
	a.False(IsTransientCode(persistentTestError))
	a.True(IsTransientCode(transientTestError))
	a.False(IsTransientCode(resourceBarNotFound))
	a.True(IsTransientCode(serviceFooUnreachable))
}

func TestIsNotFound(t *testing.T) {
	a := assert.New(t)
	a.True(IsNotFoundCode(resourceBarNotFound))
	a.False(IsNotFoundCode(transientTestError))
}
