package mixins

import (
	"context"
	"time"

	"github.com/pastelnetwork/gonode/common/blocktracker"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/pastel"
)

const (
	defaultDownloadTimeout = 300 * time.Second
)

// PastelHandler handles pastel communication
type PastelHandler struct {
	PastelClient pastel.Client
}

func NewPastelHandler(pastelClient pastel.Client) *PastelHandler {
	return &PastelHandler{
		PastelClient: pastelClient,
	}
}

// TicketInfo contains information about the ticket
type TicketInfo struct {
	IsTicketPublic        bool
	EstimatedDownloadTime time.Duration
	Filename              string
	FileType              string
	DataHash              []byte
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

// GetEstimatedSenseFee returns the estimated Action fee for the given image
func (pt *PastelHandler) GetEstimatedSenseFee(ctx context.Context, ImgSizeInMb float64) (float64, error) {
	actionFees, err := pt.PastelClient.GetActionFee(ctx, int64(ImgSizeInMb))
	if err != nil {
		return 0, err
	}

	// adding 10.0 PSL because activation ticket takes 10 PSL
	return actionFees.SenseFee, nil
}

// GetEstimatedSenseFee returns the estimated Action fee for the given image
func (pt *PastelHandler) GetEstimatedCascadeFee(ctx context.Context, ImgSizeInMb float64) (float64, error) {
	actionFees, err := pt.PastelClient.GetActionFee(ctx, int64(ImgSizeInMb))
	if err != nil {
		return 0, err
	}

	// adding 10.0 PSL because activation ticket takes 10 PSL
	return actionFees.CascadeFee, nil
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

// IsBalanceMoreThanMaxFee checks balance of provided address to be more than max fee
func (pt *PastelHandler) IsBalanceMoreThanMaxFee(ctx context.Context, address string, maxFee float64) error {
	balance, err := pt.PastelClient.GetBalance(ctx, address)
	if err != nil {
		return errors.Errorf("get balance of address(%s): %w", address, err)
	}

	if balance < maxFee {
		return errors.Errorf("current balance is less than provided max Fee - balance(%f) < max-fee(%f)", balance, maxFee)
	}

	return nil
}

// CheckBalanceToPayRegistrationFee checks balance
func (pt *PastelHandler) CheckBalanceToPayRegistrationFee(ctx context.Context, address string, fee float64, maxFee float64) error {
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

// BurnSomeCoins burns coins on the address
func (pt *PastelHandler) BurnSomeCoins(ctx context.Context, address string, amount int64, percentageOfAmountToBurn uint) (string, error) {
	if amount <= 0 {
		return "", errors.Errorf("invalid amount")
	}

	burnAmount := float64(amount) / float64(percentageOfAmountToBurn)
	burnTxid, err := pt.PastelClient.SendFromAddress(ctx, address, pt.PastelClient.BurnAddress(), burnAmount)
	if err != nil {
		return "", errors.Errorf("burn %d percent of amount: %w", percentageOfAmountToBurn, err)
	}
	log.WithContext(ctx).Debugf("burn coins txid: %s", burnTxid)
	return burnTxid, nil
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
			return errors.Errorf("wait tixid valid context done: %w", ctx.Err())
		case <-time.After(interval):
			checkConfirms := func() error {
				subCtx, cancel := context.WithTimeout(ctx, 15*time.Second)
				defer cancel()

				result, err := pt.PastelClient.GetRawTransactionVerbose1(subCtx, txID)
				if err != nil {
					log.WithContext(ctx).WithField("txid", txID).Warn("failed to get txn in mixins")
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
				return errors.Errorf("timeout when waiting for confirmation of transaction %s", txID)
			}
		}
	}
}

// RegTicket pull NFT registration ticket from cNode & decodes base64 encoded fields
func (pt *PastelHandler) RegTicket(ctx context.Context, RegTXID string) (*pastel.RegTicket, error) {
	regTicket, err := pt.PastelClient.RegTicket(ctx, RegTXID)
	if err != nil {
		return nil, errors.Errorf("fetch: %w", err)
	}

	articketData, err := pastel.DecodeNFTTicket(regTicket.RegTicketData.NFTTicket)
	if err != nil {
		return nil, errors.Errorf("convert NFT ticket: %w", err)
	}

	regTicket.RegTicketData.NFTTicketData = *articketData

	return &regTicket, nil
}

func (pt *PastelHandler) GetTicketInfo(ctx context.Context, txid, ttype string) (info TicketInfo, err error) {
	info.EstimatedDownloadTime = defaultDownloadTimeout
	switch ttype {
	case pastel.ActionTypeSense:
		// Intent here is to just verify that the txid is valid, Sense file is always set to public
		_, err := pt.PastelClient.ActionRegTicket(ctx, txid)
		if err != nil {
			log.WithContext(ctx).WithError(err).Error("could not get action registered ticket")
			return info, err
		}

		info.IsTicketPublic = false
	case pastel.ActionTypeCascade:
		ticket, err := pt.PastelClient.ActionRegTicket(ctx, txid)
		if err != nil {
			log.WithContext(ctx).WithError(err).Error("could not get action registered ticket")
			return info, err
		}

		actionTicket, err := pastel.DecodeActionTicket([]byte(ticket.ActionTicketData.ActionTicket))
		if err != nil {
			log.WithContext(ctx).WithError(err).Error("cloud not decode action ticket: %w", err)
			return info, err
		}
		ticket.ActionTicketData.ActionTicketData = *actionTicket

		cTicket, err := ticket.ActionTicketData.ActionTicketData.APICascadeTicket()
		if err != nil {
			log.WithContext(ctx).WithError(err).Error("could not get registered ticket")
			return info, err
		}

		info.DataHash = cTicket.DataHash
		info.IsTicketPublic = cTicket.MakePubliclyAccessible
		info.Filename = cTicket.FileName
		info.FileType = cTicket.FileType

		est := getEstimatedDownloadSizeOnBytes(cTicket.OriginalFileSizeInBytes)
		if est > defaultDownloadTimeout {
			info.EstimatedDownloadTime = est
		}

	default:
		regTicket, err := pt.RegTicket(ctx, txid)
		if err == nil {
			info.IsTicketPublic = regTicket.RegTicketData.NFTTicketData.AppTicketData.MakePubliclyAccessible
			info.Filename = regTicket.RegTicketData.NFTTicketData.AppTicketData.FileName
			info.FileType = regTicket.RegTicketData.NFTTicketData.AppTicketData.FileType
			info.DataHash = regTicket.RegTicketData.NFTTicketData.AppTicketData.DataHash

			est := getEstimatedDownloadSizeOnBytes(regTicket.RegTicketData.NFTTicketData.AppTicketData.OriginalFileSizeInBytes)
			if est > defaultDownloadTimeout {
				info.EstimatedDownloadTime = est
			}
		} else {
			log.WithContext(ctx).WithError(err).Error("could not get registered ticket")
			return info, err
		}
	}

	return info, nil
}

func getEstimatedDownloadSizeOnBytes(size int) time.Duration {
	return time.Duration((size / 50000)) * (time.Millisecond * 800)
}

func (pt *PastelHandler) GetBurnAddress() string {
	return pt.PastelClient.BurnAddress()
}
