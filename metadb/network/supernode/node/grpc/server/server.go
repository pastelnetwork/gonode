package server

import (
	"context"
	"net"
	"strconv"
	"strings"

	"github.com/pastelnetwork/gonode/common/errgroup"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/metadb/network/supernode/node/grpc/server/middleware"
	"google.golang.org/grpc"
)

const (
	logPrefix = "server"
)

type service interface {
	Desc() *grpc.ServiceDesc
}

// Server represents supernode server
type Server struct {
	config   *Config
	services []service
}

// Run starts the server
func (server *Server) Run(ctx context.Context) error {
	ctx = log.ContextWithPrefix(ctx, logPrefix)

	group, ctx := errgroup.WithContext(ctx)

	addresses := strings.Split(server.config.ListenAddresses, ",")
	grpcServer := server.grpcServer(ctx)

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
		return errors.Errorf("failed to listen: %w", err).WithField("address", address)
	}

	errCh := make(chan error, 1)
	go func() {
		defer errors.Recover(func(recErr error) { err = recErr })
		log.WithContext(ctx).Infof("gRPC server listening on %q", address)
		if err := grpcServer.Serve(listen); err != nil {
			errCh <- errors.Errorf("failed to serve: %w", err).WithField("address", address)
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
func New(config *Config, services ...service) *Server {
	return &Server{
		config:   config,
		services: services,
	}
}
