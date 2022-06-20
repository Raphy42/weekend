package nats

import (
	"context"

	"github.com/nats-io/nats.go"

	"github.com/Raphy42/weekend/core/errors"
	"github.com/Raphy42/weekend/core/message"
	"github.com/Raphy42/weekend/pkg/channel"
)

//todo rewrite this and stabilise Bus->Mailbox interfaces

type Mailbox struct {
	client  *Client
	subject string
}

func (m *Mailbox) Emit(ctx context.Context, message message.Message) error {
	buf, err := message.Marshall()
	if err != nil {
		return err
	}
	return m.client.conn.Publish(m.subject, buf)
}

func (m *Mailbox) ReadC(ctx context.Context) (<-chan message.Message, context.CancelFunc) {
	messages := make(chan message.Message)
	subscription, err := m.client.conn.Subscribe(m.subject, func(in *nats.Msg) {
		msg, err := message.Unmarshall(in.Data)
		if err != nil {
			close(messages)
			errors.Mustf(err, "invalid nats message data")
		}
		errors.Must(channel.Send(ctx, *msg, messages))
	})
	errors.Must(err)

	return messages, func() {
		errors.Must(subscription.Unsubscribe())
	}
}
