package service

import "google.golang.org/grpc"

type Service interface {
	Register(services *Services) error
}

type Services struct {
	*grpc.Server
}

func NewServices(server *grpc.Server) *Services {
	return &Services{server}
}
