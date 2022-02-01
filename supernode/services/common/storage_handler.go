package common

import (
	"context"
	"encoding/json"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/storage/files"
	"github.com/pastelnetwork/gonode/common/utils"
	"github.com/pastelnetwork/gonode/p2p"
	rqnode "github.com/pastelnetwork/gonode/raptorq/node"
)

type StorageHandler struct {
	P2PClient p2p.Client
	RqClient  rqnode.ClientInterface

	rqAddress string
	rqDir     string
}

func NewStorageHandler(p2p p2p.Client, rq rqnode.ClientInterface,
	rqAddress string, rqDir string,
) *StorageHandler {
	return &StorageHandler{
		P2PClient: p2p,
		RqClient:  rq,
		rqAddress: rqAddress, rqDir: rqDir,
	}
}

func (h *StorageHandler) StoreFileIntoP2P(ctx context.Context, nft *files.File) (string, error) {
	data, err := nft.Bytes()
	if err != nil {
		return "", errors.Errorf("store file %s into p2p", nft.Name())
	}
	return h.StoreBytesIntoP2P(ctx, data)
}

func (h *StorageHandler) StoreBytesIntoP2P(ctx context.Context, data []byte) (string, error) {
	return h.P2PClient.Store(ctx, data)
}

func (h *StorageHandler) StoreListOfBytesIntoP2P(ctx context.Context, list [][]byte) error {
	for _, data := range list {
		if _, err := h.StoreBytesIntoP2P(ctx, data); err != nil {
			return errors.Errorf("store data into p2p: %w", err)
		}
	}
	return nil
}

func (h *StorageHandler) GenerateRaptorQSymbols(ctx context.Context, data []byte, name string) (map[string][]byte, error) {
	if h.RqClient == nil {
		log.WithContext(ctx).Warnf("RQ Server is not initialized")
		return nil, errors.Errorf("RQ Server is not initialized")
	}

	conn, err := h.RqClient.Connect(ctx, h.rqAddress)
	if err != nil {
		return nil, errors.Errorf("connect to raptorq service: %w", err)
	}
	defer conn.Close()

	rqService := conn.RaptorQ(&rqnode.Config{
		RqFilesDir: h.rqDir,
	})

	encodeResp, err := rqService.Encode(ctx, data)
	if err != nil {
		return nil, errors.Errorf("create raptorq symbol from data %s: %w", name, err)
	}

	return encodeResp.Symbols, nil
}

func (h *StorageHandler) GetRaptorQEncodeInfo(ctx context.Context,
	data []byte, num uint32, hash string, pastelID string,
) (*rqnode.EncodeInfo, error) {
	if h.RqClient == nil {
		log.WithContext(ctx).Warnf("RQ Server is not initialized")
		return nil, errors.Errorf("RQ Server is not initialized")
	}

	conn, err := h.RqClient.Connect(ctx, h.rqAddress)
	if err != nil {
		return nil, errors.Errorf("connect to raptorq service: %w", err)
	}
	defer conn.Close()

	rqService := conn.RaptorQ(&rqnode.Config{
		RqFilesDir: h.rqDir,
	})

	encodeInfo, err := rqService.EncodeInfo(ctx, data, num, hash, pastelID)
	if err != nil {
		return nil, errors.Errorf("get raptorq encode info: %w", err)
	}

	return encodeInfo, nil
}

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
	symbols, err := h.GenerateRaptorQSymbols(ctx, data, name)
	if err != nil {
		return err
	}

	for id, symbol := range symbols {
		if _, err := h.P2PClient.Store(ctx, symbol); err != nil {
			return errors.Errorf("store symbol id=%s into kademlia: %w", id, err)
		}
	}

	return nil
}
