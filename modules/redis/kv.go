package redis

import (
	"context"
	"strings"
	"time"

	"github.com/go-redis/redis/v9"
)

type Client struct {
	//sync.RWMutex
	inner *redis.Client
	cfg   *Configuration
}

func (c *Client) SetBytes(ctx context.Context, key string, value []byte) error {
	return c.inner.Set(ctx, key, value, 0).Err()
}

func (c *Client) SetString(ctx context.Context, key, value string) error {
	return c.SetBytes(ctx, key, []byte(value))
}

func (c *Client) SetAny(ctx context.Context, key string, value any) error {
	//TODO implement me
	panic("implement me")
}

func (c *Client) SetTTL(ctx context.Context, key string, ttl time.Duration) error {
	return c.inner.PExpire(ctx, key, ttl).Err()
}

func (c *Client) GetBytes(ctx context.Context, key string) ([]byte, error) {
	result, err := c.GetString(ctx, key)
	if err != nil {
		return nil, err
	}
	return []byte(result), nil
}

func (c *Client) GetString(ctx context.Context, key string) (string, error) {
	result, err := c.inner.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}
	return result, nil
}

func (c *Client) GetAny(ctx context.Context, key string, valuePtr any) error {
	//TODO implement me
	panic("implement me")
}

func (c *Client) GetTTL(ctx context.Context, key string) (*time.Duration, error) {
	result, err := c.inner.TTL(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) Key(args ...string) string {
	return strings.Join(args, ".")
}
