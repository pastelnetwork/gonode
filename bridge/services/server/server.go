package server

import (
	"context"
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/pastelnetwork/gonode/common/log"

	"github.com/pastelnetwork/gonode/common/errors"

	"github.com/pastelnetwork/gonode/bridge/services/server/middleware"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
)

type service interface {
	Desc() *grpc.ServiceDesc
}

// Server represents supernode server
type Server struct {
	config   *Config
	services []service
	name     string
}

// Run starts the server
func (server *Server) Run(ctx context.Context) error {
	ctx = log.ContextWithPrefix(ctx, server.name)

	group, ctx := errgroup.WithContext(ctx)

	addresses := strings.Split(server.config.ListenAddresses, ",")
	grpcServer := server.grpcServer(ctx)
	if grpcServer == nil {
		return fmt.Errorf("initialize grpc server failed")
	}

	for _, address := range addresses {
		addr := net.JoinHostPort(strings.TrimSpace(address), strconv.Itoa(server.config.Port))

		group.Go(func() error {
			return server.listen(ctx, addr, grpcServer)
		})
	}

	return group.Wait()
}

func (server *Server) listen(ctx context.Context, address string, grpcServer *grpc.Server) (err error) {
	listen, err := net.Listen("tcp", address)
	if err != nil {
		return errors.Errorf("listen: %w", err).WithField("address", address)
	}

	errCh := make(chan error, 1)
	go func() {
		defer errors.Recover(func(recErr error) { err = recErr })
		log.WithContext(ctx).Infof("gRPC server listening on %q", address)
		if err := grpcServer.Serve(listen); err != nil {
			errCh <- errors.Errorf("serve: %w", err).WithField("address", address)
		}
	}()

	select {
	case <-ctx.Done():
		log.WithContext(ctx).Infof("Shutting down gRPC server at %q", address)
		grpcServer.GracefulStop()
	case err := <-errCh:
		return err
	}

	return nil
}

func (server *Server) grpcServer(ctx context.Context) *grpc.Server {
	grpcServer := grpc.NewServer(
		middleware.UnaryInterceptor(),
		middleware.StreamInterceptor(),
	)

	for _, service := range server.services {
		log.WithContext(ctx).Debugf("Register service %q", service.Desc().ServiceName)
		grpcServer.RegisterService(service.Desc(), service)
	}

	return grpcServer
}

// New returns a new Server instance.
func New(config *Config, name string, services ...service) *Server {
	return &Server{
		config:   config,
		services: services,
		name:     name,
	}
}
