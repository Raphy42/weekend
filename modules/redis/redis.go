package redis

import (
	"context"
)

func (c *Client) Ping(ctx context.Context) error {
	return c.inner.Ping(ctx).Err()
}
