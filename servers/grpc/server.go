package grpc

import (
	"context"
	"net"
	"strconv"
	"strings"

	"github.com/pastelnetwork/go-commons/errors"
	"github.com/pastelnetwork/go-commons/log"
	"github.com/pastelnetwork/supernode/servers/grpc/services"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"

	pb "github.com/pastelnetwork/supernode-proto"
)

const logPrefix = "[grpc]"

// Server represents supernode server
type Server struct {
	config *Config
}

// Run starts the server
func (server *Server) Run(ctx context.Context) error {
	group, ctx := errgroup.WithContext(ctx)

	addresses := strings.Split(server.config.ListenAddresses, ",")
	for _, address := range addresses {
		address = net.JoinHostPort(strings.TrimSpace(address), strconv.Itoa(server.config.Port))

		group.Go(func() error {
			return server.listen(ctx, address)
		})
	}

	return group.Wait()
}

func (server *Server) listen(ctx context.Context, address string) error {
	listen, err := net.Listen("tcp", address)
	if err != nil {
		return errors.Errorf("failed to listen: %v", err).WithField("address", address)
	}

	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	pb.RegisterArtworkRegisterServer(grpcServer, services.NewArtworkRegister())

	errCh := make(chan error, 1)
	go func() {
		log.Infof("%s Server listening on %q", logPrefix, address)
		if err := grpcServer.Serve(listen); err != nil {
			errCh <- errors.Errorf("failed to serve: %v", err).WithField("address", address)
		}
	}()

	select {
	case <-ctx.Done():
		log.Infof("%s Shutting down server at %q", logPrefix, address)
		grpcServer.GracefulStop()
	case err := <-errCh:
		return err
	}

	return nil
}

// NewServer returns a new Server instance.
func NewServer(config *Config) *Server {
	return &Server{
		config: config,
	}
}
