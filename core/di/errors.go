package di

import "github.com/Raphy42/weekend/core/errors"

var (
	//EDependencyMissing is an error code signaling missing providers in the di DAG
	EDependencyMissing = errors.PersistentCode(errors.DDependency, errors.ANotFound)
	//EDependencyInitFailed is an error code signaling that a provider could not be initialised
	EDependencyInitFailed = errors.PersistentCode(errors.DDependency, errors.AUnexpected)
	//EInvocationFailed is an error code signaling that an invocation returned a non-nil error
	EInvocationFailed = errors.PersistentCode(errors.DDependency, errors.AUnexpected)
)
