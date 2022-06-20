package kv

import "github.com/Raphy42/weekend/core/errors"

var (
	EEntryNotFound = errors.PersistentCode(errors.DResource, errors.ANotFound)
)
