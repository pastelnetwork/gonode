package duplicate

import (
	"bytes"
	"context"
	"fmt"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/pastel"
)

var (
	// ErrDuplicateRegistration error when duplicate ticket registration is in progress
	ErrDuplicateRegistration = errors.Errorf("duplicate ticket registration is already in progress")
)

const (
	//BlocksThresholdForDuplicateTicket is the threshold to check if the duplicate ticket is created before 10000 blocks
	BlocksThresholdForDuplicateTicket = 10000
)

// DupTicketsDetector struct implemented the funcs used to detect duplicate tickets
type DupTicketsDetector struct {
	pastelClient pastel.Client
}

// NewDupTicketsDetector is the constructor to initialize dup tickets detector
func NewDupTicketsDetector(client pastel.Client) *DupTicketsDetector {
	return &DupTicketsDetector{
		pastelClient: client,
	}
}

// DetectDuplicates represents an interface that checks mempool and inactive transactions for duplicate tickets
type DetectDuplicates interface {
	CheckDuplicateCascadeTickets(ctx context.Context, dataHash []byte) error
	CheckDuplicateSenseOrNFTTickets(ctx context.Context, dataHash []byte) error
}

// CheckDuplicateCascadeTickets checks for duplicate cascade tickets in mempool and inactive transactions list
func (task *DupTicketsDetector) CheckDuplicateCascadeTickets(ctx context.Context, dataHash []byte) error {
	if err := task.checkMempoolForInProgressCascadeTickets(ctx, dataHash); err != nil {
		log.WithContext(ctx).WithError(err).Error("error checking mempool for duplicate cascade tickets")
		return err
	}

	return nil
}

// CheckDuplicateSenseOrNFTTickets checks duplicate sense or NFT tickets in mempool and inactive transactions list
func (task *DupTicketsDetector) CheckDuplicateSenseOrNFTTickets(ctx context.Context, dataHash []byte) error {
	if err := task.checkMempoolForInProgressSenseOrNFTTickets(ctx, dataHash); err != nil {
		log.WithContext(ctx).WithError(err).Error("error checking mempool for duplicate sense or NFT tickets")
		return err
	}

	return nil
}

// checkMempoolForInProgressCascadeTickets checks the duplicate cascade transactions in mempool
func (task *DupTicketsDetector) checkMempoolForInProgressCascadeTickets(ctx context.Context, dataHash []byte) error {
	transactionIDs, err := task.pastelClient.GetRawMempool(ctx)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("error calling raw mem pool api")
		return err
	}

	if len(transactionIDs) == 0 {
		log.WithContext(ctx).Info("no in-progress transactions in mempool list")
		return nil
	}

	for _, txID := range transactionIDs {
		cascadeTicket, err := task.getCascadeTicket(ctx, txID)
		if err != nil {
			log.WithContext(ctx).
				WithError(err).
				WithField("transaction_id", txID).
				Info("failed to retrieve cascade reg ticket")
			continue
		}

		if cascadeTicket != nil {
			if bytes.Equal(cascadeTicket.DataHash, dataHash) {
				log.WithContext(ctx).
					WithField("transaction_id", txID).
					Error("duplicate ticket registration is already in progress")

				return ErrDuplicateRegistration
			}
		}
	}

	return nil
}

