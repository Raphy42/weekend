package task

import (
	"context"

	"github.com/palantir/stacktrace"

	"github.com/Raphy42/weekend/core/errors"
	"github.com/Raphy42/weekend/core/scheduler/async"
)

type Task struct {
	Name          string
	AsyncManifest async.Manifest
}

func Of[I any, O any](name string, taskFn func(ctx context.Context, input I) (O, error)) Task {
	asyncManifest := async.Of(name, func(ctx context.Context, args any) (any, error) {
		taskManifest, ok := args.(*Manifest)
		if !ok {
			return nil, stacktrace.NewErrorWithCode(
				errors.EInvalidCast,
				"unexpected '%T' for task arguments, expected *task.Manifest", args,
			)
		}
		var input I
		if err := taskManifest.Unmarshal(&input); err != nil {
			return nil, err
		}
		return taskFn(ctx, input)
	})
	return Task{
		Name:          name,
		AsyncManifest: asyncManifest,
	}
}
