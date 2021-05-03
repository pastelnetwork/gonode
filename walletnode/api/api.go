//go:generate goa gen github.com/pastelnetwork/gonode/walletnode/api/design

package api

import (
	"context"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/walletnode/api/docs"

	goahttp "goa.design/goa/v3/http"
	goahttpmiddleware "goa.design/goa/v3/http/middleware"
)

const (
	defaultShutdownTimeout = time.Second * 30
	logPrefix              = "api"
)

type service interface {
	Mount(ctx context.Context, mux goahttp.Muxer) goahttp.Server
}

// API represents RESTAPI service.
type API struct {
	config          *Config
	shutdownTimeout time.Duration
	services        []service
}

// Run startworks RESTAPI service.
func (api *API) Run(ctx context.Context) error {
	ctx = context.WithValue(ctx, log.PrefixKey, logPrefix)

	addr := net.JoinHostPort(api.config.Hostname, strconv.Itoa(api.config.Port))

	apiHTTP := api.handler(ctx)

	mux := http.NewServeMux()
	mux.Handle("/", apiHTTP)
	mux.Handle("/swagger/swagger.json", apiHTTP)

	if api.config.Swagger {
		mux.Handle("/swagger/", http.FileServer(http.FS(docs.SwaggerContent)))
	}

	srv := &http.Server{Addr: addr, Handler: mux}

	errCh := make(chan error)
	go func() {
		log.WithContext(ctx).Infof("Server listening on %q", addr)
		errCh <- srv.ListenAndServe()
	}()

	select {
	case <-ctx.Done():
		log.WithContext(ctx).Infof("Shutting down server at %q", addr)
	case err := <-errCh:
		return err
	}

	// Shutdown gracefully with a 30s timeout.
	ctx, cancel := context.WithTimeout(context.Background(), api.shutdownTimeout)
	defer cancel()

	err := srv.Shutdown(ctx)
	return err
}

func (api *API) handler(ctx context.Context) http.Handler {
	mux := goahttp.NewMuxer()

	var servers goahttp.Servers
	for _, service := range api.services {
		servers = append(servers, service.Mount(ctx, mux))
	}
	servers.Use(goahttpmiddleware.Debug(mux, os.Stdout))

	var handler http.Handler = mux

	handler = Log(ctx)(handler)
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
