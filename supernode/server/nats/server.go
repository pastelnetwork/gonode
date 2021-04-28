package nats

import (
	"context"

	"github.com/nats-io/nats-server/v2/server"
)

// Server represents nat service.
type Server struct {
	config *Config
}

// Run starts nat service.
func (srv *Server) Run(ctx context.Context) error {
	nats, err := server.NewServer(&server.Options{
		Host:   srv.config.Hostname,
		Port:   srv.config.Port,
		NoSigs: true,
	})
	if err != nil {
		return err
	}

	nats.SetLoggerV2(NewLogger(), true, true, true)
	nats.Start()

	<-ctx.Done()

	nats.Shutdown()
	nats.WaitForShutdown()
	return nil
}

// NewServer returns a new Server instance.
func NewServer(config *Config) *Server {
	return &Server{
		config: config,
	}
}
