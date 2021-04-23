//go:generate goa gen github.com/pastelnetwork/walletnode/api/design

package api

import (
	"context"
	"embed"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/pastelnetwork/go-commons/log"

	goahttp "goa.design/goa/v3/http"
	goahttpmiddleware "goa.design/goa/v3/http/middleware"
)

const (
	defaultShutdownTimeout = time.Second * 30
)

//go:embed swagger
var swaggerContent embed.FS

type service interface {
	Mount(mux goahttp.Muxer) goahttp.Server
}

// API represents RESTAPI service.
type API struct {
	config          *Config
	shutdownTimeout time.Duration
	services        []service
}

// Run startworks RESTAPI service.
func (api *API) Run(ctx context.Context) error {
	addr := net.JoinHostPort(api.config.Hostname, strconv.Itoa(api.config.Port))

	apiHTTP := api.handler()

	mux := http.NewServeMux()
	mux.Handle("/", apiHTTP)
	mux.Handle("/swagger/swagger.json", apiHTTP)

	if api.config.Swagger {
		mux.Handle("/swagger/", http.FileServer(http.FS(swaggerContent)))
	}

	srv := &http.Server{Addr: addr, Handler: mux}

	errCh := make(chan error)
	go func() {
		log.Infof("%s HTTP server listening on %q", logPrefix, addr)
		errCh <- srv.ListenAndServe()
	}()

	select {
	case <-ctx.Done():
		log.Infof("%s Shutting down HTTP server at %q", logPrefix, addr)
	case err := <-errCh:
		return err
	}

	// Shutdown gracefully with a 30s timeout.
	ctx, cancel := context.WithTimeout(context.Background(), api.shutdownTimeout)
	defer cancel()

	err := srv.Shutdown(ctx)
	return err
}

func (api *API) handler() http.Handler {
	mux := goahttp.NewMuxer()

	var servers goahttp.Servers
	for _, service := range api.services {
		servers = append(servers, service.Mount(mux))
	}
	servers.Use(goahttpmiddleware.Debug(mux, os.Stdout))

	var handler http.Handler = mux

	handler = Log()(handler)
	handler = goahttpmiddleware.RequestID()(handler)

	return handler
}

// New returns a new API instance.
func New(config *Config, services ...service) *API {
	return &API{
		config:          config,
		shutdownTimeout: defaultShutdownTimeout,
		services:        services,
	}
}
