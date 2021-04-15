//go:generate goa gen github.com/pastelnetwork/walletnode/api/design

package api

import (
	"context"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/pastelnetwork/go-commons/log"
)

const (
	defaultShutdownTimeout = time.Second * 30
)

// Rest represents RESTAPI service.
type Rest struct {
	config          *Config
	shutdownTimeout time.Duration
}

// Run starts RESTAPI service.
func (rest *Rest) Run(ctx context.Context) error {
	addr := net.JoinHostPort(rest.config.Hostname, strconv.Itoa(rest.config.Port))

	apiHTTP := apiHandler()

	mux := http.NewServeMux()
	mux.Handle("/", apiHTTP)
	mux.Handle("/swagger/swagger.json", apiHTTP)

	if rest.config.Swagger {
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
	ctx, cancel := context.WithTimeout(context.Background(), rest.shutdownTimeout)
	defer cancel()

	err := srv.Shutdown(ctx)
	return err
}

// New returns a new Rest instance.
func New(config *Config) *Rest {
	return &Rest{
		config:          config,
		shutdownTimeout: defaultShutdownTimeout,
	}
}
