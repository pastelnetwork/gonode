//go:generate goa gen github.com/pastelnetwork/gonode/walletnode/server/design

package api

import (
	"context"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/walletnode/api/docs"

	goahttp "goa.design/goa/v3/http"
	goahttpmiddleware "goa.design/goa/v3/http/middleware"
)

const (
	defaultShutdownTimeout = time.Second * 30
	logPrefix              = "server"
)

type service interface {
	Mount(ctx context.Context, mux goahttp.Muxer) goahttp.Server
}

// Server represents RESTAPI service.
type Server struct {
	config          *Config
	shutdownTimeout time.Duration
	services        []service
}

// Run startworks RESTAPI service.
func (server *Server) Run(ctx context.Context) error {
	ctx = log.ContextWithPrefix(ctx, logPrefix)

	apiHTTP := server.handler(ctx)

	mux := http.NewServeMux()
	mux.Handle("/", apiHTTP)
	mux.Handle("/swagger/swagger.json", apiHTTP)

	if server.config.Swagger {
		mux.Handle("/swagger/", http.FileServer(http.FS(docs.SwaggerContent)))
	}

	addr := net.JoinHostPort(server.config.Hostname, strconv.Itoa(server.config.Port))
	srv := &http.Server{Addr: addr, Handler: mux}

	errCh := make(chan error, 1)
	go func() {
		<-ctx.Done()
		log.WithContext(ctx).Infof("Server is shutting down...")

		ctx, cancel := context.WithTimeout(ctx, server.shutdownTimeout)
		defer cancel()

		srv.SetKeepAlivesEnabled(false)
		if err := srv.Shutdown(ctx); err != nil {
			errCh <- errors.Errorf("gracefully shutdown the server: %w", err)

		}
		close(errCh)
	}()

	log.WithContext(ctx).Infof("Server is listening on %q", addr)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return errors.Errorf("error starting server: %w", err)
	}
	defer log.WithContext(ctx).Infof("Server stoped")

	err := <-errCh
	return err
}

func (server *Server) handler(ctx context.Context) http.Handler {
	mux := goahttp.NewMuxer()

	var servers goahttp.Servers
	for _, service := range server.services {
		servers = append(servers, service.Mount(ctx, mux))
	}
	servers.Use(goahttpmiddleware.Debug(mux, os.Stdout))

	var handler http.Handler = mux

	handler = Log(ctx)(handler)
	handler = goahttpmiddleware.RequestID()(handler)

	return handler
}

// NewServer returns a new Server instance.
func NewServer(config *Config, services ...service) *Server {
	return &Server{
		config:          config,
		shutdownTimeout: defaultShutdownTimeout,
		services:        services,
	}
}
