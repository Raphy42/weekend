package scheduler

import (
	"context"
	"fmt"

	"github.com/rs/xid"
	"go.opentelemetry.io/otel"

	"github.com/Raphy42/weekend/core/scheduler/schedulable"
)

type Handle struct {
	*Context
	ID     xid.ID
	result <-chan interface{}
	error  <-chan error
}

func NewHandle(ctx context.Context, parent xid.ID) (*Handle, chan<- interface{}, chan<- error) {
	switch v := ctx.(type) {
	case Context:
		parent = v.Parent
	}

	result := make(chan interface{})
	err := make(chan error)
	return &Handle{
		Context: NewContext(ctx, parent),
		ID:      xid.New(),
		result:  result,
		error:   err,
	}, result, err
}

func (h Handle) Poll(ctx context.Context) (interface{}, error) {
	ctx, span := otel.Tracer("Handle.Poll").Start(ctx, h.ID.String())
	defer span.End()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case err := <-h.error:
		return nil, err
	case result := <-h.result:
		return result, nil
	}
}

func (h *Handle) Manifest() schedulable.Manifest {
	return schedulable.Of(fmt.Sprintf("handle.%s", h.ID), h.Poll)
}
