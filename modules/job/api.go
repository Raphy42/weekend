package job

import (
	"context"
	"time"

	"github.com/rs/xid"

	"github.com/Raphy42/weekend/core/task"
	"github.com/Raphy42/weekend/pkg/set"
	"github.com/Raphy42/weekend/pkg/slice"
)

type Api struct {
	controller *task.Controller
}

func newApi(controller *task.Controller) *Api {
	return &Api{controller: controller}
}

type CreateTaskRequest struct {
	Name        string      `body:"name"`
	Args        interface{} `body:"args"`
	ContentType string      `body:"contentType"`
}

type CreateTaskResponse struct {
	TaskID xid.ID `json:"taskID"`
}

func (a *Api) CreateTask(ctx context.Context, request *CreateTaskRequest) (*CreateTaskResponse, error) {
	return nil, nil
}

type CancelTaskRequest struct {
	TaskID xid.ID `path:"task_id"`
}

type CancelTaskResponse struct{}

func (a *Api) CancelTask(ctx context.Context, request *CancelTaskRequest) (*CancelTaskResponse, error) {
	return nil, nil
}

type ListWorkerRequest struct{}

type ListWorkerResponse struct {
	Workers []ListWorkerResponseItem `json:"workers"`
}

type ListWorkerResponseItem struct {
	ID              string           `json:"id"`
	LastUpdate      time.Time        `json:"lastUpdate"`
	RegisteredTasks []string         `json:"registeredTasks"`
	Load            map[string]int32 `json:"load"`
}

func (a *Api) ListWorkers(ctx context.Context, request *ListWorkerRequest) (*ListWorkerResponse, error) {
	workers := a.controller.Workers()
	return &ListWorkerResponse{
		Workers: slice.Map(workers, func(idx int, in task.WorkerInfo) ListWorkerResponseItem {
			return ListWorkerResponseItem{
				ID:              in.ID,
				LastUpdate:      in.LastUpdate,
				RegisteredTasks: set.Keys(in.Load),
				Load:            in.Load,
			}
		}),
	}, nil
}
