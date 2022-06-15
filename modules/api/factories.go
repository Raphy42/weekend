package api

import (
	"context"
	"net/http"
	"time"

	"github.com/Raphy42/weekend/core/app"
	"github.com/Raphy42/weekend/core/config"
)

var (
	ConfHttpServerAddr = config.Key("server", "listen")
)

func ginEngineFactory(ctx context.Context, conf *config.Config, engine *app.EngineBuilder) (*Server, error) {
	addr, err := conf.String(ctx, ConfHttpServerAddr, ":8080")
	if err != nil {
		return nil, err
	}
	server := http.Server{
		Addr:         addr,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	return newServer(&server, engine), nil
}
