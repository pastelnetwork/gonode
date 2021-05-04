package api

import (
	"context"
	"encoding/json"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/tools/pastel-api/api/services"
)

const (
	shutdownTimeout = time.Second * 5
	logPrefix       = "api"
)

type Server struct {
	services map[string]services.Service
}

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

	log.WithContext(ctx).Infof("Server is ready to handle requests at %q", addr)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return errors.Errorf("error starting server: %w", err)
	}
	defer log.WithContext(ctx).Infof("Server stoped")

	err := <-errCh
	return err
}

func (server *Server) httpHandler(ctx context.Context) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var req Request
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			log.WithContext(ctx).WithError(err).Error("Could not parse request")
			return
		}
		r.Body.Close()

		log.WithContext(ctx).WithField("req", req).Debug("Received new request")

		service, ok := server.services[req.Method]
		if !ok {
			resp := newErrorResponse(req.ID, -32601, "Method not found")
			if err := resp.Write(w, http.StatusNotImplemented); err != nil {
				log.WithContext(ctx).WithError(err).Error("Could not write response")
			}
			return
		}

		data, err := service.Handle(ctx, req.Params)
		if err != nil {
			resp := newErrorResponse(req.ID, -1, err.Error())
			if err := resp.Write(w, http.StatusInternalServerError); err != nil {
				log.WithContext(ctx).WithError(err).Error("Could not write response")
			}
			return
		}

		resp := newResponse(req.ID, data)
		if err := resp.Write(w, http.StatusOK); err != nil {
			log.WithContext(ctx).WithError(err).Error("Could not write response")
		}
	})
}

func NewServer() *Server {
	return &Server{
		services: map[string]services.Service{
			"masternode": services.NewMasterNode(),
		},
	}
}
