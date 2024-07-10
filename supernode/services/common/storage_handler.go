package common

import (
	"context"
	"fmt"
	"math"
	"sort"
	"time"

	json "github.com/json-iterator/go"

	"github.com/cenkalti/backoff"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/storage/files"
	"github.com/pastelnetwork/gonode/common/storage/rqstore"
	"github.com/pastelnetwork/gonode/common/utils"
	"github.com/pastelnetwork/gonode/p2p"
	rqnode "github.com/pastelnetwork/gonode/raptorq/node"
)

const (
	loadSymbolsBatchSize = 2500
	storeSymbolsPercent  = 10
	concurrency          = 2
)

// StorageHandler provides common logic for RQ and P2P operations
type StorageHandler struct {
	P2PClient p2p.Client
	RqClient  rqnode.ClientInterface

	rqAddress string
	rqDir     string

	TaskID string
	TxID   string

	store     rqstore.Store
	semaphore chan struct{}
}

// NewStorageHandler creates instance of StorageHandler
func NewStorageHandler(p2p p2p.Client, rq rqnode.ClientInterface,
	rqAddress string, rqDir string, store rqstore.Store) *StorageHandler {

	return &StorageHandler{
		P2PClient: p2p,
		RqClient:  rq,
		rqAddress: rqAddress,
		rqDir:     rqDir,
		store:     store,
		semaphore: make(chan struct{}, concurrency),
	}
}

// StoreFileIntoP2P stores file into P2P
func (h *StorageHandler) StoreFileIntoP2P(ctx context.Context, file *files.File, typ int) (string, error) {
	data, err := file.Bytes()
	if err != nil {
		return "", errors.Errorf("store file %s into p2p", file.Name())
	}
	return h.StoreBytesIntoP2P(ctx, data, typ)
}

// StoreBytesIntoP2P into P2P actual data
func (h *StorageHandler) StoreBytesIntoP2P(ctx context.Context, data []byte, typ int) (string, error) {
	return h.P2PClient.Store(ctx, data, typ)
}

// StoreBatch stores into P2P array of bytes arrays
func (h *StorageHandler) StoreBatch(ctx context.Context, list [][]byte, typ int) error {
	val := ctx.Value(log.TaskIDKey)
	taskID := ""
	if val != nil {
		taskID = fmt.Sprintf("%v", val)
	}
	log.WithContext(ctx).WithField("task_id", taskID).Info("task_id in storeList")

	return h.P2PClient.StoreBatch(ctx, list, typ, taskID)
}

// GenerateRaptorQSymbols calls RQ service to produce RQ Symbols
func (h *StorageHandler) GenerateRaptorQSymbols(ctx context.Context, data []byte, name string) (map[string][]byte, error) {
	if h.RqClient == nil {
		log.WithContext(ctx).Warnf("RQ Server is not initialized")
		return nil, errors.Errorf("RQ Server is not initialized")
	}

	b := backoff.NewExponentialBackOff()
	b.MaxElapsedTime = 3 * time.Minute
	b.InitialInterval = 200 * time.Millisecond

	var conn rqnode.Connection
	if err := backoff.Retry(backoff.Operation(func() error {
		var err error
		conn, err = h.RqClient.Connect(ctx, h.rqAddress)
		if err != nil {
			return errors.Errorf("connect to raptorq service: %w", err)
		}

		return nil
	}), b); err != nil {
		return nil, fmt.Errorf("retry connect to raptorq service: %w", err)
	}
	defer func() {
		if err := conn.Close(); err != nil {
			log.WithContext(ctx).WithError(err).Error("error closing rq-connection")
		}
	}()

	rqService := conn.RaptorQ(&rqnode.Config{
		RqFilesDir: h.rqDir,
	})

	b.Reset()

	encodeResp := &rqnode.Encode{}
	if err := backoff.Retry(backoff.Operation(func() error {
		var err error
		encodeResp, err = rqService.RQEncode(ctx, data, h.TxID, h.store)
		if err != nil {
			return errors.Errorf("create raptorq symbol from data %s: %w", name, err)
		}

		return nil
	}), b); err != nil {
		return nil, fmt.Errorf("retry do rqencode service: %w", err)
	}

	return encodeResp.Symbols, nil
}

// GetRaptorQEncodeInfo calls RQ service to get Encoding info and list of RQIDs
func (h *StorageHandler) GetRaptorQEncodeInfo(ctx context.Context,
	data []byte, num uint32, hash string, pastelID string,
) (encodeInfo *rqnode.EncodeInfo, err error) {
	if h.RqClient == nil {
		log.WithContext(ctx).Warnf("RQ Server is not initialized")
		return nil, errors.Errorf("RQ Server is not initialized")
	}

	b := backoff.NewExponentialBackOff()
	b.MaxElapsedTime = 3 * time.Minute
	b.InitialInterval = 500 * time.Millisecond

	var conn rqnode.Connection
	if err := backoff.Retry(backoff.Operation(func() error {
		var err error
		conn, err = h.RqClient.Connect(ctx, h.rqAddress)
		if err != nil {
			return errors.Errorf("connect to raptorq service: %w", err)
		}

		return nil
	}), b); err != nil {
		return nil, fmt.Errorf("retry connect to raptorq service: %w", err)
	}
	defer func() {
		if err := conn.Close(); err != nil {
			log.WithContext(ctx).WithError(err).Error("error closing rq-connection")
		}
	}()

	rqService := conn.RaptorQ(&rqnode.Config{
		RqFilesDir: h.rqDir,
	})

	b.Reset()
	if err := backoff.Retry(backoff.Operation(func() error {
		var err error
		encodeInfo, err = rqService.EncodeInfo(ctx, data, num, hash, pastelID)
		if err != nil {
			return errors.Errorf("get raptorq encode info: %w", err)
		}
		return nil
	}), b); err != nil {
		return nil, fmt.Errorf("retry do encode info on raptorq service: %w", err)
	}

	return encodeInfo, nil
}

