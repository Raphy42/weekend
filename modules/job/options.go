package job

import (
	"github.com/Raphy42/weekend/core/task"
	"github.com/Raphy42/weekend/internal/tasks/roundtrip"
	"github.com/Raphy42/weekend/pkg/slice"
)

type Options struct {
	Tasks []task.Task
}

func defaultOptions() Options {
	return Options{
		Tasks: slice.New(
			roundtrip.Task,
		),
	}
}

type Option func(options *Options)

func (o *Options) apply(opts ...Option) {
	for _, opt := range opts {
		opt(o)
	}
}

func WithTasks(tasks ...task.Task) Option {
	return func(options *Options) {
		options.Tasks = append(options.Tasks, tasks...)
	}
}
