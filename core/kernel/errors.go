package kernel

import (
	"github.com/palantir/stacktrace"

	"github.com/Raphy42/weekend/core/errors"
)

var (
	ETooManyRetries = errors.PersistentCode(errors.DSynchro, errors.ATooBig)
)

func IsTooManyRetries(current, max int) error {
	if current > max {
		return stacktrace.NewErrorWithCode(ETooManyRetries, "maximum number of retry reached (%d/%d)", current, max)
	}
	return nil
}
