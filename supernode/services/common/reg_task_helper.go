package common

import (
	"bytes"
	"context"
	"fmt"
	"sync"
	"time"

	json "github.com/json-iterator/go"
	"github.com/pastelnetwork/gonode/common/blocktracker"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/types"
	"github.com/pastelnetwork/gonode/common/utils"
	"github.com/pastelnetwork/gonode/mixins"
	"github.com/pastelnetwork/gonode/pastel"
)

// RegTaskHelper common operations related to (any) Ticket registration
type RegTaskHelper struct {
	*SuperNodeTask

	NetworkHandler   *NetworkHandler
	ServerPastelID   string
	serverPassPhrase string
	PastelHandler    *mixins.PastelHandler

	ActionTicketRegMetadata    *types.ActionRegMetadata
	preburntTxMinConfirmations int

	// valid only for a task run as primary
	peersTicketSignatureMtx  *sync.Mutex
	PeersTicketSignature     map[string][]byte
	AllSignaturesReceivedChn chan struct{}
}

// NewRegTaskHelper creates instance of RegTaskHelper
func NewRegTaskHelper(task *SuperNodeTask,
	pastelID string, passPhrase string,
	network *NetworkHandler,
	pastelClient pastel.Client,
	preburntTxMinConfirmations int,
) *RegTaskHelper {
	return &RegTaskHelper{
		SuperNodeTask:  task,
		ServerPastelID: pastelID, serverPassPhrase: passPhrase, NetworkHandler: network,
		PastelHandler:              &mixins.PastelHandler{PastelClient: pastelClient},
		preburntTxMinConfirmations: preburntTxMinConfirmations,
		peersTicketSignatureMtx:    &sync.Mutex{},
		PeersTicketSignature:       make(map[string][]byte),
		AllSignaturesReceivedChn:   make(chan struct{}),
	}
}

// AddPeerTicketSignature waits for ticket signatures from other SNs and adds them into internal array
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

