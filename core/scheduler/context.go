package scheduler

import (
	"context"

	"github.com/palantir/stacktrace"
	"github.com/rs/xid"
	"go.uber.org/zap"

	"github.com/Raphy42/weekend/core/logger"
	"github.com/Raphy42/weekend/core/message"
	"github.com/Raphy42/weekend/pkg/runtime"
)

const (
	schedulerInjectionKey = "wk.context.scheduler"
)

type Context struct {
	context.Context
	Cancel    context.CancelFunc
	Scheduler *Scheduler
	Parent    xid.ID
}

func NewContext(parent context.Context, parentID xid.ID) *Context {
	ctx, cancel := context.WithCancel(parent)

	return &Context{
		Context: ctx,
		Cancel: func() {
			log := logger.FromContext(parent)
			frame := runtime.Frame(1)

			log.Warn("scheduling context cancelled",
				zap.String("caller", frame.Caller),
				zap.Int("line", frame.Line),
				zap.String("filename", frame.Filename),
			)

			cancel()
		},
		Parent: parentID,
	}
}

func (c *Context) BindScheduler(scheduler *Scheduler) {
	c.Scheduler = scheduler
	c.Context = context.WithValue(c.Context, schedulerInjectionKey, scheduler)
}

func ParentID(ctx context.Context) xid.ID {
	switch v := ctx.(type) {
	case *Context:
		return ParentID(*v)
	case Context:
		return v.Parent
	default:
		return xid.NilID()
	}
}

func busFromContext(ctx context.Context) (message.Bus, error) {
	switch v := ctx.(type) {
	case *Context:
		return busFromContext(*v)
	case Context:
		if v.Scheduler == nil {
			return nil, stacktrace.NewError("no scheduler bound to this context")
		}
		return v.Scheduler.bus, nil
	default:
		return nil, stacktrace.NewError("not a scheduling context")
	}
}
