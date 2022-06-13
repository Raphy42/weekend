package async

import "context"

type Pollable interface {
	Poll(ctx context.Context) (any, error)
}
