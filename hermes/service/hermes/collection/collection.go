package collection

import (
	"context"
	"sort"
	"time"

	"github.com/pastelnetwork/gonode/common/errgroup"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/hermes/domain"
)

const (
	// CollectionFinalAllowedBlockHeightDays is to compute time for collection final allowed block height added to activation block height
	CollectionFinalAllowedBlockHeightDays = 7 * 24 * (60 / 2.5)
	runTaskInterval                       = 2 * time.Minute
)

func (s *collectionService) Run(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return errors.Errorf("context done: %w", ctx.Err())
		case <-time.After(runTaskInterval):
			log.WithContext(ctx).Info("collection service run() has been invoked")
			if err := s.sync.WaitSynchronization(ctx); err != nil {
				log.WithContext(ctx).WithError(err).Error("error syncing master-node")
				continue
			}

			groupA, cctx := errgroup.WithContext(ctx)
			groupB, gctx := errgroup.WithContext(ctx)

			groupA.Go(func() error {
				return s.parseCollectionTickets(cctx)
			})

			groupB.Go(func() error {
				return s.finalizeCollections(gctx)
			})

			if err := groupA.Wait(); err != nil {
				log.WithContext(cctx).WithError(err).Errorf("parseCollectionTickets failed")
			}

			if err := groupB.Wait(); err != nil {
				log.WithContext(gctx).WithError(err).Errorf("finalizeCollections failed")
			}
		}
	}
}

func (s collectionService) parseCollectionTickets(ctx context.Context) error {
	collectionActTickets, err := s.pastelClient.CollectionActivationTicketsFromBlockHeight(ctx, s.latestCollectionBlockHeight)
	if err != nil {
		log.WithError(err).Error("unable to get nft-collection act tickets - exit runtask now")
		return nil
	}
	if len(collectionActTickets) == 0 {
		log.WithContext(ctx).Debug("no collection act tickets have been found")
		return nil
	}

	log.WithContext(ctx).Info("nft-collection-act tickets have been retrieved, going to sort on block height")
	sort.Slice(collectionActTickets, func(i, j int) bool {
		return collectionActTickets[i].Height < collectionActTickets[j].Height
	})

	log.WithContext(ctx).WithField("count", len(collectionActTickets)).Info("nft-collection-act tickets have been sorted")
	lastKnownGoodHeight := s.latestCollectionBlockHeight

	for i := 0; i < len(collectionActTickets); i++ {
		if collectionActTickets[i].Height < s.latestCollectionBlockHeight {
			continue
		}

		actTxID := collectionActTickets[i].TXID
		regTxID := collectionActTickets[i].ActTicketData.RegTXID

		log.WithContext(ctx).WithField("collection_act_txid", actTxID).WithField("collection_reg_txid", regTxID).
			Info("Found new nft-collection-act ticket")

		regTicket, err := s.pastelClient.CollectionRegTicket(ctx, regTxID)
		if err != nil {
			log.WithContext(ctx).WithError(err).WithField("collection_reg_txid", regTxID).
				WithField("collection_act_txid", actTxID).Error("unable to get reg ticket")
			continue
		}

		log.WithContext(ctx).WithField("collection_reg_txid", regTxID).Info("Found reg ticket for nft-collection-act ticket")

		existsInDatabase, err := s.store.IfCollectionExists(ctx, actTxID)
		if existsInDatabase {
			log.WithContext(ctx).WithField("collection_reg_txid", regTxID).Info("collection already exist in database, skipping")
			//can't directly update latest block height from here - if there's another ticket in this block we don't want to skip
			if collectionActTickets[i].Height > lastKnownGoodHeight {
				lastKnownGoodHeight = collectionActTickets[i].Height
			}

			continue
		}
		if err != nil {
			log.WithContext(ctx).WithField("collection_reg_txid", regTxID).Error("Could not properly query the database for this collection")
			continue
		}

		if err := s.store.StoreCollection(ctx, domain.Collection{
			CollectionTicketTXID:                           actTxID,
			CollectionName:                                 regTicket.CollectionRegTicketData.CollectionTicketData.CollectionName,
			CollectionTicketActivationBlockHeight:          collectionActTickets[i].Height,
			CollectionFinalAllowedBlockHeight:              int(regTicket.CollectionRegTicketData.CollectionTicketData.CollectionFinalAllowedBlockHeight),
			MaxPermittedOpenNSFWScore:                      regTicket.CollectionRegTicketData.CollectionTicketData.AppTicketData.MaxPermittedOpenNSFWScore,
			MinimumSimilarityScoreToFirstEntryInCollection: regTicket.CollectionRegTicketData.CollectionTicketData.AppTicketData.MinimumSimilarityScoreToFirstEntryInCollection,
			CollectionState:                                domain.InProcessCollectionState,
			DatetimeCollectionStateUpdated:                 time.Now().Format(time.RFC3339),
		}); err != nil {
			log.WithContext(ctx).WithError(err).Error("failed to store collection")
			continue
		}
		if collectionActTickets[i].Height > lastKnownGoodHeight {
			lastKnownGoodHeight = collectionActTickets[i].Height
		}
	}

	if lastKnownGoodHeight > s.latestCollectionBlockHeight {
		s.latestCollectionBlockHeight = lastKnownGoodHeight
	}

	log.WithContext(ctx).WithField("latest NFT-Collection block-height", s.latestCollectionBlockHeight).Debugf("hermes successfully scanned to latest block height")

	return nil
}

func (s *collectionService) Stats(_ context.Context) (map[string]interface{}, error) {
	//collection service stats can be implemented here
	return nil, nil
}
