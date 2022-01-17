package common

import (
	"context"
	"github.com/pastelnetwork/gonode/common/blocktracker"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/pastel"
	"github.com/pastelnetwork/gonode/walletnode/node"
	"time"
)

type PastelHandler struct {
	PastelClient pastel.Client
}

func NewPastelHandler(pastelClient pastel.Client) *PastelHandler {
	return &PastelHandler{
		PastelClient: pastelClient,
	}
}

// VerifySignature verifies the signature of the data
func (pt *PastelHandler) VerifySignature(ctx context.Context, data []byte, signature string, pastelID string) error {
	ok, err := pt.PastelClient.Verify(ctx, data, signature, pastelID, pastel.SignAlgorithmED448)
	if err != nil {
		return err
	}
	if !ok {
		return errors.Errorf("signature verification failed")
	}

	return nil
}

// GetEstimatedActionFee returns the estimated Action fee for the given image
func (pt *PastelHandler) GetEstimatedActionFee(ctx context.Context, ImgSizeInMb int64) (float64, error) {
	actionFees, err := pt.PastelClient.GetActionFee(ctx, ImgSizeInMb)
	if err != nil {
		return 0, err
	}
	return actionFees.SenseFee, nil
}

// GetTopNodes get list of current top masternodes, validate them and return list of valid (in the order received from pasteld)
func (pt *PastelHandler) GetTopNodes(ctx context.Context,
	nodes node.NodeCollectionInterface,
	nodeClient node.ClientInterface,
	maximumFee float64,
) error {

	mns, err := pt.PastelClient.MasterNodesTop(ctx)
	if err != nil {
		return err
	}
	for _, mn := range mns {
		if mn.Fee > maximumFee {
			continue
		}
		if mn.ExtKey == "" || mn.ExtAddress == "" {
			continue
		}

		// Ensures that the PastelId(mn.ExtKey) of MN node is registered
		_, err = pt.PastelClient.FindTicketByID(ctx, mn.ExtKey)
		if err != nil {
			log.WithContext(ctx).WithField("mn", mn).Warn("FindTicketByID() failed")
			continue
		}
		n := node.NewNode(nodeClient, mn.ExtAddress, mn.ExtKey)
		nodes.Add(n)
	}

	return nil
}

// WaitTxidValid starts goroutine that wait for specified number of confirmations, but stop waiting after time interval
func (pt *PastelHandler) WaitTxidValid(ctx context.Context, txID string, expectedConfirms int64, interval time.Duration) error {
	log.WithContext(ctx).Debugf("Need %d confirmation for txid %s", expectedConfirms, txID)
	blockTracker := blocktracker.New(pt.PastelClient)
	baseBlkCnt, err := blockTracker.GetBlockCount()
	if err != nil {
		log.WithContext(ctx).WithError(err).Warn("failed to get block count")
		return err
	}

	for {
		select {
		case <-ctx.Done():
			return errors.Errorf("context done: %w", ctx.Err())
		case <-time.After(interval):
			checkConfirms := func() error {
				subCtx, cancel := context.WithTimeout(ctx, 15*time.Second)
				defer cancel()

				result, err := pt.PastelClient.GetRawTransactionVerbose1(subCtx, txID)
				if err != nil {
					return errors.Errorf("get transaction: %w", err)
				}

				if result.Confirmations >= expectedConfirms {
					return nil
				}

				return errors.Errorf("not enough confirmations: expected %d, got %d", expectedConfirms, result.Confirmations)
			}

			err := checkConfirms()
			if err != nil {
				log.WithContext(ctx).WithError(err).Warn("check confirmations failed")
			} else {
				return nil
			}

			currentBlkCnt, err := blockTracker.GetBlockCount()
			if err != nil {
				log.WithContext(ctx).WithError(err).Warn("failed to get block count")
				continue
			}

			if currentBlkCnt-baseBlkCnt >= int32(expectedConfirms)+2 {
				return errors.Errorf("timeout when wating for confirmation of transaction %s", txID)
			}
		}
	}
}
