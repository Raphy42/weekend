package nats

import (
	"github.com/nats-io/nats.go"
	"github.com/palantir/stacktrace"
	"github.com/rs/xid"

	"github.com/Raphy42/weekend/core"
	"github.com/Raphy42/weekend/core/dep"
)

func newNats(addr string) (*Client, error) {
	name := dep.Name(core.Name(), xid.New().String())
	conn, err := nats.Connect(addr,
		nats.LameDuckModeHandler(func(conn *nats.Conn) {
			// todo handle nats server graceful shutdown
			// https://docs.nats.io/running-a-nats-service/nats_admin/lame_duck_mode
		}),
		nats.Compression(true),
		nats.ErrorHandler(func(conn *nats.Conn, subscription *nats.Subscription, err error) {
			// todo handle nats connection errors
		}),
		nats.Name(name),
	)
	if err != nil {
		return nil, stacktrace.Propagate(err, "could not instantiate nats client")
	}

	js, err := conn.JetStream(
		nats.PublishAsyncMaxPending(256),
		nats.PublishAsyncErrHandler(func(stream nats.JetStream, msg *nats.Msg, err error) {
			// todo jetstream error handling
		}),
	)

	return &Client{
		name: name,
		conn: conn,
		js:   js,
	}, nil
}
