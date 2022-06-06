package scheduler

import "github.com/Raphy42/weekend/core/errors"

var (
	ENoMoreRetry = errors.PersistentCode(errors.DResource, errors.ATooBig)
)