// AddPeerCollectionTicketSignature waits for ticket signatures from other SNs and adds them into internal array
func (h *RegTaskHelper) AddPeerCollectionTicketSignature(nodeID string, signature []byte) error {
	h.peersTicketSignatureMtx.Lock()
	defer h.peersTicketSignatureMtx.Unlock()

	var err error

	<-h.NewAction(func(ctx context.Context) error {
		log.WithContext(ctx).Debugf("receive collection ticket signature from node %s", nodeID)
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
//  1. checks signatures
//  2. generates list of 50 IDs and compares them to received
func (h *RegTaskHelper) ValidateIDFiles(ctx context.Context,
	data []byte, ic uint32, max uint32, ids []string, numSignRequired int,
	pastelIDs []string,
	pastelClient pastel.Client,
) ([]byte, [][]byte, error) {

	dec, err := utils.B64Decode(data)
	if err != nil {
		return nil, nil, errors.Errorf("decode data: %w", err)
	}

	decData, err := utils.Decompress(dec)
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
		return nil, nil, errors.Errorf("file verification failed: need %d verifications, got %d", numSignRequired, verifications)
	}

	gotIDs, idFiles, err := pastel.GetIDFiles(ctx, decData, ic, max)
	if err != nil {
		return nil, nil, errors.Errorf("get ids: %w", err)
	}

	if err := utils.EqualStrList(gotIDs, ids); err != nil {
		return nil, nil, errors.Errorf("IDs don't match: %w", err)
	}

	return file, idFiles, nil
}

// WaitConfirmation wait for specific number of confirmations of some blockchain transaction by txid
func (h *RegTaskHelper) WaitConfirmation(ctx context.Context, txid string, minConfirmation int64,
	interval time.Duration, verifyBurnAmt bool, totalAmt float64, percent float64) <-chan error {
	ch := make(chan error)
	log.WithContext(ctx).Info("waiting for txn confirmation")
	go func(ctx context.Context, txid string, ch chan error) {
		defer close(ch)
		blockTracker := blocktracker.New(h.PastelHandler.PastelClient)
		baseBlkCnt, err := blockTracker.GetBlockCount()
		if err != nil {
			log.WithContext(ctx).WithError(err).Error("failed to get block count")
			ch <- err
			return
		}

		for {
			select {
			case <-ctx.Done():
				// context cancelled or abort by caller so no need to return anything
				log.WithContext(ctx).Infof("context done while waiting for confirmation: %s", ctx.Err())
				ch <- ctx.Err()
				return
			case <-time.After(interval):
				gctx, cancel := context.WithCancel(ctx)
				defer cancel()

				go func() {
					<-gctx.Done()
				}()

				txResult, err := h.PastelHandler.PastelClient.GetRawTransactionVerbose1(gctx, txid)
				if err != nil {
					log.WithContext(ctx).WithError(err).WithField("txid", txid).Error("GetRawTransactionVerbose1 err")
				} else {
					if txResult.Confirmations >= minConfirmation {
						if verifyBurnAmt {
							if err := h.verifyTxn(ctx, txResult, totalAmt, percent); err != nil {
								log.WithContext(ctx).WithField("txid", txid).WithError(err).Error("txn verification failed")
								ch <- err
								return
							}
						}

						log.WithContext(ctx).Info("transaction confirmed")
						ch <- nil
						return
					}
				}

				currentBlkCnt, err := blockTracker.GetBlockCount()
				if err != nil {
					log.WithContext(ctx).WithError(err).Error("failed to get block count")
					continue
				}

				if currentBlkCnt-baseBlkCnt >= int32(minConfirmation)+2 {
					ch <- errors.Errorf("timeout when waiting for confirmation of transaction %s", txid)
					return
				}
			}

		}
	}(ctx, txid, ch)

	return ch
}

func (h *RegTaskHelper) verifyTxn(ctx context.Context,
	txn *pastel.GetRawTransactionVerbose1Result, totalAmt float64, percent float64) error {
	inRange := func(val float64, reqVal float64, slackPercent float64) bool {
		lower := reqVal - (reqVal * slackPercent / 100)
		//upper := reqVal + (reqVal * slackPercent / 100)

		return val >= lower
	}

	isTxnAddressOk := false
	reqBurnAmount := totalAmt * percent / 100
	for _, vout := range txn.Vout {
		for _, addr := range vout.ScriptPubKey.Addresses {
			if addr == h.PastelHandler.GetBurnAddress() {
				isTxnAddressOk = true
			}
		}

		if isTxnAddressOk {
			if !inRange(vout.Value, reqBurnAmount, 2.0) {
				data, _ := json.Marshal(txn)
				return fmt.Errorf("invalid transaction amount: %v, required minimum amount : %f - raw transaction details: %s", vout.Value, reqBurnAmount, string(data))
			}

			break
		}
	}

	if !isTxnAddressOk {
		data, _ := json.Marshal(txn)
		return fmt.Errorf("invalid txn address -- correct address: %s - rawTxnData: %s  - total-amount: %f - req-burn-amount: %f",
			h.PastelHandler.GetBurnAddress(), string(data), totalAmt, reqBurnAmount)
	}

	return nil
}

// VerifyPeersTicketSignature verifies ticket signatures of other SNs
func (h *RegTaskHelper) VerifyPeersTicketSignature(ctx context.Context, ticket *pastel.ActionTicket) error {
	log.WithContext(ctx).Debug("all signature received so start validation")

	data, err := pastel.EncodeActionTicket(ticket)
	if err != nil {
		return errors.Errorf("encoded NFT ticket: %w", err)
	}
	return h.VerifyPeersSignature(ctx, data)
}

// VerifyPeersSignature verifies any data signatures of other SNs
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

// VerifyPeersCollectionTicketSignature verifies collection ticket signatures of other SNs
func (h *RegTaskHelper) VerifyPeersCollectionTicketSignature(ctx context.Context, ticket *pastel.CollectionTicket) error {
	log.WithContext(ctx).Debug("all signature received so start validation")

	data, err := pastel.EncodeCollectionTicket(ticket)
	if err != nil {
		return errors.Errorf("encoded NFT ticket: %w", err)
	}
	return h.VerifyCollectionPeersSignature(ctx, data)
}

// VerifyCollectionPeersSignature verifies collection data signatures of other SNs
func (h *RegTaskHelper) VerifyCollectionPeersSignature(ctx context.Context, data []byte) error {
	for nodeID, signature := range h.PeersTicketSignature {
		if ok, err := h.PastelHandler.PastelClient.VerifyCollectionTicket(ctx, data, string(signature), nodeID, pastel.SignAlgorithmED448); err != nil {
			return errors.Errorf("verify signature %s of node %s", signature, nodeID)
		} else if !ok {
			return errors.Errorf("signature of node %s mistmatch", nodeID)
		}
	}
	return nil
}
