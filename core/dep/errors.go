package dep

import "github.com/Raphy42/weekend/core/errors"

var (
	// EHealthCheckTimeout todo
	EHealthCheckTimeout = errors.PersistentCode(errors.DDependency, errors.ATimeout)
)
