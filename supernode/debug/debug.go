package debug

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/mux"
	json "github.com/json-iterator/go"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/storage/local"
	"github.com/pastelnetwork/gonode/common/version"
	"github.com/pastelnetwork/gonode/p2p"
	"github.com/pastelnetwork/gonode/pastel"
	"github.com/pastelnetwork/gonode/supernode/node/grpc/server"
	"github.com/pastelnetwork/gonode/supernode/services/common"
	"github.com/pastelnetwork/gonode/supernode/services/storagechallenge"
)

const (
	defaultListenAddr           = "0.0.0.0"
	defaultPollDuration         = 10 * time.Minute
	defaultP2PExpiresDuration   = 15 * time.Minute
	defaultP2PLimit             = 10
	defaultP2PMaxLimit          = 100
	rateLimitMaxRequestsPerHour = 7
)

// contains http service providing debug services to user

// RetrieveResponse indicates response structure of retrieve request
type RetrieveResponse struct {
	Key   string `json:"key"`
	Value []byte `json:"value"`
}

// DeleteRequest indicates request structure of store request
type DeleteRequest struct {
	Key string `json:"key"`
}

// StoreRequest indicates request structure of store request
type StoreRequest struct {
	Value []byte `json:"value"`
}

// LocalStoreRequest indicates request structure of store request
type LocalStoreRequest struct {
	Key   string `json:"key"`
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
	config        *Config
	p2pClient     p2p.Client
	httpServer    *http.Server
	metricsServer *http.Server
	cleanTracker  *CleanTracker
	scService     *storagechallenge.SCService
	trackedCon    map[string][]*server.TrackedConn
	rateLimiter   *RateLimiter
}

