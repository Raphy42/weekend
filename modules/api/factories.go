package api

import (
	"context"
	"net/http"
	"time"

	"github.com/Raphy42/weekend/core/app"
	"github.com/Raphy42/weekend/core/config"
	"github.com/Raphy42/weekend/core/scheduler/async"
)

var (
	ConfHttpServerAddr = config.Key("server", "listen")
)

func ginEngineFactory(ctx context.Context, conf *config.Config, builder *app.EngineBuilder) (*Server, error) {
	addr, err := conf.String(ctx, ConfHttpServerAddr, ":8080")
	if err != nil {
		return nil, err
	}
	server := http.Server{
		Addr:         addr,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	s := newServer(&server)
	builder.Background(async.Of(
		async.Name("wk.api.serve"),
		s.listenAndServe,
	))
	builder.HealthCheck(s, time.Second*10, func(_ context.Context) error {
		return nil
	})

	return s, nil
}
