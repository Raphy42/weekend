package scheduler

import (
	"context"

	"github.com/rs/xid"

	"github.com/Raphy42/weekend/pkg/channel"
)

type CoalescedResult struct {
	Origin xid.ID
	Value  interface{}
	Error  error
}

func Coalesce(ctx context.Context, handles ...*Handle) chan CoalescedResult {
	results := make(chan CoalescedResult, len(handles))
	for _, handle := range handles {
		handle := handle
		go func() {
			result, err := handle.Poll(ctx)
			_ = channel.Send(ctx, CoalescedResult{
				Origin: handle.ID,
				Value:  result,
				Error:  err,
			}, results)
		}()
	}
	return results
}
