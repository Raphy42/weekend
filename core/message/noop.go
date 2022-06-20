package message

import (
	"context"
)

type NoOpBus struct{}

func NewNoopBus() *NoOpBus {
	return &NoOpBus{}
}

func (b NoOpBus) Subject(_ context.Context, _ string) (Mailbox, error) {
	return b, nil
}

func (NoOpBus) Emit(_ context.Context, _ Message) error {
	return nil
}

func (NoOpBus) ReadC(_ context.Context) (<-chan Message, context.CancelFunc) {
	out := make(chan Message)
	return out, func() {}
}
