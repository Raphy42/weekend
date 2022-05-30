package kernel

import (
	"context"

	"github.com/Raphy42/weekend/core/errors"
	"github.com/Raphy42/weekend/core/reflect"
)

type Args interface {
	Decode(ctx context.Context, value interface{}) error
}

type LiteralArgs struct {
	value interface{}
}

func NewLiteralArgs(value interface{}) *LiteralArgs {
	return &LiteralArgs{value}
}

func (l LiteralArgs) Decode(ctx context.Context, value interface{}) error {
	if !reflect.SameType(l.value, value) {
		return errors.InvalidCast(l.value, value)
	}
	value = l.value
	return nil
}
