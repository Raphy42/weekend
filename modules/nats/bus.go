package nats

import (
	"context"

	"github.com/nats-io/nats.go"
	"go.uber.org/zap"

	"github.com/Raphy42/weekend/core/logger"
	"github.com/Raphy42/weekend/core/message"
)

type Client struct {
	name string
	conn *nats.Conn
	js   nats.JetStreamContext
}

func (c *Client) Subject(ctx context.Context, subject string) (message.Mailbox, error) {
	log := logger.FromContext(ctx)
	log.Info("subject created", zap.String("wk.message.subject", subject))

	return &Mailbox{
		client:  c,
		subject: subject,
	}, nil
}
