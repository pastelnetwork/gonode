package main

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/p2p"
)

type retrieveResponse struct {
	Key   string `json:"key"`
	Value []byte `json:"value"`
}

type storeRequest struct {
	Value []byte `json:"value"`
}

type storeReply struct {
	Key string `json:"key"`
}

type debugService struct {
	p2pClient  p2p.Client
	httpServer *http.Server
	localKeys  map[string]bool
}

func newDebugP2PService(p2pClient p2p.Client) *debugService {
	service := &debugService{
		p2pClient: p2pClient,
		localKeys: make(map[string]bool),
	}
	router := mux.NewRouter()
	router.HandleFunc("/p2p", service.p2pStore).Methods(http.MethodPost)                                   // store a data
	router.HandleFunc("/p2p/{key}", service.p2pRetrieve).Methods(http.MethodGet)                           // retrieve a key
	router.HandleFunc("/p2p/{ratio}", service.p2pRemoveRandomFilesByGivenRatio).Methods(http.MethodDelete) // remove keys given ratio

	service.httpServer = &http.Server{
		Addr:    fmt.Sprintf("%s:%s", "0.0.0.0", rawPorts[1]),
		Handler: router,
	}

	return service
}

func (service *debugService) contextWithLogPrefix(ctx context.Context) context.Context {
	return log.ContextWithPrefix(ctx, "debug-service")
}

func responseWithJSON(writer http.ResponseWriter, status int, object interface{}) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(status)
	json.NewEncoder(writer).Encode(object)
}

func (service *debugService) p2pStore(writer http.ResponseWriter, request *http.Request) {
	ctx := service.contextWithLogPrefix(request.Context())
	log.WithContext(ctx).Info("p2pStore")

	var storeRequest storeRequest
	if err := json.NewDecoder(request.Body).Decode(&storeRequest); err != nil {
		responseWithJSON(writer, http.StatusBadRequest, map[string]string{"message": "Invalid body"})
		return
	}

	key, err := service.p2pClient.Store(ctx, storeRequest.Value)
	if err != nil {
		responseWithJSON(writer, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	service.localKeys[key] = true
	responseWithJSON(writer, http.StatusOK, &storeReply{
		Key: key,
	})
}

func (service *debugService) p2pRetrieve(writer http.ResponseWriter, request *http.Request) {
	ctx := service.contextWithLogPrefix(request.Context())
	params := mux.Vars(request)

	key, ok := params["key"]
	if !ok {
		writer.Write([]byte("key param url not found"))
		writer.WriteHeader(http.StatusNotFound)
		return
	}
	log.WithContext(ctx).WithField("key", key).Info("p2pRetrieve")

	value, err := service.p2pClient.Retrieve(ctx, key)
	if err != nil {
		responseWithJSON(writer, http.StatusNotFound, map[string]string{"error": err.Error()})
		return
	}

	responseWithJSON(writer, http.StatusOK, &retrieveResponse{
		Key:   key,
		Value: value,
	})
}

func (service *debugService) p2pRemoveRandomFilesByGivenRatio(writer http.ResponseWriter, request *http.Request) {
	ctx := service.contextWithLogPrefix(request.Context())
	params := mux.Vars(request)

	ratioStr, ok := params["ratio"]
	if !ok {
		writer.Write([]byte("ratio param url not found"))
		writer.WriteHeader(http.StatusNotFound)
		return
	}
	ratio, err := strconv.ParseFloat(ratioStr, 64)
	if err != nil {
		ratio = 0.2
	}

	localKeyList, err := getLocalKeys(ctx)
	if err != nil {
		writer.Write([]byte("local keys not found"))
		writer.WriteHeader(http.StatusNotFound)
		return
	}

	for _, key := range localKeyList {
		service.localKeys[key] = true
	}

	log.WithContext(ctx).Infof("DELETE KEYS WITH RATIO %f FROM LOCAL KEYS %v", ratio, service.localKeys)
	var removedKeys = make([]string, 0)
	for key, ok := range service.localKeys {
		if ok {
			if rd := rand.Float64(); rd <= ratio {
				err = service.p2pClient.Delete(ctx, key)
				if err != nil {
					log.WithContext(ctx).WithError(err).Errorf("could not delete key %s", key)
					continue
				}
				removedKeys = append(removedKeys, key)
			}
		}
	}

	responseWithJSON(writer, http.StatusOK, removedKeys)
}

func (service *debugService) Run(ctx context.Context) error {
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
		log.WithContext(ctx).WithField("port", rawPorts[1]).Info("Http server started")
		err := service.httpServer.ListenAndServe()
		log.WithContext(ctx).WithError(err).Info("Http server stopped")
	}()

	// waiting until context is cancelled
	<-ctx.Done()
	return ctx.Err()
}
