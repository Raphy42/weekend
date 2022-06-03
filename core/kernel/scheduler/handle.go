package scheduler

import (
	"context"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/Raphy42/weekend/core/logger"
)

type Handle struct {
	*Context
	ID     uuid.UUID
	result <-chan interface{}
	error  <-chan error
}

func NewHandle(ctx context.Context, parent uuid.UUID) (*Handle, chan<- interface{}, chan<- error) {
	switch v := ctx.(type) {
	case Context:
		parent = v.Parent
	}

	result := make(chan interface{})
	err := make(chan error)
	return &Handle{
		Context: NewContext(ctx, parent),
		ID:      uuid.New(),
		result:  result,
		error:   err,
	}, result, err
}

func (h Handle) Poll(ctx context.Context) (interface{}, error) {
	log := logger.FromContext(ctx).With(zap.Stringer("wk.handle.id", h.ID))
	log.Debug("polling started")
	complete := func() {
		log.Debug("polling complete")
	}

	select {
	case <-ctx.Done():
		log.Debug("polling cancelled")
		return nil, ctx.Err()
	case err := <-h.error:
		defer complete()
		return nil, err
	case result := <-h.result:
		defer complete()
		return result, nil
	}
}
