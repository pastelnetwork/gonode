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
	"github.com/pastelnetwork/gonode/common/utils"
	"github.com/pastelnetwork/gonode/p2p"
	rqnode "github.com/pastelnetwork/gonode/raptorq/node"
)

// StorageHandler provides common logic for RQ and P2P operations
type StorageHandler struct {
	P2PClient p2p.Client
	RqClient  rqnode.ClientInterface

	rqAddress string
	rqDir     string

	TaskID string
	TxID   string
}

// NewStorageHandler creates instance of StorageHandler
func NewStorageHandler(p2p p2p.Client, rq rqnode.ClientInterface,
	rqAddress string, rqDir string) *StorageHandler {

	return &StorageHandler{
		P2PClient: p2p,
		RqClient:  rq,
		rqAddress: rqAddress,
		rqDir:     rqDir,
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

	encodeResp, err := rqService.Encode(ctx, data)
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

// StoreRaptorQSymbolsIntoP2P stores RQ symbols into P2P
func (h *StorageHandler) StoreRaptorQSymbolsIntoP2P(ctx context.Context, data []byte, name string) error {
	symbols, err := h.GenerateRaptorQSymbols(ctx, data, name)
	if err != nil {
		return err
	}

	log.WithContext(ctx).WithField("symbols count", len(symbols)).WithField("task_id", h.TaskID).WithField("reg-txid", h.TxID).
		Info("storing raptorQ symbols in p2p")

	result := make([][]byte, 0, len(symbols))
	for _, value := range symbols {
		result = append(result, value)
	}

	log.WithContext(ctx).WithField("symbols count", len(symbols)).WithField("task_id", h.TaskID).WithField("reg-txid", h.TxID).
		Info("begin batch store raptorQ symbols in p2p")
	if err := h.P2PClient.StoreBatch(ctx, result, P2PDataRaptorQSymbol); err != nil {
		return fmt.Errorf("store batch raptorq symbols in p2p: %w", err)
	}

	log.WithContext(ctx).WithField("symbols count", len(symbols)).WithField("task_id", h.TaskID).WithField("reg-txid", h.TxID).
		Info("done batch stored raptorQ symbols in p2p")

	return nil
}
