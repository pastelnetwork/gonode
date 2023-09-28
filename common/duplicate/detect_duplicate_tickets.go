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
	CheckDuplicateSenseOrNFTTickets(ctx context.Context, dataHash []byte) error
}

// CheckDuplicateSenseOrNFTTickets checks duplicate sense or NFT tickets in mempool and inactive transactions list
func (task *DupTicketsDetector) CheckDuplicateSenseOrNFTTickets(ctx context.Context, dataHash []byte) error {
	if err := task.checkMempoolForInProgressSenseOrNFTTickets(ctx, dataHash); err != nil {
		log.WithContext(ctx).WithError(err).Error("error checking mempool for duplicate sense or NFT tickets")
		return err
	}

	if err := task.checkInactiveNFTOrSenseTickets(ctx, dataHash); err != nil {
		log.WithContext(ctx).WithError(err).Error("error checking inactive duplicate tickets for sense and NFT")
		return err
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
				Debug("not valid TxID for sense reg ticket")
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

func (task *DupTicketsDetector) checkInactiveNFTOrSenseTickets(ctx context.Context, dataHash []byte) error {
	if err := task.checkInactiveSenseTickets(ctx, dataHash); err != nil {
		log.WithContext(ctx).WithError(err)
		return err
	}

	if err := task.checkInactiveNFTTickets(ctx, dataHash); err != nil {
		log.WithContext(ctx).WithError(err)
		return err
	}

	return nil
}

// checkInactiveSenseTickets checks the duplicate inactive sense tickets
func (task *DupTicketsDetector) checkInactiveSenseTickets(ctx context.Context, dataHash []byte) error {
	inactiveActionTickets, err := task.pastelClient.GetInactiveActionTickets(ctx)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("error retrieving inactive action tickets")
		return errors.Errorf("error retrieving inactive action tickets")
	}

	if len(inactiveActionTickets) == 0 {
		return nil
	}

	for _, ticket := range inactiveActionTickets {
		senseTicket, err := task.getSenseTicket(ctx, ticket.ActTicketData.RegTXID)
		if err != nil {
			continue
		}

		currentBlockCount, err := task.pastelClient.GetBlockCount(ctx)
		if err != nil {
			continue
		}
		blockCountWhenTicketCreated := ticket.ActTicketData.CreatorHeight
		DiffBWCurrentAndCreatorBlockHeight := currentBlockCount - int32(blockCountWhenTicketCreated)

		if bytes.Equal(senseTicket.DataHash, dataHash) {

			log.WithContext(ctx).Info("data hash matched with inactive sense ticket, checking the block threshold")

			if DiffBWCurrentAndCreatorBlockHeight < BlocksThresholdForDuplicateTicket {
				log.WithContext(ctx).Info("duplicate inactive sense ticket found in less then 10000 blocks")
				return errors.Errorf("duplicate inactive sense ticket found in less then 10000 blocks")
			}
		}
	}

	return nil
}

// checkInactiveNFTTickets checks the duplicate inactive NFT tickets
func (task *DupTicketsDetector) checkInactiveNFTTickets(ctx context.Context, dataHash []byte) error {
	inactiveNFTTickets, err := task.pastelClient.GetInactiveNFTTickets(ctx)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("error retrieving inactive action tickets")
		return errors.Errorf("error retrieving action tickets")
	}

	if len(inactiveNFTTickets) == 0 {
		return nil
	}

	for _, ticket := range inactiveNFTTickets {
		nftTicket, err := task.getNFTTicket(ctx, ticket.TXID)
		if err != nil {
			continue
		}

		currentBlockCount, err := task.pastelClient.GetBlockCount(ctx)
		if err != nil {
			continue
		}
		blockCountWhenTicketCreated := nftTicket.BlockNum
		DiffBWCurrentAndCreatorBlockHeight := currentBlockCount - int32(blockCountWhenTicketCreated)

		if bytes.Equal(nftTicket.AppTicketData.DataHash, dataHash) {

			log.WithContext(ctx).Info("data hash matched with inactive NFT ticket, checking the block threshold")

			if DiffBWCurrentAndCreatorBlockHeight < BlocksThresholdForDuplicateTicket {
				log.WithContext(ctx).Info("duplicate inactive ticket found in less then 10000 blocks")
				return errors.Errorf("duplicate inactive ticket found in less then 10000 blocks")
			}
		}
	}

	return nil
}

func (task *DupTicketsDetector) getSenseTicket(ctx context.Context, txID string) (*pastel.APISenseTicket, error) {
	regTicket, err := task.pastelClient.ActionRegTicket(ctx, txID)
	if err != nil {
		return nil, fmt.Errorf("not valid TxID for sense reg ticket")
	}

	if regTicket.ActionTicketData.ActionType != pastel.ActionTypeSense {
		return nil, fmt.Errorf("not valid sense ticket")
	}

	decTicket, err := pastel.DecodeActionTicket(regTicket.ActionTicketData.ActionTicket)
	if err != nil {
		return nil, fmt.Errorf("failed to decode reg sense ticket")
	}

	regTicket.ActionTicketData.ActionTicketData = *decTicket

	senseTicket, err := regTicket.ActionTicketData.ActionTicketData.APISenseTicket()
	if err != nil {
		return nil, fmt.Errorf("failed to typecast sense ticket: %w", err)
	}

	return senseTicket, nil
}

func (task *DupTicketsDetector) getNFTTicket(ctx context.Context, txID string) (*pastel.NFTTicket, error) {
	regTicket, err := task.pastelClient.RegTicket(ctx, txID)
	if err != nil {
		return nil, fmt.Errorf("not valid TxID for nft reg ticket")
	}

	decTicket, err := pastel.DecodeNFTTicket(regTicket.RegTicketData.NFTTicket)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("failed to decode reg nft ticket")
		return nil, fmt.Errorf("failed to decode reg nft ticket")
	}
	regTicket.RegTicketData.NFTTicketData = *decTicket

	return &regTicket.RegTicketData.NFTTicketData, nil
}
