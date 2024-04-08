package common

import (
	"context"
	"fmt"
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
	loadSymbolsBatchSize = 2000
)

// StorageHandler provides common logic for RQ and P2P operations
type StorageHandler struct {
	P2PClient p2p.Client
	RqClient  rqnode.ClientInterface

	rqAddress string
	rqDir     string

	TaskID string
	TxID   string

	store rqstore.Store
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

	return h.P2PClient.StoreBatch(ctx, list, typ)
}

// GenerateRaptorQSymbols calls RQ service to produce RQ Symbols
func (h *StorageHandler) GenerateRaptorQSymbols(ctx context.Context, data []byte, name string) (map[string][]byte, error) {
	if h.RqClient == nil {
		log.WithContext(ctx).Warnf("RQ Server is not initialized")
		return nil, errors.Errorf("RQ Server is not initialized")
	}

	b := backoff.NewExponentialBackOff()
	b.MaxElapsedTime = 2 * time.Minute
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

	encodeResp, err := rqService.RQEncode(ctx, data, h.TxID, h.store)
	if err != nil {
		return nil, errors.Errorf("create raptorq symbol from data %s: %w", name, err)
	}

	return encodeResp.Symbols, nil
}

// GetRaptorQEncodeInfo calls RQ service to get Encoding info and list of RQIDs
func (h *StorageHandler) GetRaptorQEncodeInfo(ctx context.Context,
	data []byte, num uint32, hash string, pastelID string,
) (*rqnode.EncodeInfo, error) {
	if h.RqClient == nil {
		log.WithContext(ctx).Warnf("RQ Server is not initialized")
		return nil, errors.Errorf("RQ Server is not initialized")
	}

	b := backoff.NewExponentialBackOff()
	b.MaxElapsedTime = 2 * time.Minute
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

	encodeInfo, err := rqService.EncodeInfo(ctx, data, num, hash, pastelID)
	if err != nil {
		return nil, errors.Errorf("get raptorq encode info: %w", err)
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

func (h *StorageHandler) StoreRaptorQSymbolsIntoP2P(ctx context.Context, data []byte, name string) error {
	// Generate the keys for RaptorQ symbols, with empty values
	keysMap, err := h.GenerateRaptorQSymbols(ctx, data, name)
	if err != nil {
		return err
	}

	processBatch := func(batchKeys map[string][]byte) error {
		// Load symbols from the database for the current batch
		loadedSymbols, err := h.store.LoadSymbols(h.TxID, batchKeys)
		if err != nil {
			return fmt.Errorf("load batch symbols from db: %w", err)
		}

		// Prepare batch for P2P storage
		result := make([][]byte, len(loadedSymbols))
		i := 0
		for key, value := range loadedSymbols {
			result[i] = value
			loadedSymbols[key] = nil // Release the reference for faster memory cleanup
			i++
		}

		// Store the loaded symbols in P2P
		if err := h.P2PClient.StoreBatch(ctx, result, P2PDataRaptorQSymbol); err != nil {
			return fmt.Errorf("store batch raptorq symbols in p2p: %w", err)
		}

		return nil
	}

	// Iterate over keysMap in batches
	batchSize := loadSymbolsBatchSize
	batchKeys := make(map[string][]byte)
	count := 0

	for key := range keysMap {
		batchKeys[key] = nil
		count++

		if count%batchSize == 0 || count == len(keysMap) {
			if err := processBatch(batchKeys); err != nil {
				return err
			}
			// Reset batchKeys for the next batch
			batchKeys = make(map[string][]byte)
		}
	}

	return nil
}
