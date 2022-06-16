package redis

import (
	"context"

	"go.uber.org/zap"

	"github.com/Raphy42/weekend/core/logger"
)

func redisVersion(ctx context.Context, client *Client) error {
	log := logger.FromContext(ctx)
	log.Debug("redis infos",
		zap.String("redis.host", client.inner.Options().Addr),
		zap.String("redis.username", client.inner.Options().Username),
		zap.Int("redis.db", client.inner.Options().DB),
	)
	return nil
}
