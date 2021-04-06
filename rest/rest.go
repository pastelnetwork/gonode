//go:generate goa gen github.com/pastelnetwork/walletnode/rest/design

package rest

import (
	"context"
	"embed"
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

//go:embed swagger
var swagger embed.FS

// Rest represents RESTAPI service.
type Rest struct {
	config          *Config
	shutdownTimeout time.Duration
}

// Run starts RESTAPI service.
func (rest *Rest) Run(ctx context.Context) error {
	addr := net.JoinHostPort(rest.config.Hostname, strconv.Itoa(rest.config.Port))

	restHTTP := httpHandler()

	mux := http.NewServeMux()
	mux.Handle("/", restHTTP)
	mux.Handle("/swagger/swagger.json", restHTTP)
	if rest.config.Swagger {
		mux.Handle("/swagger/", http.FileServer(http.FS(swagger)))
	}

	handler := cors.AllowAll().Handler(mux)

	srv := &http.Server{Addr: addr, Handler: handler}

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
