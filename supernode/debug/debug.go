package debug

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gorilla/mux"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/p2p"

	"net/http"
)

const defaultListenAddr = "127.0.0.1"

// contains http service providing debug services to user

// RetrieveResponse indicates response structure of retrieve request
type RetrieveResponse struct {
	Key   string `json:"key"`
	Value []byte `json:"value"`
}

// StoreRequest indicates request structure of store request
type StoreRequest struct {
	Value []byte `json:"value"`
}

// StoreRequest indicates reply structure of store request
type StoreReply struct {
	Key string `json:"key"`
}

// Service is main point of debug service
type Service struct {
	config     *Config
	p2pClient  p2p.Client
	httpServer *http.Server
}

// NewService returns debug service
func NewService(config *Config, p2pClient p2p.Client) *Service {
	service := &Service{
		config:    config,
		p2pClient: p2pClient,
	}

	router := mux.NewRouter()

	router.HandleFunc("/p2p/stats", service.p2pStats).Methods(http.MethodGet)
	router.HandleFunc("/p2p", service.p2pStore).Methods(http.MethodPost)
	router.HandleFunc("/p2p/{key}", service.p2pRetrieve).Methods(http.MethodGet)

	service.httpServer = &http.Server{
		Addr:    fmt.Sprintf("%s:%d", defaultListenAddr, config.HttpPort),
		Handler: router,
	}

	return service
}

func responseWithJson(writer http.ResponseWriter, status int, object interface{}) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(status)
	json.NewEncoder(writer).Encode(object)
}

func (service *Service) contextWithLogPrefix(ctx context.Context) context.Context {
	return log.ContextWithPrefix(ctx, "debug-service")
}

func (service *Service) p2pStats(writer http.ResponseWriter, request *http.Request) {
	ctx := service.contextWithLogPrefix(request.Context())
	log.WithContext(ctx).Info("p2pStats")
	stats, err := service.p2pClient.Stats(ctx)
	if err != nil {
		responseWithJson(writer, http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	responseWithJson(writer, http.StatusOK, stats)
}

func (service *Service) p2pRetrieve(writer http.ResponseWriter, request *http.Request) {
	ctx := service.contextWithLogPrefix(request.Context())
	params := mux.Vars(request)

	// TODO : validate key is found or not
	key := params["key"]
	log.WithContext(ctx).WithField("key", key).Info("p2pRetrieve")

	value, err := service.p2pClient.Retrieve(ctx, key)
	if err != nil {
		responseWithJson(writer, http.StatusNotFound, map[string]string{"error": err.Error()})
	}

	responseWithJson(writer, http.StatusOK, &RetrieveResponse{
		Key:   key,
		Value: value,
	})
}

func (service *Service) p2pStore(writer http.ResponseWriter, request *http.Request) {
	ctx := service.contextWithLogPrefix(request.Context())
	log.WithContext(ctx).Info("p2pStore")

	var storeRequest StoreRequest
	if err := json.NewDecoder(request.Body).Decode(&storeRequest); err != nil {
		responseWithJson(writer, http.StatusBadRequest, map[string]string{"message": "Invalid body"})
		return
	}

	key, err := service.p2pClient.Store(ctx, storeRequest.Value)
	if err != nil {
		responseWithJson(writer, http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	responseWithJson(writer, http.StatusOK, &StoreReply{
		Key: key,
	})
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

	// start http service
	go func() {
		log.WithContext(ctx).WithField("port", service.config.HttpPort).Info("Http server started")
		err := service.httpServer.ListenAndServe()
		log.WithContext(ctx).WithError(err).Info("Http server stopped")
	}()

	// waiting until context is cancelled
	<-ctx.Done()
	return ctx.Err()
}
