package hermes

import (
	"context"
	"sort"

	"github.com/btcsuite/btcutil/base58"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/pastel"
)

const (
	blocksDeadline = 15000
)

// CleanupInactiveTickets cleans up inactive tickets
func (s *service) CleanupInactiveTickets(ctx context.Context) error {
	count, err := s.pastelClient.GetBlockCount(ctx)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("unable to get block count, skipping cleanup")
		return nil
	}

	errReg := s.cleanupRegTickets(ctx, count)
	if errReg != nil {
		log.WithContext(ctx).WithError(errReg).Error("cleanup Inactive Tickets: cleanupRegTickets failure")
	}

	actionErr := s.cleanupActionTickets(ctx, count)
	if actionErr != nil {
		log.WithContext(ctx).WithError(actionErr).Error("cleanup Inactive Tickets: cleanupActionTickets failure")
	}

	return nil
}

func (s *service) cleanupActionTickets(ctx context.Context, count int32) error {
	actionRegTickets, err := s.pastelClient.ActionTicketsFromBlockHeight(ctx, pastel.TicketTypeInactive, uint64(s.currentActionBlock))
	if err != nil {
		return errors.Errorf("get ActionTickets: %w", err)
	}

	log.WithContext(ctx).WithField("action_tickets_count", len(actionRegTickets)).
		WithField("block-height", s.currentActionBlock).Info("Received action tickets for cleanup")

	sort.Slice(actionRegTickets, func(i, j int) bool {
		return actionRegTickets[i].Height < actionRegTickets[j].Height
	})

	//loop through nft tickets and store newly found nft reg tickets
	for i := 0; i < len(actionRegTickets); i++ {
		if actionRegTickets[i].Height < s.currentActionBlock {
			continue
		}

		if actionRegTickets[i].Height+blocksDeadline <= int(count) {
			continue
		}

		log.WithContext(ctx).WithField("reg_txid", actionRegTickets[i].TXID).
			WithField("action-block-height", s.currentActionBlock).WithField("current block", count).
			Info("Cleaning up action ticket")

		decTicket, err := pastel.DecodeActionTicket(actionRegTickets[i].ActionTicketData.ActionTicket)
		if err != nil {
			log.WithContext(ctx).WithField("reg_txid", actionRegTickets[i].TXID).
				WithError(err).Error("Failed to decode reg ticket")
			continue
		}
		actionRegTickets[i].ActionTicketData.ActionTicketData = *decTicket

		actionTicket := actionRegTickets[i].ActionTicketData.ActionTicketData
		s.cleanupActionTicketData(ctx, actionTicket)

		s.currentActionBlock = actionRegTickets[i].Height

		log.WithContext(ctx).WithField("reg_txid", actionRegTickets[i].TXID).
			Info("cleaned up action ticket")
	}

	return nil
}

func (s *service) removeRegTicketData(ctx context.Context, ticketData pastel.AppTicket) error {
	// Remove RQ File IDs
	for _, rqFileID := range ticketData.RQIDs {
		if err := s.p2p.Delete(ctx, rqFileID); err != nil {
			log.WithContext(ctx).WithError(err).Error("error deleting rqFileID")
		}
	}

	// Remove Thumbnails
	if err := s.p2p.Delete(ctx, base58.Encode(ticketData.Thumbnail1Hash)); err != nil {
		log.WithContext(ctx).WithError(err).Error("error deleting thumbnail1 hash")
	}

	if err := s.p2p.Delete(ctx, base58.Encode(ticketData.Thumbnail2Hash)); err != nil {
		log.WithContext(ctx).WithError(err).Error("error deleting thumbnail2 hash")
	}

	if err := s.p2p.Delete(ctx, base58.Encode(ticketData.PreviewHash)); err != nil {
		log.WithContext(ctx).WithError(err).Error("error deleting preview thumbnail hash")
	}

	// Remove ddAndFpFileIDs
	for _, ddFpFileID := range ticketData.DDAndFingerprintsIDs {
		if err := s.p2p.Delete(ctx, ddFpFileID); err != nil {
			log.WithContext(ctx).WithError(err).Error("error deleting ddFpFileID")
		}
	}

	return nil
}

