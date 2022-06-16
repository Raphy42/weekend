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
	"github.com/Raphy42/weekend/core/service"
)

func clientFactory(
	ctx context.Context,
	cfg *config.Config,
	builder *app.EngineBuilder,
	health *service.Registry,
) (*Client, error) {
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

	client := &Client{
		inner: inner,
		cfg:   conf,
	}
	builder.HealthCheck(client, time.Second*5,
		func(ctx context.Context) error {
			if err := client.Ping(ctx); err != nil {
				health.Set(client, err)
				return err
			}
			health.Set(client, nil)
			return nil
		},
	)

	return client, nil
}
