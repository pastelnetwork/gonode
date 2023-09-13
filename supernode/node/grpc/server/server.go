package server

import (
	"context"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/pastelnetwork/gonode/common/errgroup"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/net/credentials/alts"
	"github.com/pastelnetwork/gonode/supernode/node/grpc/server/middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/keepalive"
)

type service interface {
	Desc() *grpc.ServiceDesc
}

// Server represents supernode server
type Server struct {
	config    *Config
	services  []service
	name      string
	secClient alts.SecClient
	secInfo   *alts.SecInfo
}

// Run starts the server
func (server *Server) Run(ctx context.Context) error {
	grpclog.SetLoggerV2(log.DefaultLogger)
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
	if server.secClient == nil || server.secInfo == nil {
		log.WithContext(ctx).Errorln("secClient or secInfo don't initialize")
		return nil
	}

	// Define the keep-alive parameters
	kaParams := keepalive.ServerParameters{
		MaxConnectionIdle:     2 * time.Hour,
		MaxConnectionAge:      2 * time.Hour,
		MaxConnectionAgeGrace: 1 * time.Hour,
		Time:                  1 * time.Hour,
		Timeout:               2 * time.Minute,
	}

	// Define the keep-alive enforcement policy
	kaPolicy := keepalive.EnforcementPolicy{
		MinTime:             3 * time.Minute, // Minimum time a client should wait before sending keep-alive probes
		PermitWithoutStream: true,            // Only allow pings when there are active streams
	}

	var grpcServer *grpc.Server
	if os.Getenv("INTEGRATION_TEST_ENV") == "true" {
		grpcServer = grpc.NewServer(middleware.UnaryInterceptor(), middleware.StreamInterceptor(), grpc.MaxSendMsgSize(35000000),
			grpc.MaxRecvMsgSize(35000000), grpc.KeepaliveParams(kaParams), // Use the keep-alive parameters
			grpc.KeepaliveEnforcementPolicy(kaPolicy))
	} else {
		grpcServer = grpc.NewServer(middleware.UnaryInterceptor(), middleware.StreamInterceptor(),
			middleware.AltsCredential(server.secClient, server.secInfo), grpc.MaxSendMsgSize(35000000),
			grpc.MaxRecvMsgSize(35000000), grpc.KeepaliveParams(kaParams), // Use the keep-alive parameters
			grpc.KeepaliveEnforcementPolicy(kaPolicy))
	}

	for _, service := range server.services {
		log.WithContext(ctx).Debugf("Register service %q", service.Desc().ServiceName)
		grpcServer.RegisterService(service.Desc(), service)
	}

	return grpcServer
}

// New returns a new Server instance.
func New(config *Config, name string, secClient alts.SecClient, secInfo *alts.SecInfo, services ...service) *Server {
	return &Server{
		config:    config,
		secClient: secClient,
		secInfo:   secInfo,
		services:  services,
		name:      name,
	}
}