// NewService returns debug service
func NewService(config *Config, p2pClient p2p.Client, srvc *storagechallenge.SCService, openCons map[string][]*server.TrackedConn) *Service {
	service := &Service{
		config:      config,
		p2pClient:   p2pClient,
		scService:   srvc,
		trackedCon:  openCons,
		rateLimiter: NewRateLimiter(rateLimitMaxRequestsPerHour, time.Hour),
	}

	router := mux.NewRouter()

	router.HandleFunc("/local_p2p", service.p2pLocalStore).Methods(http.MethodPost)                    // store data only in local node
	router.HandleFunc("/p2p/stats", service.p2pStats).Methods(http.MethodGet)                          // Return stats of p2p
	router.HandleFunc("/p2p", service.p2pStore).Methods(http.MethodPost)                               // store a data
	router.HandleFunc("/p2p/{key}", service.p2pRetrieve).Methods(http.MethodGet)                       // retrieve a key
	router.HandleFunc("/storage/challenges", service.storageChallengeRetrieve).Methods(http.MethodGet) // retrieve storage challenge data
	router.HandleFunc("/storage/cleanup", service.storageChallengeCleanup).Methods(http.MethodGet)
	router.HandleFunc("/p2p/remove", service.p2pDelete).Methods(http.MethodPost)                                // retrieve a key
	router.HandleFunc("/selfhealing/challenges", service.selfHealingChallengesRetrieve).Methods(http.MethodGet) // retrieve a key
	router.HandleFunc("/selfhealing/cleanup", service.selfHealingChallengeCleanup).Methods(http.MethodDelete)   // cleanup self-healing data from history db
	router.HandleFunc("/health", service.p2pHealth).Methods(http.MethodGet)
	router.HandleFunc("/storage/initiate", service.storageInitiate).Methods(http.MethodGet)
	router.HandleFunc("/open/cons", service.getOpenCons).Methods(http.MethodGet)

	service.httpServer = &http.Server{
		Addr:    fmt.Sprintf("%s:%d", defaultListenAddr, config.HTTPPort),
		Handler: router,
	}
	service.cleanTracker = NewCleanTracker(p2pClient)

	metricsRouter := mux.NewRouter()

	router.HandleFunc("/metrics", service.metrics).Methods(http.MethodGet)
	router.HandleFunc("/sh_trigger", service.shTrigger).Methods(http.MethodGet)
	router.HandleFunc("/sh_challenge", service.shChallenge).Methods(http.MethodGet)
	service.metricsServer = &http.Server{
		Addr:    fmt.Sprintf("%s:%d", defaultListenAddr, config.MetricsPort),
		Handler: metricsRouter,
	}

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

func (service *Service) p2pStats(writer http.ResponseWriter, request *http.Request) {
	ctx := service.contextWithLogPrefix(request.Context())
	log.WithContext(ctx).Info("p2pStats")
	stats, err := service.p2pClient.Stats(ctx)
	if err != nil {
		responseWithJSON(writer, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	stats["version"] = version.Version()

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

func (service *Service) storageChallengeCleanup(writer http.ResponseWriter, request *http.Request) {
	ctx := service.contextWithLogPrefix(request.Context())

	store, err := local.OpenHistoryDB()
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("Error Opening DB")
	}
	defer store.CloseHistoryDB(ctx)

	err = store.CleanupStorageChallenges()
	if err != nil {
		responseWithJSON(writer, http.StatusNotFound, map[string]string{"error": err.Error()})
		return
	}

	responseWithJSON(writer, http.StatusOK, make(map[string]interface{}))
}

func (service *Service) storageChallengeRetrieve(writer http.ResponseWriter, request *http.Request) {
	ctx := service.contextWithLogPrefix(request.Context())

	store, err := local.OpenHistoryDB()
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("Error Opening DB")
	}
	defer store.CloseHistoryDB(ctx)

	challenges, err := store.QueryStorageChallengeMessage("", 0)
	if err != nil {
		responseWithJSON(writer, http.StatusNotFound, map[string]string{"error": err.Error()})
		return
	}

	responseWithJSON(writer, http.StatusOK, challenges)
}

func (service *Service) selfHealingChallengesRetrieve(writer http.ResponseWriter, request *http.Request) {
	ctx := service.contextWithLogPrefix(request.Context())

	store, err := local.OpenHistoryDB()
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("Error Opening DB")
	}
	defer store.CloseHistoryDB(ctx)

	challenges, err := store.QuerySelfHealingChallenges()
	if err != nil {
		responseWithJSON(writer, http.StatusNotFound, map[string]string{"error": err.Error()})
		return
	}

	responseWithJSON(writer, http.StatusOK, challenges)
}

func (service *Service) selfHealingChallengeCleanup(writer http.ResponseWriter, request *http.Request) {
	ctx := service.contextWithLogPrefix(request.Context())

	store, err := local.OpenHistoryDB()
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("Error Opening DB")
	}
	defer store.CloseHistoryDB(ctx)

	err = store.CleanupSelfHealingChallenges()
	if err != nil {
		responseWithJSON(writer, http.StatusNotFound, map[string]string{"error": err.Error()})
		return
	}

	responseWithJSON(writer, http.StatusOK, make(map[string]interface{}))
}

func (service *Service) storageInitiate(writer http.ResponseWriter, _ *http.Request) {
	ctx := context.Background()
	log.WithContext(ctx).Info("storage challenge initiate")
	task := service.scService.NewSCTask()
	go task.GenerateStorageChallenges(ctx)

	responseWithJSON(writer, http.StatusOK, fmt.Sprintf("ok %s", "storage challenge initiated"))
}

func (service *Service) getOpenCons(writer http.ResponseWriter, _ *http.Request) {
	ctx := context.Background()
	log.WithContext(ctx).Info("returning open cons")

	responseWithJSON(writer, http.StatusOK, service.trackedCon)
}

func (service *Service) p2pHealth(writer http.ResponseWriter, request *http.Request) {
	ctx := service.contextWithLogPrefix(request.Context())
	log.WithContext(ctx).Info("p2pHealth")

	responseWithJSON(writer, http.StatusOK, fmt.Sprintf("%s - %s", "p2p", version.Version()))
}

