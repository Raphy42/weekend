package main

import (
	"context"
	"net"

	"google.golang.org/grpc"

	"github.com/Raphy42/weekend/grpc/api"
)

type TestServer struct {
	running map[string]interface{}
}

func (t TestServer) Schedule(ctx context.Context, request *api.ScheduleRequest) (*api.Task, error) {
	//TODO implement me
	panic("implement me")
}

func (t TestServer) Cancel(ctx context.Context, request *api.IdRequest) (*api.Task, error) {
	//TODO implement me
	panic("implement me")
}

func (t TestServer) Info(ctx context.Context, request *api.IdRequest) (*api.Task, error) {
	//TODO implement me
	panic("implement me")
}

func main() {
	listener, err := net.Listen("tcp", "localhost:25000")
	if err != nil {
		panic(err)
	}

	server := grpc.NewServer()
	api.RegisterScheduleServiceServer(server, &TestServer{})

	if err := server.Serve(listener); err != nil {
		panic(err)
	}
}
