package message

import (
	"context"

	"github.com/Raphy42/weekend/pkg/channel"
)

const (
	InMemoryBusMaximumInFlightMessage = 4096
)

type InMemoryBus struct {
	messages chan Message
}

func NewInMemoryBus() *InMemoryBus {
	return &InMemoryBus{messages: make(chan Message, InMemoryBusMaximumInFlightMessage)}
}

func (i *InMemoryBus) Emit(ctx context.Context, message Message) error {
	return channel.Send(ctx, message, i.messages)
}

func (i *InMemoryBus) Read(ctx context.Context) (<-chan Message, context.CancelFunc) {
	out := make(chan Message)
	cancel := channel.Multicast(ctx, i.messages, out)
	return out, cancel
}