func (service *Service) p2pStore(writer http.ResponseWriter, request *http.Request) {
	ctx := service.contextWithLogPrefix(request.Context())
	log.WithContext(ctx).Info("p2pStore")

	var storeRequest StoreRequest
	if err := json.NewDecoder(request.Body).Decode(&storeRequest); err != nil {
		responseWithJSON(writer, http.StatusBadRequest, map[string]string{"message": "Invalid body"})
		return
	}

	key, err := service.p2pClient.Store(ctx, storeRequest.Value, common.P2PDebug)
	if err != nil {
		responseWithJSON(writer, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	// Store to remove key after expires
	service.cleanTracker.Track(key, time.Now().UTC().Add(defaultP2PExpiresDuration))

	responseWithJSON(writer, http.StatusOK, &StoreReply{
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

func (service *Service) p2pLocalStore(writer http.ResponseWriter, request *http.Request) {
	ctx := service.contextWithLogPrefix(request.Context())
	log.WithContext(ctx).Info("p2pLocalStore")

	var storeRequest LocalStoreRequest
	if err := json.NewDecoder(request.Body).Decode(&storeRequest); err != nil {
		responseWithJSON(writer, http.StatusBadRequest, map[string]string{"message": "Invalid body"})
		return
	}

	key, err := service.p2pClient.LocalStore(ctx, storeRequest.Key, storeRequest.Value)
	if err != nil {
		responseWithJSON(writer, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	// Store to remove key after expires
	service.cleanTracker.Track(key, time.Now().UTC().Add(defaultP2PExpiresDuration))

	responseWithJSON(writer, http.StatusOK, &StoreReply{
		Key: key,
	})
}

func (service *Service) p2pDelete(writer http.ResponseWriter, request *http.Request) {
	ctx := service.contextWithLogPrefix(request.Context())
	log.WithContext(ctx).Info("p2pDelete")

	var delRequest DeleteRequest
	if err := json.NewDecoder(request.Body).Decode(&delRequest); err != nil {
		responseWithJSON(writer, http.StatusBadRequest, map[string]string{"message": "Invalid body"})
		return
	}

	err := service.p2pClient.Delete(ctx, delRequest.Key)
	if err != nil {
		responseWithJSON(writer, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	responseWithJSON(writer, http.StatusOK, &StoreReply{
		Key: delRequest.Key,
	})
}

func (service *Service) shChallenge(writer http.ResponseWriter, request *http.Request) {
	ctx := service.contextWithLogPrefix(request.Context())

	// 1. Retrieve Authorization Header (Passphrase)
	passphrase := request.Header.Get("Authorization")

	// 2. Retrieve Pastel ID from Query Param
	pid := request.URL.Query().Get("pid")
	if pid == "" {
		responseWithJSON(writer, http.StatusBadRequest, map[string]string{"error": "Missing pid parameter"})
		return
	}

	_, err := service.scService.PastelClient.Sign(ctx, []byte{1, 2, 3}, pid, passphrase, pastel.SignAlgorithmED448)
	if err != nil {
		responseWithJSON(writer, http.StatusUnauthorized, map[string]string{"error": "Invalid pid/passphrase"})
		return
	}

	chID := request.URL.Query().Get("challenge_id")
	if chID == "" {
		responseWithJSON(writer, http.StatusBadRequest, map[string]string{"error": "Missing challenge_id query parameter"})
		return
	}

	store, err := local.OpenHistoryDB()
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("Error Opening DB")
		return
	}
	defer store.CloseHistoryDB(ctx)

	metrics, err := store.GetSelfHealingExecutionMetrics(time.Now().AddDate(20, 0, 0))
	if err != nil {
		responseWithJSON(writer, http.StatusNotFound, map[string]string{"error": err.Error()})
		return
	}

	for row := range metrics {
		if metrics[row].ChallengeID == chID {
			responseWithJSON(writer, http.StatusOK, metrics[row])
			return
		}
	}

	responseWithJSON(writer, http.StatusNotFound, map[string]string{"error": fmt.Sprintf("No metrics found for challenge_id %s", chID)})
}

func (service *Service) shTrigger(writer http.ResponseWriter, request *http.Request) {
	ctx := service.contextWithLogPrefix(request.Context())

	// 1. Retrieve Authorization Header (Passphrase)
	passphrase := request.Header.Get("Authorization")

	// 2. Retrieve Pastel ID from Query Param
	pid := request.URL.Query().Get("pid")
	if pid == "" {
		responseWithJSON(writer, http.StatusBadRequest, map[string]string{"error": "Missing pid parameter"})
		return
	}

	_, err := service.scService.PastelClient.Sign(ctx, []byte{1, 2, 3}, pid, passphrase, pastel.SignAlgorithmED448)
	if err != nil {
		responseWithJSON(writer, http.StatusUnauthorized, map[string]string{"error": "Invalid pid/passphrase"})
		return
	}

	triggerID := request.URL.Query().Get("trigger_id")
	if triggerID == "" {
		responseWithJSON(writer, http.StatusBadRequest, map[string]string{"error": "Missing trigger_id query parameter"})
		return
	}

	store, err := local.OpenHistoryDB()
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("Error Opening DB")
		return
	}
	defer store.CloseHistoryDB(ctx)

	metrics, err := store.GetSelfHealingGenerationMetrics(time.Now().AddDate(20, 0, 0))
	if err != nil {
		responseWithJSON(writer, http.StatusNotFound, map[string]string{"error": err.Error()})
		return
	}

	for row := range metrics {
		if metrics[row].TriggerID == triggerID {
			responseWithJSON(writer, http.StatusOK, metrics[row])
			return
		}
	}

	responseWithJSON(writer, http.StatusNotFound, map[string]string{"error": fmt.Sprintf("No metrics found for trigger_id %s", triggerID)})
}

func (service *Service) metrics(writer http.ResponseWriter, request *http.Request) {
	ctx := service.contextWithLogPrefix(request.Context())

	// 1. Retrieve Authorization Header (Passphrase)
	passphrase := request.Header.Get("Authorization")

	// 2. Retrieve Pastel ID from Query Param
	pid := request.URL.Query().Get("pid")
	if pid == "" {
		responseWithJSON(writer, http.StatusBadRequest, map[string]string{"error": "Missing pid parameter"})
		return
	}

	_, err := service.scService.PastelClient.Sign(ctx, []byte{1, 2, 3}, pid, passphrase, pastel.SignAlgorithmED448)
	if err != nil {
		responseWithJSON(writer, http.StatusUnauthorized, map[string]string{"error": "Invalid pid/passphrase"})
		return
	}

	// Rate Limit Check for pid
	if !service.rateLimiter.CheckRateLimit(pid) {
		responseWithJSON(writer, http.StatusTooManyRequests, map[string]string{"error": "Rate limit exceeded, Please try again later."})
		return
	}

	// 3. Get 'From' Time from Query Param
	fromStr := request.URL.Query().Get("start")
	var from time.Time
	if fromStr == "" {
		from = time.Now().UTC().Add(-time.Hour * 24 * 7) // Default to 1 week ago
	} else {
		var err error
		from, err = time.Parse(time.RFC3339, fromStr)
		if err != nil {
			responseWithJSON(writer, http.StatusBadRequest, map[string]string{"error": "Invalid start time format"})
			return
		}
	}

	// 4. Get 'To' Time from Query Param
	toStr := request.URL.Query().Get("end")
	var to *time.Time
	if toStr != "" {
		parsedTo, err := time.Parse(time.RFC3339, toStr)
		if err != nil {
			responseWithJSON(writer, http.StatusBadRequest, map[string]string{"error": "Invalid end time format"})
			return
		}
		to = &parsedTo
	}

	store, err := local.OpenHistoryDB()
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("Error Opening DB")
		return
	}
	defer store.CloseHistoryDB(ctx)

	metrics, err := store.QueryMetrics(ctx, from, to)
	if err != nil {
		responseWithJSON(writer, http.StatusNotFound, map[string]string{"error": err.Error()})
		return
	}

	responseWithJSON(writer, http.StatusOK, metrics)
}

// RateLimiter is a simple rate limiter that limits the number of requests per hour
type RateLimiter struct {
	mutex       sync.Mutex
	requests    map[string][]time.Time
	maxRequests int
	interval    time.Duration
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(maxRequests int, interval time.Duration) *RateLimiter {
	return &RateLimiter{
		requests:    make(map[string][]time.Time),
		maxRequests: maxRequests,
		interval:    interval,
	}
}

// CheckRateLimit checks if the rate limit has been reached for the given pid
func (rl *RateLimiter) CheckRateLimit(pid string) bool {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()

	now := time.Now()

	// Clean up old requests that are out of the interval
	if reqs, found := rl.requests[pid]; found {
		var updatedReqs []time.Time
		for _, t := range reqs {
			if now.Sub(t) < rl.interval {
				updatedReqs = append(updatedReqs, t)
			}
		}
		rl.requests[pid] = updatedReqs
	}

	// Check if the rate limit has been reached
	if len(rl.requests[pid]) >= rl.maxRequests {
		return false
	}

	// Add the new request
	rl.requests[pid] = append(rl.requests[pid], now)
	return true
}
