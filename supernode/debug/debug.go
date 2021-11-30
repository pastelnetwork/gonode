package debug

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/p2p"

	"net/http"
)

const (
	defaultListenAddr         = "127.0.0.1"
	defaultPollDuration       = 10 * time.Minute
	defaultP2PExpiresDuration = 15 * time.Minute
	defaultP2PLimit           = 10
	defaultP2PMaxLimit        = 100
)

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

// CleanupRequest indicates request structure of cleanup request
type CleanupRequest struct {
	DiscardRatio float64 `json:"discard_ratio"`
}

// StoreReply indicates reply structure of store request
type StoreReply struct {
	Key string `json:"key"`
}

// Service is main point of debug service
type Service struct {
	config       *Config
	p2pClient    p2p.Client
	httpServer   *http.Server
	cleanTracker *CleanTracker
}

// NewService returns debug service
func NewService(config *Config, p2pClient p2p.Client) *Service {
	service := &Service{
		config:    config,
		p2pClient: p2pClient,
	}

	router := mux.NewRouter()

	router.HandleFunc("/p2p/get", service.p2pGet).Methods(http.MethodGet)        // Return list of keys
	router.HandleFunc("/p2p/stats", service.p2pStats).Methods(http.MethodGet)    // Return stats of p2p
	router.HandleFunc("/p2p", service.p2pStore).Methods(http.MethodPost)         // store a data
	router.HandleFunc("/p2p/{key}", service.p2pRetrieve).Methods(http.MethodGet) // retrieve a key
	router.HandleFunc("/p2p/cleanup", service.p2pCleanup).Methods(http.MethodPost)

	service.httpServer = &http.Server{
		Addr:    fmt.Sprintf("%s:%d", defaultListenAddr, config.HTTPPort),
		Handler: router,
	}
	service.cleanTracker = NewCleanTracker(p2pClient)

	return service
}

func responseWithJSON(writer http.ResponseWriter, status int, object interface{}) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(status)
	json.NewEncoder(writer).Encode(object)
}

func (service *Service) contextWithLogPrefix(ctx context.Context) context.Context {
	return log.ContextWithPrefix(ctx, "debug-service")
}

func (service *Service) p2pGet(writer http.ResponseWriter, request *http.Request) {
	ctx := service.contextWithLogPrefix(request.Context())
	log.WithContext(ctx).Info("p2Get")
	var err error

	offset := 0
	limit := defaultP2PLimit

	if offsetQuery := request.URL.Query().Get("offset"); offsetQuery != "" {
		offset, err = strconv.Atoi(offsetQuery)
		if err != nil {
			responseWithJSON(writer, http.StatusForbidden, map[string]string{"error": "Invalid offset input"})
			return
		}
	}

	if limitQuery := request.URL.Query().Get("limit"); limitQuery != "" {
		limit, err = strconv.Atoi(limitQuery)
		if err != nil {
			responseWithJSON(writer, http.StatusForbidden, map[string]string{"error": "Invalid limit input"})
			return
		}
	}

	if limit > defaultP2PMaxLimit {
		responseWithJSON(writer, http.StatusForbidden, map[string]string{"error": "limit out of range, maximum 100 is supported"})
		return
	}

	log.WithContext(ctx).WithFields(log.Fields{
		"limit":  limit,
		"offset": offset,
	}).Info("p2pGetParams")

	keys := service.p2pClient.Keys(ctx, offset, limit)
	if len(keys) == 0 {
		responseWithJSON(writer, http.StatusInternalServerError, map[string]string{"error": "empty key list"})
		return
	}

	responseWithJSON(writer, http.StatusOK, map[string]interface{}{
		"len":    len(keys),
		"offset": offset,
		"limit":  limit,
		"keys":   keys,
	})
}

func (service *Service) p2pStats(writer http.ResponseWriter, request *http.Request) {
	ctx := service.contextWithLogPrefix(request.Context())
	log.WithContext(ctx).Info("p2pStats")
	stats, err := service.p2pClient.Stats(ctx)
	if err != nil {
		responseWithJSON(writer, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	responseWithJSON(writer, http.StatusOK, stats)
}

func (service *Service) p2pRetrieve(writer http.ResponseWriter, request *http.Request) {
	ctx := service.contextWithLogPrefix(request.Context())
	params := mux.Vars(request)

	// TODO : validate key is found or not
	key := params["key"]
	log.WithContext(ctx).WithField("key", key).Info("p2pRetrieve")

	value, err := service.p2pClient.Retrieve(ctx, key)
	if err != nil {
		responseWithJSON(writer, http.StatusNotFound, map[string]string{"error": err.Error()})
		return
	}

	responseWithJSON(writer, http.StatusOK, &RetrieveResponse{
		Key:   key,
		Value: value,
	})
}

func (service *Service) p2pStore(writer http.ResponseWriter, request *http.Request) {
	ctx := service.contextWithLogPrefix(request.Context())
	log.WithContext(ctx).Info("p2pStore")

	var storeRequest StoreRequest
	if err := json.NewDecoder(request.Body).Decode(&storeRequest); err != nil {
		responseWithJSON(writer, http.StatusBadRequest, map[string]string{"message": "Invalid body"})
		return
	}

	key, err := service.p2pClient.Store(ctx, storeRequest.Value)
	if err != nil {
		responseWithJSON(writer, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	// Store to remove key after expires
	service.cleanTracker.Track(key, time.Now().Add(defaultP2PExpiresDuration))

	responseWithJSON(writer, http.StatusOK, &StoreReply{
		Key: key,
	})
}

func (service *Service) p2pCleanup(writer http.ResponseWriter, request *http.Request) {
	ctx := service.contextWithLogPrefix(request.Context())
	log.WithContext(ctx).Info("p2pCleanup")

	var cleanupRequest CleanupRequest
	if err := json.NewDecoder(request.Body).Decode(&cleanupRequest); err != nil {
		responseWithJSON(writer, http.StatusBadRequest, map[string]string{"message": "Invalid body"})
		return
	}
	if cleanupRequest.DiscardRatio < 0.0 || cleanupRequest.DiscardRatio > 1.0 {
		responseWithJSON(writer, http.StatusBadRequest, map[string]string{"message": "discard_ratio must be in range 0.0 - 1.0"})
		return
	}

	err := service.p2pClient.Cleanup(ctx, cleanupRequest.DiscardRatio)
	if err != nil {
		responseWithJSON(writer, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	responseWithJSON(writer, http.StatusOK, map[string]string{"message": "cleanup successfull"})
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
		log.WithContext(ctx).WithField("port", service.config.HTTPPort).Info("Http server started")
		err := service.httpServer.ListenAndServe()
		log.WithContext(ctx).WithError(err).Info("Http server stopped")
	}()

	// waiting until context is cancelled
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(defaultPollDuration):
			service.cleanTracker.CheckExpires(ctx)
		}
	}
}
