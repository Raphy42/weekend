package message

import (
	"context"
)

type NoOpBus struct{}

func NewNoopBus() *NoOpBus {
	return &NoOpBus{}
}

func (NoOpBus) Emit(_ context.Context, _ Message) error {
	return nil
}

func (NoOpBus) Read(_ context.Context) (<-chan Message, context.CancelFunc) {
	out := make(chan Message)
	return out, func() {}
}