// checkMempoolForInProgressSenseOrNFTTickets checks the duplicate sense or NFT tickets in mempool
func (task *DupTicketsDetector) checkMempoolForInProgressSenseOrNFTTickets(ctx context.Context, dataHash []byte) error {
	transactionIDs, err := task.pastelClient.GetRawMempool(ctx)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("error calling raw mem pool api")
		return err
	}

	if len(transactionIDs) == 0 {
		log.WithContext(ctx).Info("no in-progress transactions in mempool list")
		return nil
	}

	for _, txID := range transactionIDs {
		senseTicket, err := task.getSenseTicket(ctx, txID)
		if err != nil {
			log.WithContext(ctx).WithField("transaction_id", txID).
				WithError(err).
				Info("failed to retrieve retrieve sense reg ticket")
		}

		if senseTicket != nil {
			if bytes.Equal(senseTicket.DataHash, dataHash) {
				log.WithContext(ctx).
					WithField("transaction_id", txID).
					Error("duplicate sense ticket registration is already in progress")

				return ErrDuplicateRegistration

			}
		}

		nftTicket, err := task.getNFTTicket(ctx, txID)
		if err != nil {
			log.WithContext(ctx).WithField("transaction_id", txID).
				WithError(err).
				Info("failed to retrieve retrieve sense reg ticket")

			continue
		}

		if nftTicket != nil {
			if bytes.Equal(nftTicket.AppTicketData.DataHash, dataHash) {
				log.WithContext(ctx).
					WithField("transaction_id", txID).
					Error("duplicate NFT ticket registration is already in progress")

				return ErrDuplicateRegistration
			}
		}
	}

	return nil
}

func (task *DupTicketsDetector) getCascadeTicket(ctx context.Context, txID string) (*pastel.APICascadeTicket, error) {
	regTicket, err := task.pastelClient.ActionRegTicket(ctx, txID)
	if err != nil {
		log.WithContext(ctx).WithError(err).Info("not valid transactionID for cascade reg ticket")
		return nil, err
	}

	if regTicket.ActionTicketData.ActionType != pastel.ActionTypeCascade {
		log.WithContext(ctx).Info("not valid cascade ticket")
		return nil, fmt.Errorf("not valid cascade ticket")
	}

	decTicket, err := pastel.DecodeActionTicket(regTicket.ActionTicketData.ActionTicket)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("failed to decode reg cascade ticket")
		return nil, fmt.Errorf("failed to decode reg cascade ticket")
	}

	regTicket.ActionTicketData.ActionTicketData = *decTicket

	cascadeTicket, err := regTicket.ActionTicketData.ActionTicketData.APICascadeTicket()
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("failed to typecast cascade ticket")
		return nil, fmt.Errorf("failed to typecast cascade ticket: %w", err)
	}

	return cascadeTicket, nil
}

func (task *DupTicketsDetector) getSenseTicket(ctx context.Context, txID string) (*pastel.APISenseTicket, error) {
	regTicket, err := task.pastelClient.ActionRegTicket(ctx, txID)
	if err != nil {
		log.WithContext(ctx).WithError(err).Info("not valid transactionID for sense reg ticket")
		return nil, err
	}

	if regTicket.ActionTicketData.ActionType != pastel.ActionTypeSense {
		log.WithContext(ctx).Info("not valid sense ticket")
		return nil, fmt.Errorf("not valid sense ticket")
	}

	decTicket, err := pastel.DecodeActionTicket(regTicket.ActionTicketData.ActionTicket)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("failed to decode reg sense ticket")
		return nil, fmt.Errorf("failed to decode reg sense ticket")
	}

	regTicket.ActionTicketData.ActionTicketData = *decTicket

	senseTicket, err := regTicket.ActionTicketData.ActionTicketData.APISenseTicket()
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("failed to typecast sense ticket")
		return nil, fmt.Errorf("failed to typecast sense ticket: %w", err)
	}

	return senseTicket, nil
}

func (task *DupTicketsDetector) getNFTTicket(ctx context.Context, txID string) (*pastel.NFTTicket, error) {
	regTicket, err := task.pastelClient.RegTicket(ctx, txID)
	if err != nil {
		log.WithContext(ctx).WithError(err).Info("not valid transactionID for reg ticket")
		return nil, err
	}

	decTicket, err := pastel.DecodeNFTTicket(regTicket.RegTicketData.NFTTicket)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("failed to decode reg nft ticket")
		return nil, fmt.Errorf("failed to decode reg nft ticket")
	}
	regTicket.RegTicketData.NFTTicketData = *decTicket

	return &regTicket.RegTicketData.NFTTicketData, nil
}