// ValidateRaptorQSymbolIDs calls RQ service to get Encoding info and list of RQIDs and compares them to the similar data received from WN
func (h *StorageHandler) ValidateRaptorQSymbolIDs(ctx context.Context,
	data []byte, num uint32, hash string, pastelID string,
	haveData []byte) error {

	if len(haveData) == 0 {
		return errors.Errorf("no symbols identifiers")
	}

	encodeInfo, err := h.GetRaptorQEncodeInfo(ctx, data, num, hash, pastelID)
	if err != nil {
		return err
	}

	// pick just one file generated to compare
	var gotFile, haveFile rqnode.RawSymbolIDFile
	for _, v := range encodeInfo.SymbolIDFiles {
		gotFile = v
		break
	}

	if err := json.Unmarshal(haveData, &haveFile); err != nil {
		return errors.Errorf("decode raw rq file: %w", err)
	}

	if err := utils.EqualStrList(gotFile.SymbolIdentifiers, haveFile.SymbolIdentifiers); err != nil {
		return errors.Errorf("raptor symbol mismatched: %w", err)
	}
	return nil
}

func (h *StorageHandler) storeSymbolsInP2P(ctx context.Context, dir string, batchKeys map[string][]byte) error {
	val := ctx.Value(log.TaskIDKey)
	taskID := ""
	if val != nil {
		taskID = fmt.Sprintf("%v", val)
	}
	// Load symbols from the database for the current batch
	log.WithContext(ctx).WithField("count", len(batchKeys)).Info("loading batch symbols")
	loadedSymbols, err := utils.LoadSymbols(dir, batchKeys)
	if err != nil {
		return fmt.Errorf("load batch symbols from db: %w", err)
	}

	log.WithContext(ctx).WithField("count", len(loadedSymbols)).Info("loaded batch symbols, storing now")
	// Prepare batch for P2P storage		return nil
	result := make([][]byte, len(loadedSymbols))
	i := 0
	for key, value := range loadedSymbols {
		result[i] = value
		loadedSymbols[key] = nil // Release the reference for faster memory cleanup
		i++
	}

	// Store the loaded symbols in P2P
	if err := h.P2PClient.StoreBatch(ctx, result, P2PDataRaptorQSymbol, taskID); err != nil {
		return fmt.Errorf("store batch raptorq symbols in p2p: %w", err)
	}
	log.WithContext(ctx).WithField("count", len(loadedSymbols)).Info("stored batch symbols")

	if err := utils.DeleteSymbols(ctx, dir, batchKeys); err != nil {
		return fmt.Errorf("delete batch symbols from db: %w", err)
	}
	log.WithContext(ctx).WithField("count", len(loadedSymbols)).Info("deleted batch symbols")

	return nil
}

func (h *StorageHandler) StoreRaptorQSymbolsIntoP2P(ctx context.Context, data []byte, name string) error {
	h.semaphore <- struct{}{} // Acquire slot
	defer func() {
		<-h.semaphore // Release the semaphore slot
	}()

	// Generate the keys for RaptorQ symbols, with empty values
	log.WithContext(ctx).Info("generating RaptorQ symbols")
	keysMap, err := h.GenerateRaptorQSymbols(ctx, data, name)
	if err != nil {
		return err
	}
	log.WithContext(ctx).WithField("count", len(keysMap)).Info("generated RaptorQ symbols")

	if h.TxID == "" {
		return errors.New("txid is not set, cannot store rq symbols")
	}

	dir, err := h.store.GetDirectoryByTxID(h.TxID)
	if err != nil {
		return fmt.Errorf("error fetching symbols dir from rq DB: %w", err)
	}

	// Create a slice of keys from keysMap and sort it
	keys := make([]string, 0, len(keysMap))
	for key := range keysMap {
		keys = append(keys, key)
	}
	sort.Strings(keys) // Sort the keys alphabetically

	if len(keys) > loadSymbolsBatchSize {
		// Calculate 15% of the total keys, rounded up
		requiredKeysCount := int(math.Ceil(float64(len(keys)) * storeSymbolsPercent / 100))

		// Get the subset of keys (15%)
		if requiredKeysCount > len(keys) {
			requiredKeysCount = len(keys) // Ensure we don't exceed the available keys count
		}
		keys = keys[:requiredKeysCount]
	}

	// Iterate over sorted keys in batches
	batchKeys := make(map[string][]byte)
	count := 0

	log.WithContext(ctx).WithField("count", len(keys)).Info("storing raptorQ symbols")
	for _, key := range keys {
		batchKeys[key] = nil
		count++
		if count%loadSymbolsBatchSize == 0 {
			if err := h.storeSymbolsInP2P(ctx, dir, batchKeys); err != nil {
				return err
			}
			batchKeys = make(map[string][]byte) // Reset batchKeys after storing
		}
	}

	// Store any remaining symbols in the last batch
	if len(batchKeys) > 0 {
		if err := h.storeSymbolsInP2P(ctx, dir, batchKeys); err != nil {
			return err
		}
	}

	if err := h.store.UpdateIsFirstBatchStored(h.TxID); err != nil {
		return fmt.Errorf("error updating first batch stored flag in rq DB: %w", err)
	}
	log.WithContext(ctx).WithField("curr-time", time.Now().UTC()).WithField("count", len(keys)).Info("stored RaptorQ symbols")

	return nil
}
