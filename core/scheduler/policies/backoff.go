package policies

import (
	"time"

	"github.com/palantir/stacktrace"

	"github.com/Raphy42/weekend/core/errors"
	"github.com/Raphy42/weekend/pkg/chrono/backoff"
)

type BackoffPolicy int

const (
	BackoffLinear BackoffPolicy = iota
	BackoffExponential
	BackoffImmediate
)

func (b BackoffPolicy) Backoff(interval time.Duration) backoff.Backoff {
	switch b {
	case BackoffLinear:
		return backoff.NewLinear(interval)
	case BackoffExponential:
		return backoff.NewExponentialBackoff(interval)
	case BackoffImmediate:
		return backoff.NewImmediate()
	default:
		panic(stacktrace.NewErrorWithCode(errors.ENotImplemented, "no such backoff policy implementation"))
	}
}
