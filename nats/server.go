package nats

import (
	"context"

	natsserver "github.com/nats-io/nats-server/v2/server"
)

// Server represents nat-server
type Server struct {
	config *Config
}

// Run start nat server
func (server *Server) Run(ctx context.Context) error {
	srv, err := natsserver.NewServer(&natsserver.Options{
		Host:   server.config.Hostname,
		Port:   server.config.Port,
		NoSigs: true,
	})
	if err != nil {
		return err
	}

	srv.SetLoggerV2(&Logger{}, true, true, true)
	srv.Start()

	<-ctx.Done()

	srv.Shutdown()
	srv.WaitForShutdown()
	return nil
}

// NewServer returns a new Server instance
func NewServer(config *Config) *Server {
	return &Server{
		config: config,
	}
}
