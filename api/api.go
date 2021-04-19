//go:generate goa gen github.com/pastelnetwork/walletnode/api/design

package api

import (
	"context"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/pastelnetwork/go-commons/log"
	"github.com/pastelnetwork/walletnode/services/artwork"
)

const (
	defaultShutdownTimeout = time.Second * 30
)

// API represents RESTAPI service.
type API struct {
	config          *Config
	shutdownTimeout time.Duration
	artwork         *artwork.Service
}

// Run startworks RESTAPI service.
func (api *API) Run(ctx context.Context) error {
	addr := net.JoinHostPort(api.config.Hostname, strconv.Itoa(api.config.Port))

	apiHTTP := api.handler()

	mux := http.NewServeMux()
	mux.Handle("/", apiHTTP)
	mux.Handle("/swagger/swagger.json", apiHTTP)

	if api.config.Swagger {
		mux.Handle("/swagger/", swaggerHandler())
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

// New returns a new API instance.
func New(config *Config, artwork *artwork.Service) *API {
	return &API{
		config:          config,
		shutdownTimeout: defaultShutdownTimeout,
		artwork:         artwork,
	}
}
