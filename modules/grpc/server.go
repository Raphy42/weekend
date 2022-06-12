//go:build grpc

package grpc

import (
	"context"
	"net"

	"github.com/palantir/stacktrace"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/Raphy42/weekend/core/logger"
)

type Service interface {
	Name() string
	Register(server *grpc.Server)
}

type Server struct {
	listener net.Listener
	opts     []grpc.ServerOption
	inner    *grpc.Server
	services []Service
}

func (s *Server) apply(opts ...ServerOption) error {
	for _, opt := range opts {
		if err := opt(s); err != nil {
			return stacktrace.Propagate(err, "could not build grpc.Server")
		}
	}
	return nil
}

type ServerOption func(server *Server) error

func WithTLS(cert, key string) ServerOption {
	return func(server *Server) error {
		creds, err := credentials.NewServerTLSFromFile(cert, key)
		if err != nil {
			return stacktrace.Propagate(err, "could not build TLS credentials")
		}
		server.opts = append(server.opts, grpc.Creds(creds))
		return nil
	}
}

func WithServices(services ...Service) ServerOption {
	return func(server *Server) error {
		server.services = append(server.services, services...)
		return nil
	}
}

func NewServer(listener net.Listener, opts ...ServerOption) (*Server, error) {
	s := Server{
		opts:     make([]grpc.ServerOption, 0),
		services: make([]Service, 0),
		listener: listener,
	}
	if err := s.apply(opts...); err != nil {
		return nil, err
	}
	return &s, nil
}

func (s *Server) initialise(ctx context.Context) {
	log := logger.FromContext(ctx)

	for _, service := range s.services {
		service.Register(s.inner)
		log.Debug("registered service", zap.String("kw.grpc.service", service.Name()))
	}
}

func (s *Server) start(ctx context.Context) error {
	log := logger.FromContext(ctx)

	log.Info("grpc server listening",
		zap.Stringer("listener.addr", s.listener.Addr()),
	)
	if err := s.inner.Serve(s.listener); err != nil {
		return stacktrace.Propagate(err, "server.Serve returned non-nil error")
	}
	return nil
}
