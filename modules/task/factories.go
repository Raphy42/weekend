package task

import "github.com/Raphy42/weekend/modules/api"

func apiEndpointFactory(client *Client, server *api.Server) error {
	a := newApi(client)

	group := server.Group("/v1/tasks")
	group.POST("", api.MakeJSONHandler(a.CreateTask))
	group.DELETE("/:task_id", api.MakeJSONHandler(a.CreateTask))

	return nil
}
