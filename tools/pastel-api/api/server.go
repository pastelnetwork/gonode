package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/random"
	"github.com/pastelnetwork/gonode/tools/pastel-api/api/services"
)

const (
	shutdownTimeout = time.Second * 5
	logPrefix       = "api"
)

// Server represents RPC API server.
type Server struct {
	services services.Services
}

// Run starts server.
func (server *Server) Run(ctx context.Context, config *Config) error {
	ctx = log.ContextWithPrefix(ctx, logPrefix)

	mux := http.NewServeMux()
	mux.Handle("/", server.httpHandler(ctx))

	addr := net.JoinHostPort(config.Hostname, strconv.Itoa(config.Port))
	srv := &http.Server{Addr: addr, Handler: mux}

	errCh := make(chan error, 1)
	go func() {
		<-ctx.Done()
		log.WithContext(ctx).Infof("Server is shutting down...")

		ctx, cancel := context.WithTimeout(ctx, shutdownTimeout)
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

func (server *Server) httpHandler(ctx context.Context) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer errors.Recover(errors.CheckErrorAndExit)

		w.Header().Set("Content-Type", "application/json")

		var req Request
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			log.WithContext(ctx).WithError(err).Error("Could not parse request")
			return
		}
		r.Body.Close()

		reqID, _ := random.String(8, random.Base62Chars)
		ctx = log.ContextWithPrefix(ctx, fmt.Sprintf("%s-%v", logPrefix, reqID))

		username, _, _ := r.BasicAuth()
		log.WithContext(ctx).WithField("username", username).WithField("req", req).Debug("Request")

		data, err := server.services.Handle(ctx, r, req.Method, req.Params)
		if err != nil {
			resp := newErrorResponse(req.ID, err)
			server.write(ctx, w, http.StatusBadRequest, resp)
			return
		}

		resp := newResponse(req.ID, data)
		server.write(ctx, w, http.StatusOK, resp)
	})
}

func (server *Server) write(ctx context.Context, w http.ResponseWriter, code int, resp *Response) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)

	data, err := resp.Bytes()
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("Failed to generate response")
		return
	}
	if _, err := w.Write(data); err != nil {
		log.WithContext(ctx).WithError(err).Error("Failed to send response")
		return
	}

	log.WithContext(ctx).WithField("resp", resp.Result).Debug("Response")
}

// NewServer returns a new Server instance.
func NewServer(services ...services.Service) *Server {
	return &Server{
		services: services,
	}
}
