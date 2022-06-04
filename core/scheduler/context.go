package scheduler

import (
	"context"

	"github.com/palantir/stacktrace"
	"github.com/rs/xid"

	"github.com/Raphy42/weekend/core/message"
)

var (
	mainID = xid.New()
)

func IsMainProcessID(id xid.ID) bool {
	return id.Compare(mainID) == 0
}

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
		Cancel:  cancel,
		Parent:  parentID,
	}
}

func (c *Context) BindScheduler(scheduler *Scheduler) {
	c.Scheduler = scheduler
}

func ParentID(ctx context.Context) xid.ID {
	switch v := ctx.(type) {
	case *Context:
		return ParentID(*v)
	case Context:
		return v.Parent
	default:
		return mainID
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
