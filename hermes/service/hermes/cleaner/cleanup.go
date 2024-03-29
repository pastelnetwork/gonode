package cleaner

import (
	"context"
	"github.com/pastelnetwork/gonode/hermes/common"
	"sort"
	"time"

	"github.com/btcsuite/btcutil/base58"
	"github.com/pastelnetwork/gonode/common/errgroup"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/pastel"
)

const (
	blocksDeadline  = 15000
	runTaskInterval = 2 * time.Minute
)

// Run cleans up inactive tickets

// Run stores the latest block hash and height to DB if not stored already
func (s *cleanupService) Run(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return errors.Errorf("context done: %w", ctx.Err())
		case <-time.After(runTaskInterval):
			log.WithContext(ctx).Info("cleanup service run() has been invoked")
			if err := s.sync.WaitSynchronization(ctx); err != nil {
				log.WithContext(ctx).WithError(err).Error("error syncing master-node")
				continue
			}

			group, gctx := errgroup.WithContext(ctx)
			group.Go(func() error {
				err := s.run(gctx)
				if err != nil {
					if common.IsP2PConnectionCloseError(err.Error()) {
						s.p2p, err = common.CreateNewP2PConnection(s.config.snHost, s.config.snPort, s.config.sn)
						if err != nil {
							log.WithContext(ctx).WithError(err).Error("unable to initialize new p2p connection for " +
								"cleanup service")
						}
					}
				}

				return nil
			})

			if err := group.Wait(); err != nil {
				log.WithContext(gctx).WithError(err).Errorf("run task failed")
			}
		}
	}

}

func (s cleanupService) Stats(_ context.Context) (map[string]interface{}, error) {
	//cleanup service stats can be implemented here
	return nil, nil
}

func (s *cleanupService) run(ctx context.Context) error {
	count, err := s.pastelClient.GetBlockCount(ctx)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("unable to get block count, skipping cleanup")
		return nil
	}

	errReg := s.cleanupRegTickets(ctx, count)
	if errReg != nil {
		log.WithContext(ctx).WithError(errReg).Error("cleanup Inactive Tickets: cleanupRegTickets failure")
		return err
	}

	actionErr := s.cleanupActionTickets(ctx, count)
	if actionErr != nil {
		log.WithContext(ctx).WithError(actionErr).Error("cleanup Inactive Tickets: cleanupActionTickets failure")
		return err
	}

	return nil
}

func (s *cleanupService) cleanupActionTickets(ctx context.Context, count int32) error {
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

		if !requiresCleanup(actionRegTickets[i].Height, int(count)) {
			continue
		}

		t, err := s.pastelClient.FindActionActByActionRegTxid(ctx, actionRegTickets[i].TXID)
		if err == nil && t != nil {
			continue
		}

		log.WithContext(ctx).WithField("reg_txid", actionRegTickets[i].TXID).
			WithField("action-block-height", s.currentActionBlock).WithField("ticket-height", actionRegTickets[i].Height).
			WithField("current block", int(count)).Info("Cleaning up action ticket")

		decTicket, err := pastel.DecodeActionTicket(actionRegTickets[i].ActionTicketData.ActionTicket)
		if err != nil {
			log.WithContext(ctx).WithField("reg_txid", actionRegTickets[i].TXID).
				WithError(err).Error("Failed to decode reg ticket")
			continue
		}
		actionRegTickets[i].ActionTicketData.ActionTicketData = *decTicket

		actionTicket := actionRegTickets[i].ActionTicketData.ActionTicketData
		err = s.cleanupActionTicketData(ctx, actionTicket)
		if err != nil {
			if common.IsP2PConnectionCloseError(err.Error()) {
				return err
			}
		}

		s.currentActionBlock = actionRegTickets[i].Height

		log.WithContext(ctx).WithField("reg_txid", actionRegTickets[i].TXID).
			Info("cleaned up action ticket")
	}

	return nil
}

func (s *cleanupService) removeRegTicketData(ctx context.Context, ticketData pastel.AppTicket) error {
	// Remove RQ File IDs
	for _, rqFileID := range ticketData.RQIDs {
		if err := s.p2p.Delete(ctx, rqFileID); err != nil {
			if common.IsP2PConnectionCloseError(err.Error()) {
				return err
			}

			log.WithContext(ctx).WithError(err).Error("error deleting rqFileID")
		}
	}

	// Remove Thumbnails
	if err := s.p2p.Delete(ctx, base58.Encode(ticketData.Thumbnail1Hash)); err != nil {
		if common.IsP2PConnectionCloseError(err.Error()) {
			return err
		}

		log.WithContext(ctx).WithError(err).Error("error deleting thumbnail1 hash")
	}

	if err := s.p2p.Delete(ctx, base58.Encode(ticketData.Thumbnail2Hash)); err != nil {
		if common.IsP2PConnectionCloseError(err.Error()) {
			return err
		}

		log.WithContext(ctx).WithError(err).Error("error deleting thumbnail2 hash")
	}

	if err := s.p2p.Delete(ctx, base58.Encode(ticketData.PreviewHash)); err != nil {
		if common.IsP2PConnectionCloseError(err.Error()) {
			return err
		}

		log.WithContext(ctx).WithError(err).Error("error deleting preview thumbnail hash")
	}

	// Remove ddAndFpFileIDs
	for _, ddFpFileID := range ticketData.DDAndFingerprintsIDs {
		if err := s.p2p.Delete(ctx, ddFpFileID); err != nil {
			if common.IsP2PConnectionCloseError(err.Error()) {
				return err
			}

			log.WithContext(ctx).WithError(err).Error("error deleting ddFpFileID")
		}
	}

	return nil
}

