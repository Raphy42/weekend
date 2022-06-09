package redis

import (
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/palantir/stacktrace"

	"github.com/Raphy42/weekend/core/config"
	"github.com/Raphy42/weekend/core/errors"
)

type Client struct {
	//sync.RWMutex
	inner *redis.Client
}

func newClient(ctx context.Context, cfg config.Config) (*Client, error) {
	conf, err := ConfigFrom(ctx, cfg)
	if err != nil {
		return nil, stacktrace.Propagate(err, "invalid redis config")
	}

	var inner *redis.Client
	switch conf.Mode {
	case ModeCluster:
	case ModeLocal:
		server := conf.Servers[0]
		inner = redis.NewClient(&redis.Options{
			Addr:     server.Addr,
			Username: server.Username,
			Password: server.Password,
			DB:       server.Database,
		})
	case ModeSentinel:
	default:
		return nil, stacktrace.NewErrorWithCode(errors.EUnreachable, "unexpected redis mode: '%s'", conf.Mode)
	}

	return &Client{
		inner: inner,
	}, nil
}
