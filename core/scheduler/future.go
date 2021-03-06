package scheduler

import (
	"context"
	"time"

	"github.com/palantir/stacktrace"
	"github.com/rs/xid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"

	"github.com/Raphy42/weekend/core/errors"
	"github.com/Raphy42/weekend/core/scheduler/async"
)

type Future struct {
	*Context
	ID         xid.ID
	ManifestID xid.ID
	result     <-chan any
	error      <-chan error
}

func NewFuture(ctx context.Context, parent xid.ID, manifest async.Manifest) (*Future, chan<- any, chan<- error) {
	switch v := ctx.(type) {
	case Context:
		parent = v.Parent
	}

	result := make(chan any)
	err := make(chan error)
	return &Future{
		Context:    NewContext(ctx, parent),
		ID:         xid.New(),
		ManifestID: manifest.ID,
		result:     result,
		error:      err,
	}, result, err
}

func (h Future) Poll(ctx context.Context) (any, error) {
	ctx, span := otel.Tracer("wk.core.scheduler").Start(ctx, "poll")
	span.SetAttributes(
		attribute.Stringer("wk.handle.id", h.ID),
		attribute.Stringer("wk.parent.id", h.Parent),
		attribute.Stringer("wk.manifest.id", h.ManifestID),
	)
	defer span.End()

	select {
	case <-ctx.Done():
		return nil, stacktrace.PropagateWithCode(ctx.Err(), errors.EInvalidContext, "invalid context")
	case err := <-h.error:
		return nil, err
	case result := <-h.result:
		return result, nil
	}
}

func (h Future) TryPoll(ctx context.Context, duration time.Duration) (any, bool, error) {
	timeoutCtx, cancel := context.WithTimeout(ctx, duration)
	defer cancel()

	result, err := h.Poll(timeoutCtx)
	if errors.HasCode(err, errors.EInvalidContext) {
		return nil, false, nil
	}
	return result, true, err
}

func (h Future) Error() <-chan error {
	return h.error
}
