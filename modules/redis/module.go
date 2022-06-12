package redis

import (
	"context"

	"github.com/go-redis/redis/extra/redisotel/v9"
	"github.com/go-redis/redis/v9"
	"github.com/palantir/stacktrace"
	"go.uber.org/zap"

	"github.com/Raphy42/weekend/core/app"
	"github.com/Raphy42/weekend/core/config"
	"github.com/Raphy42/weekend/core/dep"
	"github.com/Raphy42/weekend/core/errors"
	"github.com/Raphy42/weekend/core/logger"
	"github.com/Raphy42/weekend/pkg/chrono"
)

var (
	ModuleName = dep.Name("wk", "redis")
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
			Addr:     server.Addr,
			Username: server.Username,
			Password: server.Password,
			DB:       server.Database,
		})
	case ModeSentinel:
		panic("todo")
	default:
		return nil, stacktrace.NewErrorWithCode(errors.EUnreachable, "unexpected redis mode: '%s'", conf.Mode)
	}

	inner.AddHook(redisotel.NewTracingHook())

	return &Client{
		inner: inner,
		cfg:   conf,
	}, nil
}

func redisVersion(ctx context.Context, client *Client) error {
	log := logger.FromContext(ctx)
	log.Debug("connecting to redis",
		zap.String("redis.host", client.inner.Options().Addr),
		zap.String("redis.username", client.inner.Options().Username),
	)

	timer := chrono.NewChrono()
	timer.Start()
	pong, err := client.inner.Ping(ctx).Result()
	if err != nil {
		return stacktrace.Propagate(err, "could not get redis PING result")
	}
	log.Info("connected to redis", zap.String("ping.reply", pong), zap.Duration("ping.delay", timer.Elapsed()))
	return nil
}

func Module() dep.Module {
	return dep.Declare(
		ModuleName,
		dep.Factories(
			clientFactory,
		),
		dep.SideEffects(
			redisVersion,
		),
	)
}
