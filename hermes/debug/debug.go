package debug

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gorilla/mux"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/version"
	"github.com/pastelnetwork/gonode/hermes/service"

	"net/http"
)

const (
	defaultListenAddr   = "0.0.0.0"
	defaultPollDuration = 10 * time.Minute
)

// contains http service providing debug services to user

// Service is main point of debug service
type Service struct {
	config     *Config
	httpServer *http.Server
	hermes     service.SvcInterface
}

// NewService returns debug service
func NewService(config *Config, svc service.SvcInterface) *Service {
	service := &Service{
		config: config,
		hermes: svc,
	}

	router := mux.NewRouter()

	router.HandleFunc("/hermes/stats", service.hermesStats).Methods(http.MethodGet) // Return stats of hermes
	router.HandleFunc("/health", service.hermesHealth).Methods(http.MethodGet)

	service.httpServer = &http.Server{
		Addr:    fmt.Sprintf("%s:%d", defaultListenAddr, config.HTTPPort),
		Handler: router,
	}

	return service
}

func responseWithJSON(writer http.ResponseWriter, status int, object interface{}) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(status)
	json.NewEncoder(writer).Encode(object)
}

func (service *Service) contextWithLogPrefix(ctx context.Context) context.Context {
	return log.ContextWithPrefix(ctx, "hermes-debug-service")
}

func (service *Service) hermesStats(writer http.ResponseWriter, request *http.Request) {
	ctx := service.contextWithLogPrefix(request.Context())
	log.WithContext(ctx).Info("hermesStats")
	stats, err := service.hermes.Stats(ctx)
	if err != nil {
		responseWithJSON(writer, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	stats["version"] = version.Version()

	responseWithJSON(writer, http.StatusOK, stats)
}

func (service *Service) hermesHealth(writer http.ResponseWriter, request *http.Request) {
	ctx := service.contextWithLogPrefix(request.Context())
	log.WithContext(ctx).Info("hermesHealth")

	responseWithJSON(writer, http.StatusOK, fmt.Sprintf("%s - %s", "hermes", version.Version()))
}

// Run start update stats of system periodically
func (service *Service) Run(ctx context.Context) error {
	ctx = service.contextWithLogPrefix(ctx)
	log.WithContext(ctx).Info("Service started")
	defer log.WithContext(ctx).Info("Service stopped")

	defer func() {
		// stop http server
		stopCtx, stopCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer stopCancel()
		service.httpServer.Shutdown(stopCtx)
	}()

	log.WithContext(ctx).WithField("port", service.config.HTTPPort).Info("Http server started")
	err := service.httpServer.ListenAndServe()
	log.WithContext(ctx).WithError(err).Info("Http server stopped")

	return err
}
