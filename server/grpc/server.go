package grpc

import (
	"context"
	"net"
	"strconv"
	"strings"

	"github.com/pastelnetwork/go-commons/errors"
	"github.com/pastelnetwork/supernode/server"
	"github.com/pastelnetwork/supernode/server/grpc/log"
	"github.com/pastelnetwork/supernode/server/grpc/middleware"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
)

type service interface {
	Desc() *grpc.ServiceDesc
}

// Server represents supernode server
type Server struct {
	config   *server.Config
	services []service
}

// Run starts the server
func (server *Server) Run(ctx context.Context) error {
	group, ctx := errgroup.WithContext(ctx)

	addresses := strings.Split(server.config.ListenAddresses, ",")
	grpcServer := server.grpcServer()

	for _, address := range addresses {
		address = net.JoinHostPort(strings.TrimSpace(address), strconv.Itoa(server.config.Port))

		group.Go(func() (err error) {
			defer errors.Recover(func(rec error) { err = rec })
			return server.listen(ctx, address, grpcServer)
		})
	}

	return group.Wait()
}

func (server *Server) listen(ctx context.Context, address string, grpcServer *grpc.Server) (err error) {
	listen, err := net.Listen("tcp", address)
	if err != nil {
		return errors.Errorf("failed to listen: %v", err).WithField("address", address)
	}

	errCh := make(chan error, 1)
	go func() {
		defer errors.Recover(func(rec error) { err = rec })
		log.Infof("Server listening on %q", address)
		if err := grpcServer.Serve(listen); err != nil {
			errCh <- errors.Errorf("failed to serve: %v", err).WithField("address", address)
		}
	}()

	select {
	case <-ctx.Done():
		log.Infof("Shutting down server at %q", address)
		grpcServer.Stop()
	case err := <-errCh:
		return err
	}

	return nil
}

func (server *Server) grpcServer() *grpc.Server {
	grpcServer := grpc.NewServer(
		middleware.UnaryInterceptor(),
		middleware.StreamInterceptor(),
	)

	for _, service := range server.services {
		log.Debugf("Register service %q", service.Desc().ServiceName)
		grpcServer.RegisterService(service.Desc(), service)
	}

	return grpcServer
}

// NewServer returns a new Server instance.
func NewServer(config *server.Config, services ...service) *Server {
	return &Server{
		config:   config,
		services: services,
	}
}
