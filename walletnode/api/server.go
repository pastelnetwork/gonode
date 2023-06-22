// //go:generate goa gen github.com/pastelnetwork/gonode/walletnode/api/design

package api

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"

	"path/filepath"
	"strings"
	"sync"

	"github.com/pastelnetwork/gonode/common/errgroup"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/pastel"
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
	Run(ctx context.Context) error
}

// Server represents REST API server.
type Server struct {
	config          *Config
	shutdownTimeout time.Duration
	pastelClient    pastel.Client
	services        []service
	fileMappings    *sync.Map // maps unique ID to file path
}

// Run starts listening for the requests.
func (server *Server) Run(ctx context.Context) error {
	ctx = log.ContextWithPrefix(ctx, logPrefix)

	goamux := goahttp.NewMuxer()

	groupServices, ctx := errgroup.WithContext(ctx)
	var servers goahttp.Servers
	for _, service := range server.services {
		groupServices.Go(func() error {
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

	// Serve static files from the "static" directory at the "/static/" path
	mux.Handle("/files/", http.StripPrefix("/files/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract unique ID from request
		uniqueID := filepath.Base(r.URL.Path)

		// Look up file path
		value, ok := server.fileMappings.Load(uniqueID)
		if !ok {
			http.NotFound(w, r)
			return
		}

		filePath, ok := value.(string)
		if !ok {
			http.NotFound(w, r)
			return
		}

		// Extract txid from the file path
		parts := strings.Split(filePath, "/")
		if len(parts) < 2 {
			http.NotFound(w, r)
			return
		}

		txid := parts[len(parts)-2]

		// Set the original filename in the response headers
		w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, parts[len(parts)-1]))

		// Serve the file from the txid directory
		filePath = filepath.Join(server.config.StaticFilesDir, txid, filepath.Base(filePath))
		http.ServeFile(w, r, filePath)

	})))

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

// NewAPIServer returns a new Server instance.
func NewAPIServer(config *Config, filesMap *sync.Map, pastelClient pastel.Client, services ...service) *Server {
	server := &Server{
		pastelClient:    pastelClient,
		config:          config,
		shutdownTimeout: defaultShutdownTimeout,
		services:        services,
		fileMappings:    filesMap,
	}

	return server
}
