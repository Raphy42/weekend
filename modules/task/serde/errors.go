package serde

import "github.com/Raphy42/weekend/core/errors"

var (
	EEncodingNotFound = errors.PersistentCode(errors.DEncoding, errors.ANotFound)
)
