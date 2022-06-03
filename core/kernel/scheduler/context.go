package scheduler

import (
	"context"

	"github.com/google/uuid"
)

type Context struct {
	context.Context
	Cancel    context.CancelFunc
	Scheduler *Scheduler
	Parent    uuid.UUID
}

func NewContext(parent context.Context, parentID uuid.UUID) *Context {
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
