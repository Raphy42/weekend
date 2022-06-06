package scheduler

import (
	"context"
	"time"

	"github.com/palantir/stacktrace"

	"github.com/Raphy42/weekend/core/errors"
	"github.com/Raphy42/weekend/core/message"
	"github.com/Raphy42/weekend/core/scheduler/policies"
	"github.com/Raphy42/weekend/core/scheduler/schedulable"
)

type Supervisor struct {
	policy policies.Policy
}

func NewSupervisor() *Supervisor {
	return &Supervisor{
		policy: policies.Default(),
	}
}

func (s Supervisor) retry(ctx context.Context, manifest schedulable.Manifest, args interface{}, policy policies.Policy, err error) (interface{}, error) {
	code := stacktrace.GetCode(err)
	if errors.IsPersistentCode(int16(code)) {
		return nil, stacktrace.PropagateWithCode(err, ENoMoreRetry, "non transient error encountered")
	}

	timer := policy.BackoffPolicy.Backoff(time.Millisecond * 500)
	var retryBlock <-chan time.Time

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case duration := <-timer.Duration(ctx):
		retryBlock = time.After(duration)
	}

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-retryBlock:
		return s.Execute(ctx, manifest, args, &policy)
	}
}

func (s Supervisor) retryConstant(ctx context.Context, manifest schedulable.Manifest, args interface{}, policy policies.Policy, err error) (interface{}, error) {
	policy.Retry -= 1
	if policy.Retry == 0 {
		return nil, stacktrace.PropagateWithCode(err, ENoMoreRetry, "maximum retry count exceeded")
	}
	return s.retry(ctx, manifest, args, policy, err)
}

func (s Supervisor) Execute(ctx context.Context, manifest schedulable.Manifest, args interface{}, current *policies.Policy) (interface{}, error) {
	bus := message.NewInMemoryBus()
	scheduler := New(bus)
	handle, err := scheduler.Schedule(ctx, manifest, args)
	if err != nil {
		return nil, err
	}
	policy := s.policy
	if current != nil {
		policy = *current
	}

	result, err := handle.Poll(ctx)
	if err != nil {
		switch policy.RetryPolicy {
		case policies.RetryNever:
			return nil, err
		case policies.RetryAlways:
			return s.retry(ctx, manifest, args, policy, err)
		case policies.RetryConstant:
			return s.retryConstant(ctx, manifest, args, policy, err)
		default:
			return nil, stacktrace.NewErrorWithCode(errors.EUnreachable, "invalid retry policy")
		}
	}
	return result, nil
}