func (s *cleanupService) removeCascadeTicketData(ctx context.Context, rqIDs []string, dataHash []byte) error {
	for _, rqFileID := range rqIDs {
		if err := s.p2p.Delete(ctx, rqFileID); err != nil {
			if common.IsP2PConnectionCloseError(err.Error()) {
				return err
			}
			log.WithContext(ctx).WithError(err).Error("error deleting rqFileID")
		}
	}

	if err := s.p2p.Delete(ctx, base58.Encode(dataHash)); err != nil {
		if common.IsP2PConnectionCloseError(err.Error()) {
			return err
		}

		log.WithContext(ctx).WithError(err).Error("error deleting thumbnail hash")
	}

	return nil
}

func (s *cleanupService) removeSenseTicketData(ctx context.Context, ddAndFingerprintIDs []string, dataHash []byte) error {
	for _, ddFpFileID := range ddAndFingerprintIDs {
		if err := s.p2p.Delete(ctx, ddFpFileID); err != nil {
			if common.IsP2PConnectionCloseError(err.Error()) {
				return err
			}

			log.WithContext(ctx).WithError(err).Error("error deleting ddFpFileID")
		}
	}

	if err := s.p2p.Delete(ctx, base58.Encode(dataHash)); err != nil {
		if common.IsP2PConnectionCloseError(err.Error()) {
			return err
		}

		log.WithContext(ctx).WithError(err).Error("error deleting thumbnail hash")
	}

	return nil
}

func (s *cleanupService) cleanupRegTickets(ctx context.Context, count int32) error {
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

		if !requiresCleanup(nftRegTickets[i].Height, int(count)) {
			continue
		}

		t, err := s.pastelClient.FindActByRegTxid(ctx, nftRegTickets[i].TXID)
		if err == nil && t != nil {
			continue
		}

		log.WithContext(ctx).WithField("reg_txid", nftRegTickets[i].TXID).
			WithField("reg-block-height", s.currentNFTBlock).WithField("current block", int(count)).WithField("ticket-height", nftRegTickets[i].Height).
			Info("cleaning up reg ticket")

		decTicket, err := pastel.DecodeNFTTicket(nftRegTickets[i].RegTicketData.NFTTicket)
		if err != nil {
			log.WithContext(ctx).WithError(err).Error("Failed to decode reg ticket")
			continue
		}
		nftRegTickets[i].RegTicketData.NFTTicketData = *decTicket
		ticketData := nftRegTickets[i].RegTicketData.NFTTicketData.AppTicketData

		err = s.removeRegTicketData(ctx, ticketData)
		if err != nil {
			if common.IsP2PConnectionCloseError(err.Error()) {
				return err
			}
		}

		s.currentNFTBlock = nftRegTickets[i].Height

		log.WithContext(ctx).WithField("reg_txid", nftRegTickets[i].TXID).
			Info("cleaned up reg ticket successfully")
	}

	return nil
}

func (s *cleanupService) cleanupActionTicketData(ctx context.Context, actionTicket pastel.ActionTicket) error {
	switch actionTicket.ActionType {
	case pastel.ActionTypeCascade:
		cascadeTicket, err := actionTicket.APICascadeTicket()
		if err != nil {
			log.WithContext(ctx).WithField("actionRegTickets.ActionTicketData", actionTicket).
				Warnf("Could not get sense ticket for action ticket data")
		}

		err = s.removeCascadeTicketData(ctx, cascadeTicket.RQIDs, cascadeTicket.DataHash)
		if err != nil {
			if common.IsP2PConnectionCloseError(err.Error()) {
				return err
			}
		}

	case pastel.ActionTypeSense:
		senseTicket, err := actionTicket.APISenseTicket()
		if err != nil {
			log.WithContext(ctx).WithField("actionRegTickets.ActionTicketData", actionTicket).
				Warnf("Could not get sense ticket for action ticket data")
		}

		err = s.removeSenseTicketData(ctx, senseTicket.DDAndFingerprintsIDs, senseTicket.DataHash)
		if err != nil {
			if common.IsP2PConnectionCloseError(err.Error()) {
				return err
			}
		}
	}

	return nil
}

func requiresCleanup(ticketHeight int, currentBlockHeight int) bool {
	return (ticketHeight + 15000) < currentBlockHeight
}
