package http

import "context"

type ClientRequestMiddleware func(ctx context.Context, r *Request, next MiddlewareNextFunc)
type ClientResponseMiddleware func(ctx context.Context, r *Response, next MiddlewareNextFunc)
type MiddlewareNextFunc func(ctx context.Context, err error)

type ClientMiddleware interface {
	ClientRequestMiddleware | ClientResponseMiddleware
}

func ApplyRequestMiddlewares(ctx context.Context, request *Request, middlewares ...ClientRequestMiddleware) (*Request, error) {
	for _, middleware := range middlewares {
		var err error
		middleware(ctx, request, func(middlewareCtx context.Context, middlewareErr error) {
			ctx = middlewareCtx
			err = middlewareErr
		})
		if err != nil {
			return nil, err
		}
	}
	return request, nil
}

func ApplyResponseMiddlewares(ctx context.Context, response *Response, middlewares ...ClientResponseMiddleware) (*Response, error) {
	for _, middleware := range middlewares {
		var err error
		middleware(ctx, response, func(middlewareCtx context.Context, middlewareErr error) {
			ctx = middlewareCtx
			err = middlewareErr
		})
		if err != nil {
			return nil, err
		}
	}
	return response, nil
}
