package roundtrip

import (
	"context"

	"github.com/Raphy42/weekend/core/dep"
	"github.com/Raphy42/weekend/core/task"
)

var Name = dep.Name("wk", "internal", "round_trip")

func New(args any) (*task.Manifest, error) {
	return task.MsgPack(Name, args)
}

func RoundTrip(_ context.Context, in any) (any, error) {
	return in, nil
}

var Task = task.Of(
	Name,
	RoundTrip,
)
