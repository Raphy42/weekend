package scheduler

import (
	"context"

	"weekend/grpc/api"
	"weekend/grpc/service"
)

const (
	SERVICE_NAME = "scheduler.service"
)

type SchedulerService struct {
}

func (s *SchedulerService) Register(services *service.Services) {
	api.RegisterScheduleServiceServer(services.Server, s)
}

func (s SchedulerService) Schedule(ctx context.Context, request *api.ScheduleRequest) (*api.Task, error) {
	//TODO implement me
	panic("implement me")
}

func (s SchedulerService) Cancel(ctx context.Context, request *api.IdRequest) (*api.Task, error) {
	//TODO implement me
	panic("implement me")
}

func (s SchedulerService) Info(ctx context.Context, request *api.IdRequest) (*api.Task, error) {
	//TODO implement me
	panic("implement me")
}
