package common

import (
	"bytes"
	"context"
	"github.com/DataDog/zstd"
	"github.com/pastelnetwork/gonode/common/blocktracker"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/utils"
	"github.com/pastelnetwork/gonode/mixins"
	"github.com/pastelnetwork/gonode/pastel"
	"sync"
	"time"
)

type RegTaskHelper struct {
	*SuperNodeTask

	NetworkHandler   *NetworkHandler
	ServerPastelID   string
	serverPassPhrase string
	PastelHandler    *mixins.PastelHandler

	// valid only for a task run as primary
	peersTicketSignatureMtx  *sync.Mutex
	PeersTicketSignature     map[string][]byte
	AllSignaturesReceivedChn chan struct{}
}

func NewRegTaskHelper(task *SuperNodeTask,
	pastelID string, passPhrase string,
	network *NetworkHandler,
	pastelClient pastel.Client,
) *RegTaskHelper {
	return &RegTaskHelper{
		SuperNodeTask:  task,
		ServerPastelID: pastelID, serverPassPhrase: passPhrase, NetworkHandler: network,
		PastelHandler:            &mixins.PastelHandler{PastelClient: pastelClient},
		peersTicketSignatureMtx:  &sync.Mutex{},
		PeersTicketSignature:     make(map[string][]byte),
		AllSignaturesReceivedChn: make(chan struct{}),
	}
}

func (h *RegTaskHelper) AddPeerTicketSignature(nodeID string, signature []byte, reqStatus Status) error {
	h.peersTicketSignatureMtx.Lock()
	defer h.peersTicketSignatureMtx.Unlock()

	if err := h.RequiredStatus(reqStatus); err != nil {
		return err
	}

	var err error

	<-h.NewAction(func(ctx context.Context) error {
		log.WithContext(ctx).Debugf("receive NFT ticket signature from node %s", nodeID)
		if node := h.NetworkHandler.Accepted.ByID(nodeID); node == nil {
			log.WithContext(ctx).WithField("node", nodeID).Errorf("node is not in Accepted list")
			err = errors.Errorf("node %s not in Accepted list", nodeID)
			return nil
		}

		h.PeersTicketSignature[nodeID] = signature
		if len(h.PeersTicketSignature) == len(h.NetworkHandler.Accepted) {
			log.WithContext(ctx).Debug("all signature received")
			go func() {
				close(h.AllSignaturesReceivedChn)
			}()
		}
		return nil
	})
	return err
}

// ValidateIDFiles validates received (IDs) file and its (50) IDs:
// 	1. checks signatures
//	2. generates list of 50 IDs and compares them to received
func (h *RegTaskHelper) ValidateIDFiles(ctx context.Context,
	data []byte, ic uint32, max uint32, ids []string, numSignRequired int,
	pastelIDs []string,
	pastelClient pastel.Client,
) ([]byte, [][]byte, error) {

	dec, err := utils.B64Decode(data)
	if err != nil {
		return nil, nil, errors.Errorf("decode data: %w", err)
	}

	decData, err := zstd.Decompress(nil, dec)
	if err != nil {
		return nil, nil, errors.Errorf("decompress: %w", err)
	}

	splits := bytes.Split(decData, []byte{pastel.SeparatorByte})
	if len(splits) != numSignRequired+1 {
		return nil, nil, errors.New("invalid data")
	}

	file, err := utils.B64Decode(splits[0])
	if err != nil {
		return nil, nil, errors.Errorf("decode file: %w", err)
	}

	verifications := 0
	verifiedNodes := make(map[int]bool)
	for i := 1; i < numSignRequired+1; i++ {
		for j := 0; j < len(pastelIDs); j++ {
			if _, ok := verifiedNodes[j]; ok {
				continue
			}

			verified, err := pastelClient.Verify(ctx, file, string(splits[i]), pastelIDs[j], pastel.SignAlgorithmED448)
			if err != nil {
				return nil, nil, errors.Errorf("verify file signature %w", err)
			}

			if verified {
				verifiedNodes[j] = true
				verifications++
				break
			}
		}
	}

	if verifications != numSignRequired {
		return nil, nil, errors.Errorf("file verification failed: need %d verifications, got %d", 3, verifications)
	}

	gotIDs, idFiles, err := pastel.GetIDFiles(decData, ic, max)
	if err != nil {
		return nil, nil, errors.Errorf("get ids: %w", err)
	}

	if err := utils.EqualStrList(gotIDs, ids); err != nil {
		return nil, nil, errors.Errorf("IDs don't match: %w", err)
	}

	return file, idFiles, nil
}

func (h *RegTaskHelper) VerifyPeersSignature(ctx context.Context, data []byte) error {
	for nodeID, signature := range h.PeersTicketSignature {
		if ok, err := h.PastelHandler.PastelClient.Verify(ctx, data, string(signature), nodeID, pastel.SignAlgorithmED448); err != nil {
			return errors.Errorf("verify signature %s of node %s", signature, nodeID)
		} else if !ok {
			return errors.Errorf("signature of node %s mistmatch", nodeID)
		}
	}
	return nil
}

func (h *RegTaskHelper) WaitConfirmation(ctx context.Context, txid string, minConfirmation int64, interval time.Duration) <-chan error {
	ch := make(chan error)

	go func(ctx context.Context, txid string) {
		defer close(ch)
		blockTracker := blocktracker.New(h.PastelHandler.PastelClient)
		baseBlkCnt, err := blockTracker.GetBlockCount()
		if err != nil {
			log.WithContext(ctx).WithError(err).Warn("failed to get block count")
			ch <- err
			return
		}

		for {
			select {
			case <-ctx.Done():
				// context cancelled or abort by caller so no need to return anything
				log.WithContext(ctx).Debugf("context done: %s", ctx.Err())
				ch <- ctx.Err()
				return
			case <-time.After(interval):
				txResult, err := h.PastelHandler.PastelClient.GetRawTransactionVerbose1(ctx, txid)
				if err != nil {
					log.WithContext(ctx).WithError(err).Warn("GetRawTransactionVerbose1 err")
				} else {
					if txResult.Confirmations >= minConfirmation {
						log.WithContext(ctx).Debug("transaction confirmed")
						ch <- nil
						return
					}
				}

				currentBlkCnt, err := blockTracker.GetBlockCount()
				if err != nil {
					log.WithContext(ctx).WithError(err).Warn("failed to get block count")
					continue
				}

				if currentBlkCnt-baseBlkCnt >= int32(minConfirmation)+2 {
					ch <- errors.Errorf("timeout when wating for confirmation of transaction %s", txid)
					return
				}
			}

		}
	}(ctx, txid)
	return ch
}
