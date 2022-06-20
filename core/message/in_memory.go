package message

import (
	"context"

	"github.com/Raphy42/weekend/pkg/channel"
)

const (
	InMemoryBusMaximumInFlightMessage = 4096
)

type InMemoryBus struct{}

func NewInMemoryBus() *InMemoryBus {
	return &InMemoryBus{}
}

func (i *InMemoryBus) Subject(ctx context.Context, subject string) (Mailbox, error) {
	return NewInMemoryMailbox(), nil
}

type InMemoryMailbox struct {
	messages chan Message
}

func NewInMemoryMailbox() *InMemoryMailbox {
	return &InMemoryMailbox{messages: make(chan Message, InMemoryBusMaximumInFlightMessage)}
}

func (i *InMemoryMailbox) Emit(ctx context.Context, message Message) error {
	return channel.Send(ctx, message, i.messages)
}

func (i *InMemoryMailbox) ReadC(ctx context.Context) (<-chan Message, context.CancelFunc) {
	out := make(chan Message)
	cancel := channel.Multicast(ctx, i.messages, out)
	return out, cancel
}
