package kernel

import (
	"context"
)

type Schedulable interface {
	Options() Options
	Execute(ctx context.Context, args interface{}) (interface{}, error)
}

type TaskSchedulable struct {
	options Options
	task    func(ctx context.Context, args interface{}) (interface{}, error)
}

func (t TaskSchedulable) Options() Options {
	return t.options
}

func (t TaskSchedulable) Execute(ctx context.Context, args interface{}) (interface{}, error) {
	return t.task(ctx, args)
}

func Task(name string, fn func(ctx context.Context, args interface{}) (interface{}, error), opts ...Option) *TaskSchedulable {
	options := DefaultOptions()
	options.Apply(opts...)
	options.Name = func() string {
		return name
	}
	return &TaskSchedulable{
		options: *options,
		task:    fn,
	}
}
