package redis

import (
	"context"

	"github.com/go-redis/redis/v9"
)

type Client struct {
	//sync.RWMutex
	inner *redis.Client
	cfg   *Configuration
}

func (c *Client) Ping(ctx context.Context) error {
	return c.inner.Ping(ctx).Err()
}
