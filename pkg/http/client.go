package http

import (
	"context"
)

type Client interface {
	RequestUse(middlewares ...ClientRequestMiddleware)
	ResponseUse(middlewares ...ClientResponseMiddleware)
	Execute(ctx context.Context, request Request) (*Response, error)
}
