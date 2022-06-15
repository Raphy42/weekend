package redis

import (
	"context"
	"time"

	"go.uber.org/zap"

	"github.com/Raphy42/weekend/core/app"
	"github.com/Raphy42/weekend/core/logger"
	"github.com/Raphy42/weekend/core/service"
)

func redisVersion(ctx context.Context, client *Client) error {
	log := logger.FromContext(ctx)
	log.Debug("connecting to redis",
		zap.String("redis.host", client.inner.Options().Addr),
		zap.String("redis.username", client.inner.Options().Username),
	)
	return nil
}

func redisHealthCheck(client *Client, builder *app.EngineBuilder, health *service.Registry) error {
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
	return nil
}