func (s *service) cleanupRegTickets(ctx context.Context, count int32) error {
	nftRegTickets, err := s.pastelClient.RegTicketsFromBlockHeight(ctx, pastel.TicketTypeInactive, uint64(s.currentNFTBlock))
	if err != nil {
		return errors.Errorf("get RegTickets: %w", err)
	}

	log.WithContext(ctx).WithField("reg_tickets_count", len(nftRegTickets)).
		WithField("block-height", s.currentNFTBlock).Info("Received reg tickets for cleanup")

	sort.Slice(nftRegTickets, func(i, j int) bool {
		return nftRegTickets[i].Height < nftRegTickets[j].Height
	})

	//loop through nft tickets and store newly found nft reg tickets
	for i := 0; i < len(nftRegTickets); i++ {
		if nftRegTickets[i].Height < s.currentNFTBlock {
			continue
		}

		if nftRegTickets[i].Height+blocksDeadline <= int(count) {
			continue
		}

		log.WithContext(ctx).WithField("reg_txid", nftRegTickets[i].TXID).
			WithField("reg-block-height", s.currentNFTBlock).WithField("current block", count).
			Info("cleaning up reg ticket")

		decTicket, err := pastel.DecodeNFTTicket(nftRegTickets[i].RegTicketData.NFTTicket)
		if err != nil {
			log.WithContext(ctx).WithError(err).Error("Failed to decode reg ticket")
			continue
		}
		nftRegTickets[i].RegTicketData.NFTTicketData = *decTicket
		ticketData := nftRegTickets[i].RegTicketData.NFTTicketData.AppTicketData

		s.removeRegTicketData(ctx, ticketData)
		s.currentNFTBlock = nftRegTickets[i].Height

		log.WithContext(ctx).WithField("reg_txid", nftRegTickets[i].TXID).
			Info("cleaned up reg ticket successfully")
	}

	return nil
}

func (s *service) cleanupActionTicketData(ctx context.Context, actionTicket pastel.ActionTicket) error {
	switch actionTicket.ActionType {
	case pastel.ActionTypeCascade:
		cascadeTicket, err := actionTicket.APICascadeTicket()
		if err != nil {
			log.WithContext(ctx).WithField("actionRegTickets.ActionTicketData", actionTicket).
				Warnf("Could not get sense ticket for action ticket data")
		}

		// Remove RQ File IDs
		for _, rqFileID := range cascadeTicket.RQIDs {
			if err := s.p2p.Delete(ctx, rqFileID); err != nil {
				log.WithContext(ctx).WithError(err).Error("error deleting rqFileID")
			}
		}

		if err := s.p2p.Delete(ctx, base58.Encode(cascadeTicket.DataHash)); err != nil {
			log.WithContext(ctx).WithError(err).Error("error deleting thumbnail hash")
		}

	case pastel.ActionTypeSense:
		senseTicket, err := actionTicket.APISenseTicket()
		if err != nil {
			log.WithContext(ctx).WithField("actionRegTickets.ActionTicketData", actionTicket).
				Warnf("Could not get sense ticket for action ticket data")
		}

		for _, ddFpFileID := range senseTicket.DDAndFingerprintsIDs {
			if err := s.p2p.Delete(ctx, ddFpFileID); err != nil {
				log.WithContext(ctx).WithError(err).Error("error deleting ddFpFileID")
			}
		}

		if err := s.p2p.Delete(ctx, base58.Encode(senseTicket.DataHash)); err != nil {
			log.WithContext(ctx).WithError(err).Error("error deleting thumbnail hash")
		}
	}

	return nil
}
