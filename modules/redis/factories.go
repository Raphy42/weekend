package redis

import (
	"context"
	"time"

	"github.com/go-redis/redis/extra/redisotel/v9"
	"github.com/go-redis/redis/v9"
	"github.com/palantir/stacktrace"

	"github.com/Raphy42/weekend/core/app"
	"github.com/Raphy42/weekend/core/config"
	"github.com/Raphy42/weekend/core/errors"
)

func clientFactory(ctx context.Context, cfg *config.Config, _ *app.EngineBuilder) (*Client, error) {
	conf, err := ConfigFrom(ctx, *cfg)
	if err != nil {
		return nil, stacktrace.Propagate(err, "invalid redis config")
	}

	var inner *redis.Client
	switch conf.Mode {
	case ModeCluster:
		panic("todo")
	case ModeLocal:
		server := conf.Servers[0]
		inner = redis.NewClient(&redis.Options{
			Addr:         server.Addr,
			Username:     server.Username,
			Password:     server.Password,
			DB:           server.Database,
			DialTimeout:  500 * time.Millisecond,
			ReadTimeout:  500 * time.Millisecond,
			WriteTimeout: 2 * time.Second,
		})
	case ModeSentinel:
		panic("todo")
	default:
		return nil, stacktrace.NewErrorWithCode(errors.EUnreachable, "unexpected redis mode: '%s'", conf.Mode)
	}

	inner = inner.WithContext(ctx)
	inner.AddHook(redisotel.NewTracingHook())

	return &Client{
		inner: inner,
		cfg:   conf,
	}, nil
}
