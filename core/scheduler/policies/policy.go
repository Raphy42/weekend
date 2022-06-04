package policies

import (
	"github.com/palantir/stacktrace"

	"github.com/Raphy42/weekend/core/errors"
)

type PolicyAction struct {
	Restart bool
}

type Policy struct {
	RetryPolicy       RetryPolicy
	Retry             int
	ConcurrencyPolicy ConcurrencyPolicy
	Concurrency       int
}

func (p Policy) ShouldRetry(current int) bool {
	switch p.RetryPolicy {
	case RetryAlways:
		return true
	case RetryConstant:
		return p.Retry <= current
	case RetryNever:
		return false
	default:
		panic(stacktrace.NewMessageWithCode(errors.EUnreachable, "no such retry policy: '%d'", p.RetryPolicy))
	}
}

func (p Policy) Scale(current int) int {
	switch p.ConcurrencyPolicy {
	case ConcurrencyAuto:
		return 1
	case ConcurrencyConstant:
		if p.Concurrency == -1 {
			return 1
		}
		return p.Concurrency - current
	default:
		panic(stacktrace.NewMessageWithCode(errors.EUnreachable, "no such concurrency policy: '%d'", p.ConcurrencyPolicy))
	}
}

func Default() Policy {
	return Policy{
		RetryPolicy:       RetryConstant,
		Retry:             5,
		ConcurrencyPolicy: ConcurrencyConstant,
		Concurrency:       16,
	}
}
