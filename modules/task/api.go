package task

import (
	"context"

	"github.com/rs/xid"
)

type Api struct {
	client *Client
}

func newApi(client *Client) *Api {
	return &Api{client: client}
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
