package task

import (
	"context"

	"github.com/rs/xid"
)

type Client struct {
}

func (c *Client) Run(ctx context.Context, name string, args any) (xid.ID, error) {
	return xid.NilID(), nil
}

func (c *Client) RunSync(ctx context.Context, name string, args any) (any, error) {
	return nil, nil
}
