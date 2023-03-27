package hermes

import (
	"context"
	"sort"
	"time"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/hermes/service/hermes/domain"
	"github.com/pastelnetwork/gonode/pastel"
)

const (
	// CollectionFinalAllowedBlockHeightDays is to compute time for collection final allowed block height added to activation block height
	CollectionFinalAllowedBlockHeightDays = 7 * 24 * (60 / 2.5)
)

func (s *service) parseCollectionTickets(ctx context.Context) error {
	collectionActTickets, err := s.pastelClient.CollectionActivationTicketsFromBlockHeight(ctx, s.latestNFTBlockHeight)
	if err != nil {
		log.WithError(err).Error("unable to get nft-collection act tickets - exit runtask now")
		return nil
	}
	if len(collectionActTickets) == 0 {
		log.WithContext(ctx).Error("no collection act tickets have been found")
		return nil
	}

	log.WithContext(ctx).Info("nft-collection-act tickets have been retrieved, going to sort on block height")
	sort.Slice(collectionActTickets, func(i, j int) bool {
		return collectionActTickets[i].Height < collectionActTickets[j].Height
	})

	log.WithContext(ctx).WithField("count", len(collectionActTickets)).Info("nft-collection-act tickets have been sorted")

	//track latest block height, but don't set it until we check all the nft reg tickets and the sense tickets.
	lastKnownGoodHeight := s.latestCollectionBlockHeight

	//loop through nft-collection-act tickets and store newly found nft reg tickets
	for i := 0; i < len(collectionActTickets); i++ {
		if collectionActTickets[i].Height < s.latestNFTBlockHeight {
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

		decTicket, err := pastel.DecodeCollectionTicket(regTicket.CollectionRegTicketData.CollectionTicket)
		if err != nil {
			log.WithContext(ctx).WithField("collection_reg_txid", regTxID).WithError(err).Error("Failed to decode reg ticket")
			continue
		}
		regTicket.CollectionRegTicketData.CollectionTicketData = *decTicket

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

		collectionState, datetimeCollectionStateUpdated, err := getCollectionState(ctx,
			regTicket.CollectionRegTicketData.CollectionTicketData.CollectionItemCopyCount,
			regTicket.CollectionRegTicketData.CollectionTicketData.MaxCollectionEntries,
		)
		if err != nil {
			log.WithContext(ctx).WithError(err).Error("undefined collection state")

			return err
		}

		if err := s.store.StoreCollection(ctx, domain.Collection{
			CollectionTicketTXID:                           actTxID,
			CollectionName:                                 regTicket.CollectionRegTicketData.CollectionTicketData.CollectionName,
			CollectionTicketActivationBlockHeight:          collectionActTickets[i].Height,
			CollectionFinalAllowedBlockHeight:              collectionActTickets[i].Height + CollectionFinalAllowedBlockHeightDays,
			MaxPermittedOpenNSFWScore:                      regTicket.CollectionRegTicketData.CollectionTicketData.AppTicketData.MaxPermittedOpenNSFWScore,
			MinimumSimilarityScoreToFirstEntryInCollection: regTicket.CollectionRegTicketData.CollectionTicketData.AppTicketData.MinimumSimilarityScoreToFirstEntryInCollection,
			CollectionState:                                collectionState,
			DatetimeCollectionStateUpdated:                 datetimeCollectionStateUpdated,
		}); err != nil {
			log.WithContext(ctx).WithError(err).Error("failed to store collection")
			continue
		}
		if collectionActTickets[i].Height > lastKnownGoodHeight {
			lastKnownGoodHeight = collectionActTickets[i].Height
		}
	}
	//loop through action tickets and store newly found nft reg tickets

	if lastKnownGoodHeight > s.latestCollectionBlockHeight {
		s.latestCollectionBlockHeight = lastKnownGoodHeight
	}

	log.WithContext(ctx).WithField("latest NFT-Collection block-height", s.latestCollectionBlockHeight).Debugf("hermes successfully scanned to latest block height")

	return nil
}

func getCollectionState(ctx context.Context, currentNoOfCollectionEntries, maxNoOfCollectionEntries uint) (collectionState domain.CollectionState, datetimeCollectionStateUpdated string, err error) {
	ent := log.WithContext(ctx).WithField("current_no_of_col_entries", currentNoOfCollectionEntries).
		WithField("max_no_of_collection_entries", maxNoOfCollectionEntries)

	if currentNoOfCollectionEntries == maxNoOfCollectionEntries {
		ent.Info("current and max no of collection entries are equal, collection state is finalized")
		return domain.FinalizedCollectionState, time.Now().Format(time.RFC3339), nil
	}

	if currentNoOfCollectionEntries < maxNoOfCollectionEntries {
		ent.Info("current no of collection entries are less than max no of collection entries, collection state is in-process")

		return domain.InProcessCollectionState, time.Now().Format(time.RFC3339), nil
	}

	ent.Error("current no of collection entries are less than max no of collection entries, collection state is in-process")
	err = errors.Errorf("current no of collection entries are greater than max no of collection entries")
	return domain.UndefinedCollectionState, "", err
}
