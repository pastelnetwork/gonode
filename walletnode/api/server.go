// //go:generate goa gen github.com/pastelnetwork/gonode/walletnode/api/design

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
	"golang.org/x/sync/errgroup"

	goahttp "goa.design/goa/v3/http"
	goahttpmiddleware "goa.design/goa/v3/http/middleware"
)

const (
	defaultShutdownTimeout = time.Second * 30
	logPrefix              = "server"
)

type service interface {
	Mount(ctx context.Context, mux goahttp.Muxer) goahttp.Server
	Run(ctx context.Context) error
}

// Server represents REST API server.
type Server struct {
	config          *Config
	shutdownTimeout time.Duration
	services        []service
}

// Run starts listening for the requests.
func (server *Server) Run(ctx context.Context) error {
	ctx = log.ContextWithPrefix(ctx, logPrefix)

	goamux := goahttp.NewMuxer()

	groupServices, ctx := errgroup.WithContext(ctx)
	var servers goahttp.Servers
	for _, service := range server.services {
		groupServices.Go(func() (err error) {
			defer errors.Recover(func(recErr error) { err = recErr })
			return service.Run(ctx)
		})
		servers = append(servers, service.Mount(ctx, goamux))
	}
	servers.Use(goahttpmiddleware.Debug(goamux, os.Stdout))

	var handler http.Handler = goamux

	handler = Recovery()(handler)
	handler = Log(ctx)(handler)
	handler = goahttpmiddleware.RequestID()(handler)

	mux := http.NewServeMux()
	mux.Handle("/", handler)
	mux.Handle("/swagger/swagger.json", handler)

	if server.config.Swagger {
		mux.Handle("/swagger/", http.FileServer(http.FS(docs.SwaggerContent)))
	}

	addr := net.JoinHostPort(server.config.Hostname, strconv.Itoa(server.config.Port))
	srv := &http.Server{Addr: addr, Handler: mux}

	errCh := make(chan error, 1)
	go func() {
		defer errors.Recover(log.Fatal)

		<-ctx.Done()
		log.WithContext(ctx).Infof("Shutting down REST server on %q", addr)

		ctx, cancel := context.WithTimeout(ctx, server.shutdownTimeout)
		defer cancel()

		srv.SetKeepAlivesEnabled(false)
		if err := srv.Shutdown(ctx); err != nil {
			errCh <- errors.Errorf("gracefully shutdown the server: %w", err)
		}
		close(errCh)
	}()

	log.WithContext(ctx).Infof("REST server is listening on %q", addr)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return errors.Errorf("error starting server: %w", err)
	}
	defer log.WithContext(ctx).Infof("Server stoped")

	if err := groupServices.Wait(); err != nil {
		return err
	}
	return <-errCh
}

// NewServer returns a new Server instance.
func NewServer(config *Config, services ...service) *Server {
	return &Server{
		config:          config,
		shutdownTimeout: defaultShutdownTimeout,
		services:        services,
	}
}
