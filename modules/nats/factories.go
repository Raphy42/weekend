package nats

import (
	"context"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/palantir/stacktrace"

	"github.com/Raphy42/weekend/core/app"
	"github.com/Raphy42/weekend/core/config"
)

var (
	ConfNatsDsn = config.Key("nats", "dsn")
)

func natsFactory(ctx context.Context, config *config.Config, builder *app.EngineBuilder) (*Client, error) {
	dsn, err := config.String(ctx, ConfNatsDsn, nats.DefaultURL)
	if err != nil {
		return nil, stacktrace.Propagate(err, "no entry found in config for NATS url")
	}

	n, err := newNats(dsn)
	if err != nil {
		return nil, err
	}

	builder.HealthCheck(n, 15*time.Second, func(ctx context.Context) error {
		_, err := n.conn.RTT()
		return err
	})

	return n, nil
}
