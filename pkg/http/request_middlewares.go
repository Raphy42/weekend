package http

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

func InstrumentedRequestMiddleware(name string) ClientRequestMiddleware {
	return func(ctx context.Context, r *Request, next func(ctx context.Context, err error)) {
		ctx, span := otel.Tracer(name).Start(ctx, "wk.http.execute")
		span.SetAttributes(
			attribute.String("http.method", r.Method),
			attribute.String("http.path", r.Path),
			attribute.String("http.contentType", r.ContentType),
			attribute.String("wk.http.service", r.Service),
		)
		defer span.End()

		next(ctx, nil)
	}
}
