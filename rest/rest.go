//go:generate goa gen github.com/pastelnetwork/walletnode/rest/design

package rest

import (
	"context"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/pastelnetwork/go-commons/log"
	"github.com/rs/cors"
)

const (
	defaultShutdownTimeout = time.Second * 30
)

type Rest struct {
	config          *Config
	shutdownTimeout time.Duration
}

func Run(ctx context.Context) error {
	return nil
}

func (rest *Rest) Run(ctx context.Context) error {
	addr := net.JoinHostPort(rest.config.Hostname, strconv.Itoa(rest.config.Port))

	mux := http.NewServeMux()
	mux.Handle("/", httpHandler())
	handler := cors.AllowAll().Handler(mux)

	srv := &http.Server{Addr: addr, Handler: handler}

	errCh := make(chan error)
	go func() {
		log.Infof("[rest] HTTP server listening on %q", addr)
		errCh <- srv.ListenAndServe()
	}()

	select {
	case <-ctx.Done():
		log.Infof("[rest] Shutting down HTTP server at %q", addr)
	case err := <-errCh:
		return err
	}

	// Shutdown gracefully with a 30s timeout.
	ctx, cancel := context.WithTimeout(context.Background(), rest.shutdownTimeout)
	defer cancel()

	err := srv.Shutdown(ctx)
	return err
}

func New(config *Config) *Rest {
	return &Rest{
		config:          config,
		shutdownTimeout: defaultShutdownTimeout,
	}
}
