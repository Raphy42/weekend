package job

import (
	"context"
	"time"

	"github.com/Raphy42/weekend/core/app"
	"github.com/Raphy42/weekend/core/task"
	"github.com/Raphy42/weekend/modules/api"
	"github.com/Raphy42/weekend/modules/nats"
	"github.com/Raphy42/weekend/modules/redis"
)

func apiEndpointFactory(controller *task.Controller, server *api.Server) error {
	a := newApi(controller)

	taskGroup := server.Group("/api/v1/tasks")
	taskGroup.POST("", api.MakeJSONHandler(a.CreateTask))
	taskGroup.DELETE("/:task_id", api.MakeJSONHandler(a.CancelTask))

	workerGroup := server.Group("/api/v1/workers")
	workerGroup.GET("", api.MakeJSONHandler(a.ListWorkers))

	return nil
}

func newWorkerFactory(options Options) func(builder *app.EngineBuilder, nats *nats.Client) (*task.Worker, error) {
	return func(builder *app.EngineBuilder, nats *nats.Client) (*task.Worker, error) {
		worker := task.NewWorker(nats, options.Tasks...)
		builder.HealthCheck(worker, time.Second*1, func(ctx context.Context) error {
			return worker.Announce(ctx)
		})

		return worker, nil
	}
}

func controllerFactory(nats *nats.Client, redis *redis.Client, builder *app.EngineBuilder) *task.Controller {
	controller := task.NewController(nats, redis)
	builder.Background(controller.Producer())
	return controller
}
