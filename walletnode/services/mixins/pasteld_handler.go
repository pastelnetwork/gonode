package mixins

import (
	"context"
	"github.com/pastelnetwork/gonode/common/blocktracker"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/pastel"
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
func (pt *PastelHandler) VerifySignature(ctx context.Context, data []byte, signature string, pastelID string, algo string) (bool, error) {
	ok, err := pt.PastelClient.Verify(ctx, data, signature, pastelID, algo)
	if err != nil {
		return false, err
	}
	if !ok {
		return false, errors.Errorf("signature verification failed")
	}

	return true, nil
}

// GetEstimatedActionFee returns the estimated Action fee for the given image
func (pt *PastelHandler) GetEstimatedActionFee(ctx context.Context, ImgSizeInMb int64) (float64, error) {
	actionFees, err := pt.PastelClient.GetActionFee(ctx, ImgSizeInMb)
	if err != nil {
		return 0, err
	}
	return actionFees.SenseFee, nil
}

// determine current block height & hash of it
func (pt *PastelHandler) GetBlock(ctx context.Context) (int, string, error) {
	// Get block num
	blockNum, err := pt.PastelClient.GetBlockCount(ctx)
	if err != nil {
		return -1, "", errors.Errorf("get block num: %w", err)
	}

	// Get block hash string
	blockInfo, err := pt.PastelClient.GetBlockVerbose1(ctx, blockNum)
	if err != nil {
		return -1, "", errors.Errorf("get block info blocknum=%d: %w", blockNum, err)
	}

	// Decode hash string to byte
	if err != nil {
		return -1, "", errors.Errorf("convert hash string %s to bytes: %w", blockInfo.Hash, err)
	}

	return int(blockNum), blockInfo.Hash, nil
}

func (pt *PastelHandler) CheckRegistrationFee(ctx context.Context, address string, fee float64, maxFee float64) error {
	if fee == 0 {
		return errors.Errorf("invalid fee amount to check: %f", fee)
	}

	if fee > maxFee {
		return errors.Errorf("registration fee is to expensive - maximum-fee (%f) < registration-fee(%f)", maxFee, fee)
	}

	balance, err := pt.PastelClient.GetBalance(ctx, address)
	if err != nil {
		return errors.Errorf("get balance of address(%s): %w", address, err)
	}

	if balance < fee {
		return errors.Errorf("not enough PSL - balance(%f) < registration-fee(%f)", balance, fee)
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
