package schedulable

import "context"

type Pollable interface {
	Poll(ctx context.Context) (interface{}, error)
}
